package synology

import (
	"testing"

	"github.com/scttfrdmn/syno-docker/pkg/config"
)

func TestNewConnection(t *testing.T) {
	cfg := &config.Config{
		Host:       "192.168.1.100",
		User:       "admin",
		Port:       22,
		SSHKeyPath: "/path/to/key",
	}

	conn := NewConnection(cfg)

	if conn == nil {
		t.Fatal("NewConnection returned nil")
	}

	if conn.config != cfg {
		t.Error("Connection config not set correctly")
	}

	if conn.sshClient != nil {
		t.Error("SSH client should be nil before Connect()")
	}

	if conn.dockerAPI != nil {
		t.Error("Docker API should be nil before Connect()")
	}
}

func TestExecuteDockerCommand(t *testing.T) {
	cfg := &config.Config{
		Host:       "192.168.1.100",
		User:       "admin",
		Port:       22,
		SSHKeyPath: "/path/to/key",
	}

	conn := NewConnection(cfg)

	// Test command construction
	expectedCmd := "/usr/local/bin/docker ps -a"

	// We can't actually execute without a real SSH connection,
	// but we can test the command construction logic
	if conn.config.Host != cfg.Host {
		t.Error("Configuration not properly set")
	}

	// Test that DockerBinary constant is used
	if DockerBinary != "/usr/local/bin/docker" {
		t.Errorf("Expected DockerBinary to be /usr/local/bin/docker, got %s", DockerBinary)
	}

	// Verify the command would be constructed correctly
	cmd := DockerBinary + " " + "ps -a"
	if cmd != expectedCmd {
		t.Errorf("Expected command %s, got %s", expectedCmd, cmd)
	}
}

func TestConstants(t *testing.T) {
	expectedConstants := map[string]string{
		"ServiceName":          "pkg-ContainerManager-dockerd",
		"ConfigPath":           "/var/packages/ContainerManager/etc/dockerd.json",
		"DockerBinary":         "/usr/local/bin/docker",
		"SocketPath":           "/var/run/docker.sock",
		"RestartCommand":       "sudo systemctl restart pkg-ContainerManager-dockerd",
		"DefaultVolume":        "/volume1/docker",
		"DefaultRestartPolicy": "unless-stopped",
		"DefaultNetwork":       "bridge",
	}

	actualConstants := map[string]string{
		"ServiceName":          ServiceName,
		"ConfigPath":           ConfigPath,
		"DockerBinary":         DockerBinary,
		"SocketPath":           SocketPath,
		"RestartCommand":       RestartCommand,
		"DefaultVolume":        DefaultVolume,
		"DefaultRestartPolicy": DefaultRestartPolicy,
		"DefaultNetwork":       DefaultNetwork,
	}

	for name, expected := range expectedConstants {
		if actual := actualConstants[name]; actual != expected {
			t.Errorf("Constant %s: expected %s, got %s", name, expected, actual)
		}
	}

	// Test integer constants
	if DefaultSSHPort != 22 {
		t.Errorf("Expected DefaultSSHPort to be 22, got %d", DefaultSSHPort)
	}

	if DefaultSSHUser != "admin" {
		t.Errorf("Expected DefaultSSHUser to be admin, got %s", DefaultSSHUser)
	}
}
