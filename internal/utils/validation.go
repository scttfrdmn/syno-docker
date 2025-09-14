package utils

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var (
	// Valid Docker image name pattern
	dockerImagePattern = regexp.MustCompile(`^([a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*(/[a-z0-9]+((\.|_|__|-+)[a-z0-9]+)*)*)?(:[a-zA-Z0-9_][a-zA-Z0-9._-]*)?$`)

	// Valid container name pattern
	containerNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)
)

// ValidateDockerImage validates a Docker image name format
func ValidateDockerImage(image string) error {
	if image == "" {
		return fmt.Errorf("image name cannot be empty")
	}

	if !dockerImagePattern.MatchString(image) {
		return fmt.Errorf("invalid Docker image name: %s", image)
	}

	return nil
}

// ValidateContainerName validates a Docker container name format
func ValidateContainerName(name string) error {
	if name == "" {
		return fmt.Errorf("container name cannot be empty")
	}

	if len(name) > 253 {
		return fmt.Errorf("container name too long (max 253 characters): %d", len(name))
	}

	if !containerNamePattern.MatchString(name) {
		return fmt.Errorf("invalid container name: %s (must start with alphanumeric, contain only alphanumeric, underscore, period, or hyphen)", name)
	}

	return nil
}

// ValidatePortMapping validates a port mapping string (host:container format)
func ValidatePortMapping(portMapping string) error {
	parts := strings.Split(portMapping, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid port mapping format: %s (expected host:container)", portMapping)
	}

	hostPort := strings.TrimSpace(parts[0])
	containerPort := strings.TrimSpace(parts[1])

	// Validate host port
	if err := validatePort(hostPort); err != nil {
		return fmt.Errorf("invalid host port in mapping %s: %w", portMapping, err)
	}

	// Validate container port
	if err := validatePort(containerPort); err != nil {
		return fmt.Errorf("invalid container port in mapping %s: %w", portMapping, err)
	}

	return nil
}

func validatePort(port string) error {
	if port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("port must be a number: %s", port)
	}

	if portNum < 1 || portNum > 65535 {
		return fmt.Errorf("port must be between 1 and 65535: %d", portNum)
	}

	return nil
}

// ValidateVolumeMapping validates a volume mapping string (host:container format)
func ValidateVolumeMapping(volumeMapping string) error {
	parts := strings.Split(volumeMapping, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return fmt.Errorf("invalid volume mapping format: %s (expected host:container[:options])", volumeMapping)
	}

	hostPath := strings.TrimSpace(parts[0])
	containerPath := strings.TrimSpace(parts[1])

	if hostPath == "" {
		return fmt.Errorf("host path cannot be empty in volume mapping: %s", volumeMapping)
	}

	if containerPath == "" {
		return fmt.Errorf("container path cannot be empty in volume mapping: %s", volumeMapping)
	}

	// Container path must be absolute
	if !strings.HasPrefix(containerPath, "/") {
		return fmt.Errorf("container path must be absolute in volume mapping: %s", volumeMapping)
	}

	// Validate options if present
	if len(parts) == 3 {
		options := strings.TrimSpace(parts[2])
		validOptions := map[string]bool{
			"ro":         true,
			"rw":         true,
			"z":          true,
			"Z":          true,
			"consistent": true,
			"delegated":  true,
			"cached":     true,
		}

		optionParts := strings.Split(options, ",")
		for _, option := range optionParts {
			option = strings.TrimSpace(option)
			if !validOptions[option] {
				return fmt.Errorf("invalid volume option '%s' in mapping: %s", option, volumeMapping)
			}
		}
	}

	return nil
}

// ValidateEnvironmentVariable validates an environment variable (KEY=value format)
func ValidateEnvironmentVariable(envVar string) error {
	if envVar == "" {
		return fmt.Errorf("environment variable cannot be empty")
	}

	parts := strings.SplitN(envVar, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid environment variable format: %s (expected KEY=value)", envVar)
	}

	key := strings.TrimSpace(parts[0])
	if key == "" {
		return fmt.Errorf("environment variable key cannot be empty: %s", envVar)
	}

	// Check for valid environment variable key pattern
	if !regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`).MatchString(key) {
		return fmt.Errorf("invalid environment variable key: %s (must start with letter or underscore, contain only alphanumeric and underscore)", key)
	}

	return nil
}

// ValidateRestartPolicy validates a Docker restart policy
func ValidateRestartPolicy(policy string) error {
	validPolicies := map[string]bool{
		"no":             true,
		"always":         true,
		"unless-stopped": true,
		"on-failure":     true,
	}

	if !validPolicies[policy] {
		return fmt.Errorf("invalid restart policy: %s (valid options: no, always, unless-stopped, on-failure)", policy)
	}

	return nil
}

// ValidateHostname validates a hostname or IP address format
func ValidateHostname(hostname string) error {
	if hostname == "" {
		return fmt.Errorf("hostname cannot be empty")
	}

	// Try to parse as IP address first
	if ip := net.ParseIP(hostname); ip != nil {
		return nil
	}

	// Validate as hostname
	if len(hostname) > 253 {
		return fmt.Errorf("hostname too long (max 253 characters): %d", len(hostname))
	}

	// Check for valid hostname pattern
	hostnamePattern := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	if !hostnamePattern.MatchString(hostname) {
		return fmt.Errorf("invalid hostname format: %s", hostname)
	}

	return nil
}

// ValidateSynologyPath validates that a path is under a valid Synology volume
func ValidateSynologyPath(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("path must be absolute: %s", path)
	}

	// Check if it's under a valid Synology volume
	validPrefixes := []string{
		"/volume1/",
		"/volume2/",
		"/volume3/",
		"/volume4/",
		"/volumeUSB1/",
		"/volumeUSB2/",
		"/volumeeSATA1/",
		"/volumeeSATA2/",
	}

	for _, prefix := range validPrefixes {
		if strings.HasPrefix(path, prefix) {
			return nil
		}
	}

	return fmt.Errorf("path must be under a valid Synology volume (/volume1/, /volume2/, etc.): %s", path)
}
