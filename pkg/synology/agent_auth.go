package synology

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// connectSSHWithAgent attempts to connect using ssh-agent if available
func (c *Connection) connectSSHWithAgent() error {
	// Check if SSH_AUTH_SOCK is set
	authSock := os.Getenv("SSH_AUTH_SOCK")
	if authSock == "" {
		return fmt.Errorf("ssh-agent not available (SSH_AUTH_SOCK not set)")
	}

	// Connect to ssh-agent
	agentConn, err := net.Dial("unix", authSock)
	if err != nil {
		return fmt.Errorf("failed to connect to ssh-agent: %w", err)
	}
	defer agentConn.Close()

	agentClient := agent.NewClient(agentConn)

	// Get available signers from agent
	signers, err := agentClient.Signers()
	if err != nil {
		return fmt.Errorf("failed to get signers from ssh-agent: %w", err)
	}

	if len(signers) == 0 {
		return fmt.Errorf("no keys available in ssh-agent")
	}

	// Configure SSH client with agent authentication
	sshConfig := &ssh.ClientConfig{
		User: c.config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Implement proper host key verification
	}

	// Connect to SSH server
	address := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to dial SSH with agent: %w", err)
	}

	c.sshClient = client
	return nil
}

// hasSSHAgent checks if ssh-agent is available and has keys
func hasSSHAgent() bool {
	authSock := os.Getenv("SSH_AUTH_SOCK")
	if authSock == "" {
		return false
	}

	// Try to connect to agent
	agentConn, err := net.Dial("unix", authSock)
	if err != nil {
		return false
	}
	defer agentConn.Close()

	agentClient := agent.NewClient(agentConn)
	signers, err := agentClient.Signers()
	if err != nil {
		return false
	}

	return len(signers) > 0
}
