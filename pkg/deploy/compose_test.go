package deploy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseComposeFile(t *testing.T) {
	// Create a temporary compose file
	tempDir := t.TempDir()
	composeFile := filepath.Join(tempDir, "docker-compose.yml")

	composeContent := `version: '3.8'
services:
  web:
    image: nginx:latest
    ports:
      - "8080:80"
    volumes:
      - ./html:/usr/share/nginx/html
    environment:
      - ENV=production
    restart: unless-stopped

  db:
    image: postgres:13
    environment:
      POSTGRES_DB: mydb
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    volumes:
      - db_data:/var/lib/postgresql/data

volumes:
  db_data:
`

	if err := os.WriteFile(composeFile, []byte(composeContent), 0644); err != nil {
		t.Fatalf("Failed to create compose file: %v", err)
	}

	// Parse the compose file
	compose, err := parseComposeFile(composeFile)
	if err != nil {
		t.Fatalf("Failed to parse compose file: %v", err)
	}

	// Verify parsing results
	if compose.Version != "3.8" {
		t.Errorf("Expected version 3.8, got %s", compose.Version)
	}

	if len(compose.Services) != 2 {
		t.Errorf("Expected 2 services, got %d", len(compose.Services))
	}

	// Check web service
	web, exists := compose.Services["web"]
	if !exists {
		t.Fatal("Web service not found")
	}

	if web.Image != "nginx:latest" {
		t.Errorf("Expected web image nginx:latest, got %s", web.Image)
	}

	if len(web.Ports) != 1 || web.Ports[0] != "8080:80" {
		t.Errorf("Expected web ports [8080:80], got %v", web.Ports)
	}

	if web.Restart != "unless-stopped" {
		t.Errorf("Expected web restart unless-stopped, got %s", web.Restart)
	}

	// Check db service
	db, exists := compose.Services["db"]
	if !exists {
		t.Fatal("DB service not found")
	}

	if db.Image != "postgres:13" {
		t.Errorf("Expected db image postgres:13, got %s", db.Image)
	}
}

func TestLoadEnvFile(t *testing.T) {
	// Create a temporary env file
	tempDir := t.TempDir()
	envFile := filepath.Join(tempDir, ".env")

	envContent := `# This is a comment
DATABASE_URL=postgres://localhost:5432/mydb
API_TOKEN=test-token-value
DEBUG=true
EMPTY_VAR=

# Another comment
PORT=3000`

	if err := os.WriteFile(envFile, []byte(envContent), 0644); err != nil {
		t.Fatalf("Failed to create env file: %v", err)
	}

	// Load the env file
	envVars, err := loadEnvFile(envFile)
	if err != nil {
		t.Fatalf("Failed to load env file: %v", err)
	}

	expectedVars := map[string]string{
		"DATABASE_URL": "postgres://localhost:5432/mydb",
		"API_TOKEN":    "test-token-value",
		"DEBUG":        "true",
		"EMPTY_VAR":    "",
		"PORT":         "3000",
	}

	if len(envVars) != len(expectedVars) {
		t.Errorf("Expected %d variables, got %d", len(expectedVars), len(envVars))
	}

	for key, expected := range expectedVars {
		if actual, exists := envVars[key]; !exists {
			t.Errorf("Expected variable %s not found", key)
		} else if actual != expected {
			t.Errorf("Variable %s: expected %s, got %s", key, expected, actual)
		}
	}
}

func TestProcessEnvironment(t *testing.T) {
	envVars := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	tests := []struct {
		name     string
		env      interface{}
		expected []string
		hasError bool
	}{
		{
			name:     "array format",
			env:      []interface{}{"DATABASE_URL=${DB_HOST}:${DB_PORT}/mydb", "DEBUG=true"},
			expected: []string{"DATABASE_URL=localhost:5432/mydb", "DEBUG=true"},
			hasError: false,
		},
		{
			name: "map format",
			env: map[string]interface{}{
				"DATABASE_URL": "${DB_HOST}:${DB_PORT}/mydb",
				"DEBUG":        "true",
			},
			expected: []string{"DATABASE_URL=localhost:5432/mydb", "DEBUG=true"},
			hasError: false,
		},
		{
			name:     "nil environment",
			env:      nil,
			expected: []string{},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processEnvironment(tt.env, envVars)

			if tt.hasError && err == nil {
				t.Error("Expected error but got none")
				return
			}

			if !tt.hasError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}

			if tt.hasError {
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d env vars, got %d", len(tt.expected), len(result))
				return
			}

			// Convert to map for easier comparison (order may vary)
			resultMap := make(map[string]bool)
			for _, env := range result {
				resultMap[env] = true
			}

			for _, expected := range tt.expected {
				if !resultMap[expected] {
					t.Errorf("Expected env var %s not found in result", expected)
				}
			}
		})
	}
}

func TestProcessCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      interface{}
		expected []string
		hasError bool
	}{
		{
			name:     "string command",
			cmd:      "npm start",
			expected: []string{"npm start"},
			hasError: false,
		},
		{
			name:     "array command",
			cmd:      []interface{}{"npm", "run", "dev"},
			expected: []string{"npm", "run", "dev"},
			hasError: false,
		},
		{
			name:     "empty array",
			cmd:      []interface{}{},
			expected: []string{},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processCommand(tt.cmd)

			if tt.hasError && err == nil {
				t.Error("Expected error but got none")
				return
			}

			if !tt.hasError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}

			if tt.hasError {
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d command parts, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Command part %d: expected %s, got %s", i, expected, result[i])
				}
			}
		})
	}
}

func TestExpandEnvVar(t *testing.T) {
	envVars := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
		"NAME": "myapp",
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"${HOST}", "localhost"},
		{"$HOST", "localhost"},
		{"${HOST}:${PORT}", "localhost:5432"},
		{"$HOST:$PORT", "localhost:5432"},
		{"prefix-${NAME}-suffix", "prefix-myapp-suffix"},
		{"no-variables", "no-variables"},
		{"${NONEXISTENT}", "${NONEXISTENT}"}, // Undefined variables remain unchanged
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := expandEnvVar(tt.input, envVars)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGenerateProjectName(t *testing.T) {
	tests := []struct {
		composePath string
		expected    string
	}{
		{"/home/user/myproject/docker-compose.yml", "myproject"},
		{"/path/to/my-app/docker-compose.yml", "myapp"},
		{"/path/to/My_Project/docker-compose.yml", "myproject"},
		{"./docker-compose.yml", "synodeploy"},
		{"/docker-compose.yml", "synodeploy"}, // This should generate the default name
	}

	for _, tt := range tests {
		t.Run(tt.composePath, func(t *testing.T) {
			result := GenerateProjectName(tt.composePath)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
