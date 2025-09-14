package utils

import (
	"testing"
)

func TestValidateDockerImage(t *testing.T) {
	tests := []struct {
		image     string
		shouldErr bool
	}{
		{"nginx", false},
		{"nginx:latest", false},
		{"nginx:1.20", false},
		{"docker.io/nginx", false},
		{"docker.io/nginx:latest", false},
		{"gcr.io/project/image:tag", false},
		{"localhost:5000/myimage:v1.0.0", false},
		{"", true},
		{"NGINX", true}, // uppercase not allowed
		{"nginx::", true},
		{"nginx/", true},
		{"/nginx", true},
		{"nginx:-tag", true},
	}

	for _, tt := range tests {
		t.Run(tt.image, func(t *testing.T) {
			err := ValidateDockerImage(tt.image)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateContainerName(t *testing.T) {
	tests := []struct {
		name      string
		shouldErr bool
	}{
		{"nginx", false},
		{"nginx-server", false},
		{"nginx_server", false},
		{"nginx.server", false},
		{"nginx123", false},
		{"123nginx", false},
		{"", true},
		{"-nginx", true},                  // cannot start with hyphen
		{".nginx", true},                  // cannot start with period
		{"_nginx", true},                  // cannot start with underscore
		{"nginx server", true},            // spaces not allowed
		{"nginx/server", true},            // slashes not allowed
		{string(make([]byte, 254)), true}, // too long
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateContainerName(tt.name)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidatePortMapping(t *testing.T) {
	tests := []struct {
		mapping   string
		shouldErr bool
	}{
		{"8080:80", false},
		{"80:80", false},
		{"1:65535", false},
		{"", true},
		{"8080", true},        // missing container port
		{"8080:80:tcp", true}, // too many parts
		{"0:80", true},        // invalid host port
		{"8080:0", true},      // invalid container port
		{"65536:80", true},    // host port too high
		{"8080:65536", true},  // container port too high
		{"abc:80", true},      // non-numeric host port
		{"8080:abc", true},    // non-numeric container port
		{":80", true},         // empty host port
		{"8080:", true},       // empty container port
	}

	for _, tt := range tests {
		t.Run(tt.mapping, func(t *testing.T) {
			err := ValidatePortMapping(tt.mapping)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateVolumeMapping(t *testing.T) {
	tests := []struct {
		mapping   string
		shouldErr bool
	}{
		{"/host/path:/container/path", false},
		{"/host/path:/container/path:ro", false},
		{"/host/path:/container/path:rw", false},
		{"/host/path:/container/path:z", false},
		{"/host/path:/container/path:ro,z", false},
		{"", true},
		{"/host/path", true},                          // missing container path
		{":/container/path", true},                    // empty host path
		{"/host/path:", true},                         // empty container path
		{"/host/path:relative/path", true},            // container path not absolute
		{"/host/path:/container/path:invalid", true},  // invalid option
		{"/host/path:/container/path:ro:extra", true}, // too many parts
	}

	for _, tt := range tests {
		t.Run(tt.mapping, func(t *testing.T) {
			err := ValidateVolumeMapping(tt.mapping)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateEnvironmentVariable(t *testing.T) {
	tests := []struct {
		envVar    string
		shouldErr bool
	}{
		{"VAR=value", false},
		{"MY_VAR=my_value", false},
		{"DEBUG=true", false},
		{"PORT=3000", false},
		{"VAR=", false},       // empty value is allowed
		{"_VAR=value", false}, // underscore at start is allowed
		{"", true},
		{"VAR", true},          // missing equals sign
		{"=value", true},       // empty key
		{"123VAR=value", true}, // key cannot start with number
		{"MY-VAR=value", true}, // hyphen not allowed in key
		{"MY VAR=value", true}, // space not allowed in key
	}

	for _, tt := range tests {
		t.Run(tt.envVar, func(t *testing.T) {
			err := ValidateEnvironmentVariable(tt.envVar)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateRestartPolicy(t *testing.T) {
	tests := []struct {
		policy    string
		shouldErr bool
	}{
		{"no", false},
		{"always", false},
		{"unless-stopped", false},
		{"on-failure", false},
		{"", true},
		{"invalid", true},
		{"Never", true},  // case sensitive
		{"ALWAYS", true}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.policy, func(t *testing.T) {
			err := ValidateRestartPolicy(tt.policy)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateHostname(t *testing.T) {
	tests := []struct {
		hostname  string
		shouldErr bool
	}{
		{"192.168.1.100", false}, // IP address
		{"example.com", false},
		{"sub.example.com", false},
		{"localhost", false},
		{"server-01", false},
		{"server01", false},
		{"", true},
		{"-server", true},                 // cannot start with hyphen
		{"server-", true},                 // cannot end with hyphen
		{"server..com", true},             // double dots not allowed
		{".example.com", true},            // cannot start with dot
		{"example.com.", true},            // cannot end with dot
		{string(make([]byte, 254)), true}, // too long
	}

	for _, tt := range tests {
		t.Run(tt.hostname, func(t *testing.T) {
			err := ValidateHostname(tt.hostname)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateSynologyPath(t *testing.T) {
	tests := []struct {
		path      string
		shouldErr bool
	}{
		{"/volume1/docker", false},
		{"/volume2/data", false},
		{"/volumeUSB1/backup", false},
		{"/volumeeSATA1/storage", false},
		{"", true},
		{"relative/path", true}, // must be absolute
		{"/invalid/path", true}, // not under valid volume
		{"/usr/local", true},    // system path not allowed
		{"/home/user", true},    // user home not allowed
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			err := ValidateSynologyPath(tt.path)
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
