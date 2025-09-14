package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	config := New()

	if config.Port != DefaultPort {
		t.Errorf("Expected port %d, got %d", DefaultPort, config.Port)
	}

	if config.User != DefaultUser {
		t.Errorf("Expected user %s, got %s", DefaultUser, config.User)
	}

	if config.Defaults.VolumePath != DefaultVolumePath {
		t.Errorf("Expected volume path %s, got %s", DefaultVolumePath, config.Defaults.VolumePath)
	}

	if config.Defaults.Network != DefaultNetwork {
		t.Errorf("Expected network %s, got %s", DefaultNetwork, config.Defaults.Network)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		shouldErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Host:       "192.168.1.100",
				User:       "admin",
				Port:       22,
				SSHKeyPath: createTempSSHKey(t),
			},
			shouldErr: false,
		},
		{
			name: "missing host",
			config: &Config{
				User:       "admin",
				Port:       22,
				SSHKeyPath: createTempSSHKey(t),
			},
			shouldErr: true,
		},
		{
			name: "missing user",
			config: &Config{
				Host:       "192.168.1.100",
				Port:       22,
				SSHKeyPath: createTempSSHKey(t),
			},
			shouldErr: true,
		},
		{
			name: "invalid port",
			config: &Config{
				Host:       "192.168.1.100",
				User:       "admin",
				Port:       -1,
				SSHKeyPath: createTempSSHKey(t),
			},
			shouldErr: true,
		},
		{
			name: "missing SSH key",
			config: &Config{
				Host:       "192.168.1.100",
				User:       "admin",
				Port:       22,
				SSHKeyPath: "/nonexistent/key",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()

	// Override home directory for test
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tempDir)

	// Create test config
	config := &Config{
		Host:       "192.168.1.100",
		User:       "testuser",
		Port:       2222,
		SSHKeyPath: createTempSSHKey(t),
		Defaults: struct {
			VolumePath string `yaml:"volume_path"`
			Network    string `yaml:"network,omitempty"`
		}{
			VolumePath: "/volume1/test",
			Network:    "testnet",
		},
	}

	// Save config
	if err := config.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loadedConfig, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Compare configs
	if loadedConfig.Host != config.Host {
		t.Errorf("Expected host %s, got %s", config.Host, loadedConfig.Host)
	}
	if loadedConfig.User != config.User {
		t.Errorf("Expected user %s, got %s", config.User, loadedConfig.User)
	}
	if loadedConfig.Port != config.Port {
		t.Errorf("Expected port %d, got %d", config.Port, loadedConfig.Port)
	}
	if loadedConfig.SSHKeyPath != config.SSHKeyPath {
		t.Errorf("Expected SSH key path %s, got %s", config.SSHKeyPath, loadedConfig.SSHKeyPath)
	}
	if loadedConfig.Defaults.VolumePath != config.Defaults.VolumePath {
		t.Errorf("Expected volume path %s, got %s", config.Defaults.VolumePath, loadedConfig.Defaults.VolumePath)
	}
	if loadedConfig.Defaults.Network != config.Defaults.Network {
		t.Errorf("Expected network %s, got %s", config.Defaults.Network, loadedConfig.Defaults.Network)
	}
}

func TestGetConfigPath(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()

	// Override home directory for test
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tempDir)

	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("Failed to get config path: %v", err)
	}

	expectedPath := filepath.Join(tempDir, ConfigDir, ConfigFile)
	if path != expectedPath {
		t.Errorf("Expected config path %s, got %s", expectedPath, path)
	}

	// Check if directory was created
	configDir := filepath.Dir(path)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("Config directory was not created: %s", configDir)
	}
}

func createTempSSHKey(t *testing.T) string {
	tempFile, err := os.CreateTemp("", "ssh_key_*")
	if err != nil {
		t.Fatalf("Failed to create temp SSH key file: %v", err)
	}
	defer tempFile.Close()

	// Write dummy SSH key content
	if _, err := tempFile.WriteString("-----BEGIN OPENSSH PRIVATE KEY-----\ntest\n-----END OPENSSH PRIVATE KEY-----\n"); err != nil {
		t.Fatalf("Failed to write temp SSH key: %v", err)
	}

	return tempFile.Name()
}
