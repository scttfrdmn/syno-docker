package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestConnectionToChubChub tests direct connection to your local Synology system
func TestConnectionToChubChub(t *testing.T) {
	// Skip in CI environments (GitHub Actions can't reach local networks)
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping local network test in CI environment")
	}
	// Test hostname resolution
	t.Run("HostnameResolution", func(t *testing.T) {
		cmd := exec.Command("ping", "-c", "1", "chubchub.local")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to ping chubchub.local: %v\nOutput: %s", err, output)
		}
		t.Logf("✅ Hostname chubchub.local resolves successfully")
	})

	// Test SSH connectivity
	t.Run("SSHConnectivity", func(t *testing.T) {
		cmd := exec.Command("ssh", "-o", "ConnectTimeout=10", "-o", "BatchMode=yes",
			"scttfrdmn@chubchub.local", "echo 'SSH connection successful'; whoami")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("SSH connection failed: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "SSH connection successful") {
			t.Fatalf("SSH connection test failed. Output: %s", output)
		}

		if !strings.Contains(string(output), "scttfrdmn") {
			t.Fatalf("Wrong user returned. Output: %s", output)
		}

		t.Logf("✅ SSH connection to scttfrdmn@chubchub.local successful")
	})

	// Test Docker availability
	t.Run("DockerAvailability", func(t *testing.T) {
		cmd := exec.Command("ssh", "scttfrdmn@chubchub.local", "/usr/local/bin/docker version --format 'Server: {{.Server.Version}}'")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Docker version command failed: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "Server:") {
			t.Fatalf("Invalid Docker version output: %s", output)
		}

		t.Logf("✅ Docker accessible via SSH: %s", strings.TrimSpace(string(output)))
	})

	// Test Docker container listing
	t.Run("DockerContainerList", func(t *testing.T) {
		cmd := exec.Command("ssh", "scttfrdmn@chubchub.local", "/usr/local/bin/docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}'")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Docker ps command failed: %v\nOutput: %s", err, output)
		}

		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(lines) < 1 {
			t.Fatal("No output from docker ps command")
		}

		// First line should be header
		if !strings.Contains(lines[0], "NAMES") {
			t.Fatalf("Invalid docker ps output format: %s", lines[0])
		}

		t.Logf("✅ Docker container listing successful")
		t.Logf("Current containers on chubchub.local:")
		for _, line := range lines {
			t.Logf("  %s", line)
		}
	})

	// Test Container Manager service status
	t.Run("ContainerManagerStatus", func(t *testing.T) {
		cmd := exec.Command("ssh", "scttfrdmn@chubchub.local", "systemctl is-active pkg-ContainerManager-dockerd")
		output, err := cmd.CombinedOutput()

		status := strings.TrimSpace(string(output))
		if err != nil || status != "active" {
			t.Logf("Warning: Container Manager service status: %s (error: %v)", status, err)
			// Don't fail - service might be running but not via systemctl
		} else {
			t.Logf("✅ Container Manager service is active")
		}
	})

	// Test volume path access
	t.Run("VolumePathAccess", func(t *testing.T) {
		testPath := "/volume1/docker/syno-docker-test"

		// Create test directory
		cmd := exec.Command("ssh", "scttfrdmn@chubchub.local", fmt.Sprintf("mkdir -p %s && echo 'Directory created successfully'", testPath))
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to create test directory: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "Directory created successfully") {
			t.Fatalf("Unexpected output from directory creation: %s", output)
		}

		// Test write permission
		testFile := filepath.Join(testPath, "test.txt")
		cmd = exec.Command("ssh", "scttfrdmn@chubchub.local", fmt.Sprintf("echo 'test content' > %s && cat %s", testFile, testFile))
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to write test file: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "test content") {
			t.Fatalf("File write/read test failed: %s", output)
		}

		// Cleanup test file
		cmd = exec.Command("ssh", "scttfrdmn@chubchub.local", fmt.Sprintf("rm -f %s", testFile))
		cmd.Run()

		t.Logf("✅ Volume path %s accessible and writable", testPath)
	})
}

// TestDirectDockerCommands tests Docker commands directly via SSH
func TestDirectDockerCommands(t *testing.T) {
	// Skip in CI environments
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping local network test in CI environment")
	}

	// Test Docker info command
	t.Run("DockerInfo", func(t *testing.T) {
		cmd := exec.Command("ssh", "scttfrdmn@chubchub.local",
			"/usr/local/bin/docker info --format '{{.ServerVersion}}'")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Docker info failed: %v\nOutput: %s", err, output)
		}

		version := strings.TrimSpace(string(output))
		if version == "" {
			t.Fatal("Docker info returned empty version")
		}

		t.Logf("✅ Docker server version: %s", version)
	})

	// Test image pulling (using a small image)
	t.Run("ImagePull", func(t *testing.T) {
		cmd := exec.Command("ssh", "scttfrdmn@chubchub.local",
			"/usr/local/bin/docker pull hello-world:latest")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Docker pull failed: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "Pull complete") && !strings.Contains(string(output), "Image is up to date") {
			t.Fatalf("Unexpected pull output: %s", output)
		}

		t.Logf("✅ Docker image pull successful")
	})

	// Test container run (quick test)
	t.Run("ContainerRun", func(t *testing.T) {
		containerName := fmt.Sprintf("syno-docker-test-%d", os.Getpid())

		// Cleanup any existing container with same name
		cleanupCmd := exec.Command("ssh", "scttfrdmn@chubchub.local",
			fmt.Sprintf("/usr/local/bin/docker rm -f %s 2>/dev/null || true", containerName))
		cleanupCmd.Run()

		// Run hello-world container
		cmd := exec.Command("ssh", "scttfrdmn@chubchub.local",
			fmt.Sprintf("/usr/local/bin/docker run --name %s hello-world:latest", containerName))
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Docker run failed: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "Hello from Docker!") {
			t.Fatalf("Hello-world container didn't run correctly: %s", output)
		}

		// Cleanup
		cleanupCmd = exec.Command("ssh", "scttfrdmn@chubchub.local",
			fmt.Sprintf("/usr/local/bin/docker rm -f %s", containerName))
		cleanupCmd.Run()

		t.Logf("✅ Docker container run/remove successful")
	})
}
