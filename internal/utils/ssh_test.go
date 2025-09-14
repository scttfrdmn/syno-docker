package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseHostPort(t *testing.T) {
	tests := []struct {
		input        string
		expectedHost string
		expectedPort int
		shouldError  bool
	}{
		{"192.168.1.100", "192.168.1.100", 22, false},
		{"192.168.1.100:2222", "192.168.1.100", 2222, false},
		{"hostname", "hostname", 22, false},
		{"hostname:22", "hostname", 22, false},
		{"hostname:", "hostname", 22, false},
		{"hostname:invalid", "", 0, true},
		{"hostname:-1", "", 0, true},
		{"hostname:65536", "", 0, true},
		{"", "", 0, true},
		{"host:port:extra", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			host, port, err := ParseHostPort(tt.input)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if host != tt.expectedHost {
				t.Errorf("Expected host %s, got %s", tt.expectedHost, host)
			}

			if port != tt.expectedPort {
				t.Errorf("Expected port %d, got %d", tt.expectedPort, port)
			}
		})
	}
}

func TestNormalizeHost(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"192.168.1.100", "192.168.1.100"},
		{"  192.168.1.100  ", "192.168.1.100"},
		{"HOSTNAME", "hostname"},
		{"http://example.com", "example.com"},
		{"https://example.com", "example.com"},
		{"example.com/", "example.com"},
		{"https://example.com/path/", "example.com/path"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := NormalizeHost(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestValidateSSHKey(t *testing.T) {
	// Create a temporary SSH key file
	tempDir := t.TempDir()

	// Create valid SSH key file
	validKeyPath := filepath.Join(tempDir, "valid_key")
	validKeyContent := `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEAuKv1SNQM2QWJ7uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJ
QJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ
2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2u
J2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2
ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJ
QJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ
2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2u
J2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2ZJQJ2uJ2
ZJQwIDAQABAAABAQCyaQqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzq
mZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqm
ZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZ
K5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK
5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5
qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5q
XxzqmZK5qXxzqmZK5qXxzqAAAAgQCZ5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqm
ZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZ
K5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK
5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5qXxzqmZK5
qXxzqmZK5qXxzq
-----END OPENSSH PRIVATE KEY-----`

	if err := os.WriteFile(validKeyPath, []byte(validKeyContent), 0600); err != nil {
		t.Fatalf("Failed to create valid key file: %v", err)
	}

	// Create invalid SSH key file
	invalidKeyPath := filepath.Join(tempDir, "invalid_key")
	if err := os.WriteFile(invalidKeyPath, []byte("not a valid key"), 0600); err != nil {
		t.Fatalf("Failed to create invalid key file: %v", err)
	}

	tests := []struct {
		name      string
		keyPath   string
		shouldErr bool
	}{
		{"nonexistent key", filepath.Join(tempDir, "nonexistent"), true},
		{"invalid key content", invalidKeyPath, true},
		{"valid key", validKeyPath, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSSHKey(tt.keyPath)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestFindSSHKey(t *testing.T) {
	// Create a temporary home directory
	tempDir := t.TempDir()
	sshDir := filepath.Join(tempDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatalf("Failed to create .ssh directory: %v", err)
	}

	// Override home directory for test
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tempDir)

	tests := []struct {
		name       string
		keyFiles   []string
		shouldFind bool
		expected   string
	}{
		{
			name:       "id_rsa found",
			keyFiles:   []string{"id_rsa"},
			shouldFind: true,
			expected:   "id_rsa",
		},
		{
			name:       "id_ed25519 found",
			keyFiles:   []string{"id_ed25519"},
			shouldFind: true,
			expected:   "id_ed25519",
		},
		{
			name:       "multiple keys - prefers id_rsa",
			keyFiles:   []string{"id_ed25519", "id_rsa", "id_ecdsa"},
			shouldFind: true,
			expected:   "id_rsa",
		},
		{
			name:       "no keys found",
			keyFiles:   []string{},
			shouldFind: false,
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up SSH directory
			os.RemoveAll(sshDir)
			os.MkdirAll(sshDir, 0700)

			// Create key files
			for _, keyFile := range tt.keyFiles {
				keyPath := filepath.Join(sshDir, keyFile)
				if err := os.WriteFile(keyPath, []byte("dummy key"), 0600); err != nil {
					t.Fatalf("Failed to create key file %s: %v", keyFile, err)
				}
			}

			keyPath, err := FindSSHKey()

			if tt.shouldFind {
				if err != nil {
					t.Errorf("Expected to find key but got error: %v", err)
					return
				}

				expectedPath := filepath.Join(sshDir, tt.expected)
				if keyPath != expectedPath {
					t.Errorf("Expected key path %s, got %s", expectedPath, keyPath)
				}
			} else {
				if err == nil {
					t.Error("Expected error but found key")
				}
			}
		})
	}
}
