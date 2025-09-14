package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

const (
	// DefaultSSHKeyName is the default SSH key filename
	DefaultSSHKeyName = "id_rsa"
	// DefaultSSHPort is the default SSH port
	DefaultSSHPort = 22
)

// FindSSHKey finds the first available SSH key in the user's .ssh directory
func FindSSHKey() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	sshDir := filepath.Join(homeDir, ".ssh")
	keyPaths := []string{
		filepath.Join(sshDir, "id_rsa"),
		filepath.Join(sshDir, "id_ed25519"),
		filepath.Join(sshDir, "id_ecdsa"),
	}

	for _, keyPath := range keyPaths {
		if _, err := os.Stat(keyPath); err == nil {
			return keyPath, nil
		}
	}

	return "", fmt.Errorf("no SSH key found in %s", sshDir)
}

// ValidateSSHKey validates that an SSH key file exists and is parseable
func ValidateSSHKey(keyPath string) error {
	// Check if file exists
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return fmt.Errorf("SSH key file does not exist: %s", keyPath)
	}

	// Try to parse the key
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("failed to read SSH key file: %w", err)
	}

	_, err = ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return fmt.Errorf("invalid SSH private key: %w", err)
	}

	return nil
}

// ParseHostPort parses a host:port string and returns the host and port
func ParseHostPort(hostPort string) (host string, port int, err error) {
	if hostPort == "" {
		return "", 0, fmt.Errorf("host cannot be empty")
	}

	parts := strings.Split(hostPort, ":")
	if len(parts) == 1 {
		return parts[0], DefaultSSHPort, nil
	}

	if len(parts) == 2 {
		host = parts[0]
		if parts[1] != "" {
			var portInt int
			if _, err := fmt.Sscanf(parts[1], "%d", &portInt); err != nil {
				return "", 0, fmt.Errorf("invalid port: %s", parts[1])
			}
			if portInt <= 0 || portInt > 65535 {
				return "", 0, fmt.Errorf("port must be between 1 and 65535, got %d", portInt)
			}
			port = portInt
		} else {
			port = DefaultSSHPort
		}
		return host, port, nil
	}

	return "", 0, fmt.Errorf("invalid host:port format: %s", hostPort)
}

// NormalizeHost normalizes a hostname by removing protocols and extra characters
func NormalizeHost(host string) string {
	host = strings.TrimSpace(host)
	host = strings.ToLower(host)

	// Remove protocol if present
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")

	// Remove trailing slash
	host = strings.TrimSuffix(host, "/")

	return host
}
