package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/scttfrdmn/synodeploy/pkg/synology"
)

// ComposeService represents a service in a docker-compose file
type ComposeService struct {
	Image       string            `yaml:"image,omitempty"`
	Build       interface{}       `yaml:"build,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Environment interface{}       `yaml:"environment,omitempty"`
	Restart     string            `yaml:"restart,omitempty"`
	Networks    interface{}       `yaml:"networks,omitempty"`
	DependsOn   interface{}       `yaml:"depends_on,omitempty"`
	Command     interface{}       `yaml:"command,omitempty"`
	WorkingDir  string            `yaml:"working_dir,omitempty"`
	User        string            `yaml:"user,omitempty"`
	Hostname    string            `yaml:"hostname,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
}

// ComposeFile represents a complete docker-compose file structure
type ComposeFile struct {
	Version  string                    `yaml:"version,omitempty"`
	Services map[string]ComposeService `yaml:"services"`
	Networks map[string]interface{}    `yaml:"networks,omitempty"`
	Volumes  map[string]interface{}    `yaml:"volumes,omitempty"`
}

// ComposeOptions represents options for deploying a compose file
type ComposeOptions struct {
	ComposeFile string
	ProjectName string
	EnvFile     string
}

// Compose deploys a docker-compose file to the Synology NAS
func Compose(conn *synology.Connection, opts *ComposeOptions) error {
	// Read and parse compose file
	composeData, err := parseComposeFile(opts.ComposeFile)
	if err != nil {
		return errors.Wrap(err, "failed to parse compose file")
	}

	// Load environment variables if env file specified
	envVars := make(map[string]string)
	if opts.EnvFile != "" {
		envVars, err = loadEnvFile(opts.EnvFile)
		if err != nil {
			return errors.Wrap(err, "failed to load environment file")
		}
	}

	// Deploy each service as a container
	fmt.Printf("Deploying compose project: %s\n", opts.ProjectName)

	for serviceName, service := range composeData.Services {
		containerName := fmt.Sprintf("%s_%s_1", opts.ProjectName, serviceName)

		fmt.Printf("Deploying service: %s (container: %s)\n", serviceName, containerName)

		// Convert compose service to container options
		containerOpts, err := convertServiceToContainer(service, containerName, envVars)
		if err != nil {
			return errors.Wrapf(err, "failed to convert service %s to container options", serviceName)
		}

		// Deploy container
		_, err = Container(conn, containerOpts)
		if err != nil {
			return errors.Wrapf(err, "failed to deploy service %s", serviceName)
		}
	}

	fmt.Printf("✅ Compose project %s deployed successfully!\n", opts.ProjectName)
	return nil
}

func parseComposeFile(composePath string) (*ComposeFile, error) {
	// Check if file exists
	if _, err := os.Stat(composePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("compose file not found: %s", composePath)
	}

	// Read file
	data, err := os.ReadFile(composePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read compose file: %w", err)
	}

	// Parse YAML
	var compose ComposeFile
	if err := yaml.Unmarshal(data, &compose); err != nil {
		return nil, fmt.Errorf("failed to parse compose YAML: %w", err)
	}

	return &compose, nil
}

func loadEnvFile(envPath string) (map[string]string, error) {
	envVars := make(map[string]string)

	data, err := os.ReadFile(envPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read env file: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Remove quotes if present
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}

			envVars[key] = value
		}
	}

	return envVars, nil
}

func convertServiceToContainer(service ComposeService, containerName string, envVars map[string]string) (*ContainerOptions, error) {
	opts := &ContainerOptions{
		Image:       service.Image,
		Name:        containerName,
		Ports:       service.Ports,
		Volumes:     service.Volumes,
		Restart:     service.Restart,
		User:        service.User,
		WorkingDir:  service.WorkingDir,
		NetworkMode: synology.DefaultNetwork,
	}

	// Set default restart policy if not specified
	if opts.Restart == "" {
		opts.Restart = synology.DefaultRestartPolicy
	}

	// Process environment variables
	env, err := processEnvironment(service.Environment, envVars)
	if err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}
	opts.Env = env

	// Process command if specified
	if service.Command != nil {
		cmd, err := processCommand(service.Command)
		if err != nil {
			return nil, fmt.Errorf("failed to process command: %w", err)
		}
		opts.Command = cmd
	}

	return opts, nil
}

func processEnvironment(env interface{}, envVars map[string]string) ([]string, error) {
	var result []string

	switch e := env.(type) {
	case []interface{}:
		// Array format: ["KEY=value", "KEY2=value2"]
		for _, item := range e {
			if str, ok := item.(string); ok {
				result = append(result, expandEnvVar(str, envVars))
			}
		}
	case map[string]interface{}:
		// Object format: {KEY: value, KEY2: value2}
		for key, value := range e {
			if str, ok := value.(string); ok {
				result = append(result, fmt.Sprintf("%s=%s", key, expandEnvVar(str, envVars)))
			}
		}
	case nil:
		// No environment variables
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported environment format: %T", env)
	}

	return result, nil
}

func processCommand(cmd interface{}) ([]string, error) {
	switch c := cmd.(type) {
	case string:
		// Single command string
		return []string{c}, nil
	case []interface{}:
		// Array of command parts
		var result []string
		for _, part := range c {
			if str, ok := part.(string); ok {
				result = append(result, str)
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported command format: %T", cmd)
	}
}

func expandEnvVar(value string, envVars map[string]string) string {
	// Simple environment variable expansion
	// Support ${VAR} and $VAR formats
	result := value

	for key, val := range envVars {
		result = strings.ReplaceAll(result, "${"+key+"}", val)
		result = strings.ReplaceAll(result, "$"+key, val)
	}

	return result
}

// GenerateProjectName generates a project name from a compose file path
func GenerateProjectName(composePath string) string {
	// Use directory name as project name
	dir := filepath.Dir(composePath)
	projectName := filepath.Base(dir)

	// Clean up project name (remove special characters)
	projectName = strings.ToLower(projectName)
	projectName = strings.ReplaceAll(projectName, " ", "")
	projectName = strings.ReplaceAll(projectName, "_", "")
	projectName = strings.ReplaceAll(projectName, "-", "")

	if projectName == "" || projectName == "." || projectName == "/" {
		projectName = "synodeploy"
	}

	return projectName
}
