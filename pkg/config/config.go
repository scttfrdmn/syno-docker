package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DefaultUser       = "admin"
	DefaultPort       = 22
	DefaultVolumePath = "/volume1/docker"
	DefaultNetwork    = "bridge"
	ConfigDir         = ".synodeploy"
	ConfigFile        = "config.yaml"
)

type Config struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port,omitempty"`
	User       string `yaml:"user"`
	SSHKeyPath string `yaml:"ssh_key_path"`

	Defaults struct {
		VolumePath string `yaml:"volume_path"`
		Network    string `yaml:"network,omitempty"`
	} `yaml:"defaults"`
}

func New() *Config {
	return &Config{
		Port: DefaultPort,
		User: DefaultUser,
		Defaults: struct {
			VolumePath string `yaml:"volume_path"`
			Network    string `yaml:"network,omitempty"`
		}{
			VolumePath: DefaultVolumePath,
			Network:    DefaultNetwork,
		},
	}
}

func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host is required")
	}
	if c.User == "" {
		return fmt.Errorf("user is required")
	}
	if c.SSHKeyPath == "" {
		return fmt.Errorf("ssh_key_path is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	// Check if SSH key exists
	if _, err := os.Stat(c.SSHKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("SSH key not found at %s", c.SSHKeyPath)
	}

	return nil
}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ConfigDir)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(configDir, ConfigFile), nil
}

func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration not found. Run 'synodeploy init <host>' first")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := New()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
