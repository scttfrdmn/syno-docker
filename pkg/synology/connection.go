package synology

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"

	"github.com/scttfrdmn/synodeploy/pkg/config"
)

type Connection struct {
	config    *config.Config
	sshClient *ssh.Client
	dockerAPI *client.Client
}

func NewConnection(cfg *config.Config) *Connection {
	return &Connection{
		config: cfg,
	}
}

func (c *Connection) Connect() error {
	if err := c.connectSSH(); err != nil {
		return errors.Wrap(err, "failed to establish SSH connection")
	}

	if err := c.setupDockerClient(); err != nil {
		return errors.Wrap(err, "failed to setup Docker client")
	}

	return nil
}

func (c *Connection) connectSSH() error {
	// Read SSH private key
	keyBytes, err := os.ReadFile(c.config.SSHKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read SSH key: %w", err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return fmt.Errorf("failed to parse SSH key: %w", err)
	}

	// Configure SSH client
	sshConfig := &ssh.ClientConfig{
		User: c.config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Implement proper host key verification
	}

	// Connect to SSH server
	address := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to dial SSH: %w", err)
	}

	c.sshClient = client
	return nil
}

func (c *Connection) setupDockerClient() error {
	// Create Docker client that uses SSH connection
	dockerHost := fmt.Sprintf("ssh://%s@%s:%d", c.config.User, c.config.Host, c.config.Port)

	dockerClient, err := client.NewClientWithOpts(
		client.WithHost(dockerHost),
		client.WithVersion("1.43"),
	)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Test Docker connection
	ctx := context.Background()
	_, err = dockerClient.Ping(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping Docker daemon: %w", err)
	}

	c.dockerAPI = dockerClient
	return nil
}

func (c *Connection) Close() error {
	var errs []string

	if c.dockerAPI != nil {
		if err := c.dockerAPI.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("Docker client: %v", err))
		}
	}

	if c.sshClient != nil {
		if err := c.sshClient.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("SSH client: %v", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %s", strings.Join(errs, ", "))
	}

	return nil
}

func (c *Connection) ExecuteCommand(cmd string) (string, error) {
	if c.sshClient == nil {
		return "", fmt.Errorf("SSH client not connected")
	}

	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Set up environment to include Docker binary path
	session.Setenv("PATH", "/usr/local/bin:/usr/bin:/bin")

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(output), fmt.Errorf("command failed: %w", err)
	}

	return string(output), nil
}

func (c *Connection) ExecuteDockerCommand(args []string) (string, error) {
	// Always use full path to Docker binary
	cmd := DockerBinary + " " + strings.Join(args, " ")
	return c.ExecuteCommand(cmd)
}

func (c *Connection) GetDockerClient() *client.Client {
	return c.dockerAPI
}

func (c *Connection) TestConnection() error {
	// Test SSH connection
	if _, err := c.ExecuteCommand("echo 'SSH connection test'"); err != nil {
		return fmt.Errorf("SSH connection test failed: %w", err)
	}

	// Test Docker command
	if _, err := c.ExecuteDockerCommand([]string{"version", "--format", "'{{.Server.Version}}'"}); err != nil {
		return fmt.Errorf("Docker connection test failed: %w", err)
	}

	return nil
}

func (c *Connection) StreamCommand(cmd string, stdout, stderr io.Writer) error {
	if c.sshClient == nil {
		return fmt.Errorf("SSH client not connected")
	}

	session, err := c.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Set up environment
	session.Setenv("PATH", "/usr/local/bin:/usr/bin:/bin")

	// Connect streams
	session.Stdout = stdout
	session.Stderr = stderr

	// Run command
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}
