package deploy

import (
	"fmt"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"

	"github.com/scttfrdmn/syno-docker/pkg/synology"
)

// ContainerOptions defines options for container deployment
type ContainerOptions struct {
	Image       string
	Name        string
	Ports       []string // ["8080:80", "443:443"]
	Volumes     []string // ["/volume1/data:/app/data"]
	Env         []string // ["KEY=value"]
	Restart     string   // "unless-stopped"
	NetworkMode string
	WorkingDir  string
	Command     []string
	User        string
}

// ContainerInfo represents container information
type ContainerInfo struct {
	ID     string
	Name   string
	Image  string
	Status string
	Ports  []string
}

// NewContainerOptions creates new container options with defaults
func NewContainerOptions(image string) *ContainerOptions {
	return &ContainerOptions{
		Image:       image,
		Restart:     synology.DefaultRestartPolicy,
		NetworkMode: synology.DefaultNetwork,
	}
}

// Container deploys a container using direct Docker commands over SSH
func Container(conn *synology.Connection, opts *ContainerOptions) (string, error) {
	if opts.Name == "" {
		opts.Name = generateContainerName(opts.Image)
	}

	// Build docker run command
	dockerArgs := []string{"run", "-d"}

	// Add name
	dockerArgs = append(dockerArgs, "--name", opts.Name)

	// Add ports
	for _, port := range opts.Ports {
		dockerArgs = append(dockerArgs, "-p", port)
	}

	// Add volumes
	for _, volume := range opts.Volumes {
		dockerArgs = append(dockerArgs, "-v", volume)
	}

	// Add environment variables
	for _, env := range opts.Env {
		dockerArgs = append(dockerArgs, "-e", env)
	}

	// Add restart policy
	if opts.Restart != "" {
		dockerArgs = append(dockerArgs, "--restart", opts.Restart)
	}

	// Add network
	if opts.NetworkMode != "" {
		dockerArgs = append(dockerArgs, "--network", opts.NetworkMode)
	}

	// Add user
	if opts.User != "" {
		dockerArgs = append(dockerArgs, "--user", opts.User)
	}

	// Add working directory
	if opts.WorkingDir != "" {
		dockerArgs = append(dockerArgs, "-w", opts.WorkingDir)
	}

	// Add image
	dockerArgs = append(dockerArgs, opts.Image)

	// Add command
	dockerArgs = append(dockerArgs, opts.Command...)

	// Pull image first
	fmt.Printf("Pulling image %s...\n", opts.Image)
	if _, err := conn.ExecuteDockerCommand([]string{"pull", opts.Image}); err != nil {
		return "", errors.Wrap(err, "failed to pull image")
	}

	// Run container
	fmt.Printf("Creating and starting container %s...\n", opts.Name)
	output, err := conn.ExecuteDockerCommand(dockerArgs)
	if err != nil {
		return "", errors.Wrapf(err, "failed to run container: %s", output)
	}

	containerID := strings.TrimSpace(output)
	return containerID, nil
}

// ListContainers lists containers using direct Docker commands
func ListContainers(conn *synology.Connection, all bool) ([]ContainerInfo, error) {
	args := []string{"ps", "--format", "'table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}'"}
	if all {
		args = append(args, "-a")
	}

	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list containers")
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) <= 1 {
		return []ContainerInfo{}, nil
	}

	var containers []ContainerInfo
	for _, line := range lines[1:] { // Skip header
		fields := strings.Fields(line)
		if len(fields) >= 4 {
			container := ContainerInfo{
				ID:     fields[0],
				Name:   fields[1],
				Image:  fields[2],
				Status: fields[3],
			}
			if len(fields) > 4 {
				container.Ports = fields[4:]
			}
			containers = append(containers, container)
		}
	}

	return containers, nil
}

// RemoveContainer removes a container using direct Docker commands
func RemoveContainer(conn *synology.Connection, nameOrID string, force bool) error {
	args := []string{"rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, nameOrID)

	fmt.Printf("Removing container %s...\n", nameOrID)
	output, err := conn.ExecuteDockerCommand(args)
	if err != nil {
		return errors.Wrapf(err, "failed to remove container: %s", output)
	}

	return nil
}

func generateContainerName(image string) string {
	// Extract image name without registry and tag
	parts := strings.Split(image, "/")
	name := parts[len(parts)-1]

	// Remove tag if present
	if idx := strings.Index(name, ":"); idx != -1 {
		name = name[:idx]
	}

	return name
}

// GetDockerClient returns a placeholder - not used in simple implementation
func GetDockerClient() (*client.Client, error) {
	// This is a placeholder for compatibility
	// The simple implementation uses SSH commands directly
	return nil, fmt.Errorf("docker client not needed in simple implementation")
}

// TestDockerConnection tests Docker availability over SSH
func TestDockerConnection(conn *synology.Connection) error {
	// Test Docker command
	if _, err := conn.ExecuteDockerCommand([]string{"version", "--format", "'{{.Server.Version}}'"}); err != nil {
		return fmt.Errorf("docker connection test failed: %w", err)
	}

	return nil
}
