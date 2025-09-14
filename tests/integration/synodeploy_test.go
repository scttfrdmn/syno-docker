package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/scttfrdmn/synodeploy/pkg/config"
	"github.com/scttfrdmn/synodeploy/pkg/deploy"
	"github.com/scttfrdmn/synodeploy/pkg/synology"
	"github.com/scttfrdmn/synodeploy/tests/integration/helpers"
)

// TestSynoDeployEndToEnd tests the complete SynoDeploy workflow
func TestSynoDeployEndToEnd(t *testing.T) {
	// Setup configuration for chubchub.local
	homeDir, _ := os.UserHomeDir()
	sshKeyPath := filepath.Join(homeDir, ".ssh", "id_rsa")

	cfg := &config.Config{
		Host:       "chubchub.local",
		User:       "scttfrdmn",
		Port:       22,
		SSHKeyPath: sshKeyPath,
		Defaults: struct {
			VolumePath string `yaml:"volume_path"`
			Network    string `yaml:"network,omitempty"`
		}{
			VolumePath: "/volume1/docker",
			Network:    "bridge",
		},
	}

	// Test connection
	t.Run("Connection", func(t *testing.T) {
		conn := synology.NewConnection(cfg)
		if err := conn.Connect(); err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		if err := conn.TestConnection(); err != nil {
			t.Fatalf("Connection test failed: %v", err)
		}

		t.Logf("✅ SynoDeploy connection to chubchub.local successful")
	})

	// Test single container deployment
	t.Run("SingleContainerDeployment", func(t *testing.T) {
		conn := synology.NewConnection(cfg)
		if err := conn.Connect(); err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		// Deploy nginx container
		containerName := "synodeploy-test-nginx-" + helpers.RandomString(6)
		opts := &deploy.ContainerOptions{
			Image:       "nginx:alpine",
			Name:        containerName,
			Ports:       []string{"8082:80"},
			Volumes:     []string{"/volume1/docker/test-html:/usr/share/nginx/html"},
			Restart:     "unless-stopped",
			NetworkMode: "bridge",
		}

		// Create test HTML content
		htmlDir := "/volume1/docker/test-html"
		if err := helpers.CreateTestDirectory(conn, htmlDir); err != nil {
			t.Fatalf("Failed to create HTML directory: %v", err)
		}

		htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html><head><title>SynoDeploy Test</title></head>
<body><h1>SynoDeploy Integration Test</h1>
<p>Container: %s</p>
<p>Timestamp: %s</p></body></html>`,
			containerName, time.Now().Format("2006-01-02 15:04:05"))

		if err := helpers.CreateTestFile(conn, filepath.Join(htmlDir, "index.html"), htmlContent); err != nil {
			t.Fatalf("Failed to create test HTML file: %v", err)
		}

		// Deploy container
		containerID, err := deploy.DeployContainer(conn, opts)
		if err != nil {
			t.Fatalf("Failed to deploy container: %v", err)
		}

		t.Logf("✅ Container deployed: %s (ID: %s)", containerName, containerID)

		// Wait for container to start
		time.Sleep(5 * time.Second)

		// Verify container is running
		containers, err := deploy.ListContainers(conn, false)
		if err != nil {
			t.Fatalf("Failed to list containers: %v", err)
		}

		found := false
		for _, container := range containers {
			if container.Name == containerName {
				found = true
				if !strings.Contains(container.Status, "Up") {
					t.Errorf("Container %s is not running. Status: %s", containerName, container.Status)
				}
				t.Logf("✅ Container status: %s", container.Status)
				break
			}
		}

		if !found {
			t.Errorf("Container %s not found in running containers", containerName)
		}

		// Test HTTP endpoint
		if err := helpers.TestHTTPEndpoint("http://chubchub.local:8082", "SynoDeploy Integration Test", 30*time.Second); err != nil {
			t.Errorf("HTTP test failed: %v", err)
		} else {
			t.Logf("✅ HTTP endpoint accessible and serving content")
		}

		// Cleanup
		if err := deploy.RemoveContainer(conn, containerName, true); err != nil {
			t.Errorf("Failed to cleanup container: %v", err)
		} else {
			t.Logf("✅ Container cleaned up successfully")
		}

		// Cleanup test files
		conn.ExecuteCommand(fmt.Sprintf("rm -rf %s", htmlDir))
	})

	// Test container lifecycle management
	t.Run("ContainerLifecycle", func(t *testing.T) {
		conn := synology.NewConnection(cfg)
		if err := conn.Connect(); err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		containerName := "synodeploy-lifecycle-test-" + helpers.RandomString(6)

		// Deploy a long-running container
		opts := &deploy.ContainerOptions{
			Image:       "alpine:latest",
			Name:        containerName,
			Command:     []string{"sleep", "60"},
			Restart:     "no",
			NetworkMode: "bridge",
		}

		// Deploy
		containerID, err := deploy.DeployContainer(conn, opts)
		if err != nil {
			t.Fatalf("Failed to deploy container: %v", err)
		}
		t.Logf("✅ Long-running container deployed: %s", containerID)

		// Wait for container to start
		time.Sleep(3 * time.Second)

		// Verify it's running
		containers, err := deploy.ListContainers(conn, false)
		if err != nil {
			t.Fatalf("Failed to list containers: %v", err)
		}

		found := false
		for _, container := range containers {
			if container.Name == containerName {
				found = true
				if !strings.Contains(container.Status, "Up") {
					t.Errorf("Container should be running. Status: %s", container.Status)
				}
				break
			}
		}

		if !found {
			t.Error("Container not found in running containers")
		}

		// Remove the container
		if err := deploy.RemoveContainer(conn, containerName, true); err != nil {
			t.Fatalf("Failed to remove container: %v", err)
		}
		t.Logf("✅ Container lifecycle test completed successfully")
	})

	// Test error handling
	t.Run("ErrorHandling", func(t *testing.T) {
		conn := synology.NewConnection(cfg)
		if err := conn.Connect(); err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		// Test deploying non-existent image
		opts := &deploy.ContainerOptions{
			Image:       "nonexistent/image:latest",
			Name:        "should-fail",
			NetworkMode: "bridge",
		}

		_, err := deploy.DeployContainer(conn, opts)
		if err == nil {
			t.Error("Expected deployment to fail with non-existent image")
			// Cleanup if it somehow succeeded
			deploy.RemoveContainer(conn, "should-fail", true)
		} else {
			t.Logf("✅ Correctly failed to deploy non-existent image: %v", err)
		}

		// Test removing non-existent container
		err = deploy.RemoveContainer(conn, "non-existent-container", false)
		if err == nil {
			t.Error("Expected removal of non-existent container to fail")
		} else {
			t.Logf("✅ Correctly failed to remove non-existent container")
		}
	})
}