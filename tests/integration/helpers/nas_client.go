package helpers

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/scttfrdmn/synodeploy/pkg/deploy"
	"github.com/scttfrdmn/synodeploy/pkg/synology"
)

// RandomString generates a random string of specified length
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// CreateTestDirectory creates a directory on the NAS for testing
func CreateTestDirectory(conn *synology.Connection, path string) error {
	cmd := fmt.Sprintf("mkdir -p %s && chmod 755 %s", path, path)
	_, err := conn.ExecuteCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

// CreateTestFile creates a file with specified content on the NAS
func CreateTestFile(conn *synology.Connection, filePath, content string) error {
	// Ensure parent directory exists
	dir := filePath[:strings.LastIndex(filePath, "/")]
	if err := CreateTestDirectory(conn, dir); err != nil {
		return err
	}

	// Create file with content
	cmd := fmt.Sprintf("cat > %s << 'EOF'\n%s\nEOF", filePath, content)
	_, err := conn.ExecuteCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	return nil
}

// TestHTTPEndpoint tests if an HTTP endpoint is accessible and returns expected content
func TestHTTPEndpoint(url, expectedContent string, timeout time.Duration) error {
	client := &http.Client{Timeout: timeout}

	// Retry logic for container startup
	maxRetries := 10
	retryInterval := 3 * time.Second

	for i := 0; i < maxRetries; i++ {
		resp, err := client.Get(url)
		if err != nil {
			if i == maxRetries-1 {
				return fmt.Errorf("failed to connect to %s after %d retries: %w", url, maxRetries, err)
			}
			time.Sleep(retryInterval)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			if i == maxRetries-1 {
				return fmt.Errorf("HTTP %s returned status %d", url, resp.StatusCode)
			}
			time.Sleep(retryInterval)
			continue
		}

		// Read response body
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		content := string(body[:n])

		if !strings.Contains(content, expectedContent) {
			return fmt.Errorf("response from %s does not contain expected content '%s', got: %s",
				url, expectedContent, content)
		}

		return nil // Success
	}

	return fmt.Errorf("failed to get successful response from %s after %d retries", url, maxRetries)
}

// WaitForContainer waits for a container to reach the expected state
func WaitForContainer(conn *synology.Connection, containerName, expectedState string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		containers, err := deploy.ListContainers(conn, true)
		if err != nil {
			return fmt.Errorf("failed to list containers: %w", err)
		}

		for _, container := range containers {
			if container.Name == containerName {
				if strings.Contains(container.Status, expectedState) {
					return nil
				}
				break
			}
		}

		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("container %s did not reach state '%s' within %v", containerName, expectedState, timeout)
}

// CheckContainerLogs retrieves and checks container logs
func CheckContainerLogs(conn *synology.Connection, containerName, expectedLogContent string) error {
	output, err := conn.ExecuteDockerCommand([]string{"logs", containerName})
	if err != nil {
		return fmt.Errorf("failed to get logs for container %s: %w", containerName, err)
	}

	if !strings.Contains(output, expectedLogContent) {
		return fmt.Errorf("container %s logs do not contain expected content '%s', got: %s",
			containerName, expectedLogContent, output)
	}

	return nil
}

// GetContainerIP retrieves the IP address of a container
func GetContainerIP(conn *synology.Connection, containerName string) (string, error) {
	output, err := conn.ExecuteDockerCommand([]string{
		"inspect", "--format", "'{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}'", containerName,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get IP for container %s: %w", containerName, err)
	}

	ip := strings.Trim(strings.TrimSpace(output), "'")
	if ip == "" {
		return "", fmt.Errorf("no IP address found for container %s", containerName)
	}

	return ip, nil
}

// TestContainerConnectivity tests network connectivity between containers
func TestContainerConnectivity(conn *synology.Connection, fromContainer, toContainer string, port int) error {
	toIP, err := GetContainerIP(conn, toContainer)
	if err != nil {
		return err
	}

	// Test connectivity using nc (netcat)
	cmd := []string{"exec", fromContainer, "nc", "-z", toIP, fmt.Sprintf("%d", port)}
	_, err = conn.ExecuteDockerCommand(cmd)
	if err != nil {
		return fmt.Errorf("connectivity test from %s to %s:%d failed: %w", fromContainer, toContainer, port, err)
	}

	return nil
}