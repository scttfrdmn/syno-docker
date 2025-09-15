package integration

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/scttfrdmn/syno-docker/pkg/deploy"
	"github.com/scttfrdmn/syno-docker/tests/integration/helpers"
)

// TestComprehensiveCommandSuite tests all v0.2.x commands with real Synology hardware
func TestComprehensiveCommandSuite(t *testing.T) {
	if !*integrationTest {
		t.Skip("Integration tests not enabled. Use -integration flag.")
	}

	if *nasHost == "" {
		t.Skip("NAS host not specified. Use -nas-host flag.")
	}

	// Setup test environment
	runner, err := setupTestEnvironment()
	if err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}
	defer runner.Cleanup.CleanupAll()

	// Test all command phases
	t.Run("ContainerOperations", func(t *testing.T) {
		testContainerOperations(t, runner)
	})

	t.Run("ImageManagement", func(t *testing.T) {
		testImageManagement(t, runner)
	})

	t.Run("VolumeManagement", func(t *testing.T) {
		testVolumeManagement(t, runner)
	})

	t.Run("NetworkManagement", func(t *testing.T) {
		testNetworkManagement(t, runner)
	})

	t.Run("SystemOperations", func(t *testing.T) {
		testSystemOperations(t, runner)
	})
}

// testContainerOperations tests logs, exec, start/stop/restart, stats
func testContainerOperations(t *testing.T, runner *TestRunner) {
	// Deploy a test container for operations testing
	containerName := "test-ops-" + helpers.RandomString(6)
	runner.Cleanup.AddContainer(containerName)

	t.Run("DeployTestContainer", func(t *testing.T) {
		opts := deploy.NewContainerOptions("nginx:alpine")
		opts.Name = containerName
		opts.Ports = []string{"0:80"} // Use random port
		opts.Env = []string{"NGINX_TEST=integration"}

		containerID, err := deploy.Container(runner.Connection, opts)
		if err != nil {
			t.Fatalf("Failed to deploy test container: %v", err)
		}

		if containerID == "" {
			t.Fatal("Container ID is empty")
		}

		t.Logf("✅ Test container deployed: %s (ID: %s)", containerName, containerID)

		// Wait for container to be running
		if err := helpers.WaitForContainer(runner.Connection, containerName, "Up", 30*time.Second); err != nil {
			t.Fatalf("Container did not start: %v", err)
		}
	})

	t.Run("TestLogs", func(t *testing.T) {
		// Test basic log retrieval
		logs, err := deploy.GetContainerLogs(runner.Connection, containerName, "all", "", false)
		if err != nil {
			t.Fatalf("Failed to get container logs: %v", err)
		}

		if !strings.Contains(logs, "nginx") {
			t.Errorf("Expected nginx logs, got: %s", logs)
		}

		// Test logs with timestamps
		logsWithTime, err := deploy.GetContainerLogs(runner.Connection, containerName, "10", "", true)
		if err != nil {
			t.Fatalf("Failed to get logs with timestamps: %v", err)
		}

		if len(logsWithTime) == 0 {
			t.Error("Expected timestamped logs but got empty result")
		}

		t.Logf("✅ Container logs retrieved successfully")
	})

	t.Run("TestExec", func(t *testing.T) {
		// Test non-interactive exec
		opts := &deploy.ExecOptions{
			Interactive: false,
			TTY:         false,
		}

		output, err := deploy.ExecCommand(runner.Connection, containerName, []string{"whoami"}, opts)
		if err != nil {
			t.Fatalf("Failed to execute command: %v", err)
		}

		if !strings.Contains(output, "root") {
			t.Errorf("Expected 'root' in output, got: %s", output)
		}

		// Test exec with custom user
		opts.User = "nginx"
		output, err = deploy.ExecCommand(runner.Connection, containerName, []string{"whoami"}, opts)
		if err != nil {
			t.Fatalf("Failed to execute command as nginx user: %v", err)
		}

		if !strings.Contains(output, "nginx") {
			t.Errorf("Expected 'nginx' in output, got: %s", output)
		}

		t.Logf("✅ Container exec commands executed successfully")
	})

	t.Run("TestLifecycleOperations", func(t *testing.T) {
		// Test stop
		if err := deploy.StopContainer(runner.Connection, containerName, 10); err != nil {
			t.Fatalf("Failed to stop container: %v", err)
		}

		// Wait for stopped state
		if err := helpers.WaitForContainer(runner.Connection, containerName, "Exited", 30*time.Second); err != nil {
			t.Fatalf("Container did not stop: %v", err)
		}

		// Test start
		if err := deploy.StartContainer(runner.Connection, containerName); err != nil {
			t.Fatalf("Failed to start container: %v", err)
		}

		// Wait for running state
		if err := helpers.WaitForContainer(runner.Connection, containerName, "Up", 30*time.Second); err != nil {
			t.Fatalf("Container did not start: %v", err)
		}

		// Test restart
		if err := deploy.RestartContainer(runner.Connection, containerName, 10); err != nil {
			t.Fatalf("Failed to restart container: %v", err)
		}

		// Wait for running state after restart
		if err := helpers.WaitForContainer(runner.Connection, containerName, "Up", 30*time.Second); err != nil {
			t.Fatalf("Container did not restart: %v", err)
		}

		t.Logf("✅ Container lifecycle operations completed successfully")
	})

	t.Run("TestStats", func(t *testing.T) {
		// Test stats functionality by checking if we can get resource usage
		// Note: This is a simplified test since stats normally streams
		opts := &deploy.StatsOptions{
			NoStream: true,
			All:      false,
		}

		// We can't easily capture streaming output in tests, so we'll just verify the command executes
		err := deploy.ShowContainerStats(runner.Connection, []string{containerName}, opts)
		if err != nil {
			t.Fatalf("Failed to get container stats: %v", err)
		}

		t.Logf("✅ Container stats command executed successfully")
	})
}

// testImageManagement tests pull, images, rmi, export/import
func testImageManagement(t *testing.T, runner *TestRunner) {
	testImage := "alpine:3.18"

	t.Run("TestPull", func(t *testing.T) {
		opts := &deploy.PullOptions{
			Quiet: true,
		}

		if err := deploy.PullImage(runner.Connection, testImage, opts); err != nil {
			t.Fatalf("Failed to pull image: %v", err)
		}

		t.Logf("✅ Image %s pulled successfully", testImage)
	})

	t.Run("TestImages", func(t *testing.T) {
		opts := &deploy.ImagesOptions{
			All: true,
		}

		images, err := deploy.ListImages(runner.Connection, "", opts)
		if err != nil {
			t.Fatalf("Failed to list images: %v", err)
		}

		found := false
		for _, img := range images {
			if strings.Contains(img.Repository, "alpine") && strings.Contains(img.Tag, "3.18") {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected to find alpine:3.18 in image list")
		}

		// Test quiet mode
		opts.Quiet = true
		imageIDs, err := deploy.ListImageIDs(runner.Connection, "", opts)
		if err != nil {
			t.Fatalf("Failed to list image IDs: %v", err)
		}

		if len(imageIDs) == 0 {
			t.Error("Expected to find at least one image ID")
		}

		t.Logf("✅ Image listing completed successfully, found %d images", len(images))
	})

	t.Run("TestExportImport", func(t *testing.T) {
		// Create a small test container first
		testContainerName := "test-export-" + helpers.RandomString(6)
		runner.Cleanup.AddContainer(testContainerName)

		opts := deploy.NewContainerOptions("alpine:3.18")
		opts.Name = testContainerName
		opts.Command = []string{"sh", "-c", "echo 'test data' > /tmp/test.txt && sleep 3"}

		containerID, err := deploy.Container(runner.Connection, opts)
		if err != nil {
			t.Fatalf("Failed to deploy container for export test: %v", err)
		}

		// Wait for container to be running
		if err := helpers.WaitForContainer(runner.Connection, testContainerName, "Up", 60*time.Second); err != nil {
			t.Fatalf("Container did not start: %v", err)
		}

		// Give the container a moment to execute the command
		time.Sleep(5 * time.Second)

		// Test export (to stdout - simplified test)
		exportOpts := &deploy.ExportOptions{}
		if err := deploy.ExportContainer(runner.Connection, testContainerName, exportOpts); err != nil {
			t.Fatalf("Failed to export container: %v", err)
		}

		shortID := containerID
		if len(containerID) > 12 {
			shortID = containerID[:12]
		}
		t.Logf("✅ Container %s exported successfully (ID: %s)", testContainerName, shortID)
	})

	t.Run("TestRemoveImage", func(t *testing.T) {
		// Remove the test image we pulled
		opts := &deploy.RmiOptions{
			Force: true,
		}

		if err := deploy.RemoveImage(runner.Connection, testImage, opts); err != nil {
			t.Fatalf("Failed to remove image: %v", err)
		}

		// Verify image is removed
		images, err := deploy.ListImages(runner.Connection, "alpine", &deploy.ImagesOptions{})
		if err != nil {
			t.Fatalf("Failed to list images after removal: %v", err)
		}

		for _, img := range images {
			if img.Repository == "alpine" && img.Tag == "3.18" {
				t.Error("Image should have been removed but still exists")
			}
		}

		t.Logf("✅ Image %s removed successfully", testImage)
	})
}

// testVolumeManagement tests volume lifecycle operations
func testVolumeManagement(t *testing.T, runner *TestRunner) {
	volumeName := "test-vol-" + helpers.RandomString(6)
	runner.Cleanup.AddVolume(volumeName)

	// Create volume once for all tests at the function level
	opts := &deploy.VolumeCreateOptions{
		Driver: "local",
		Labels: []string{"test=integration", "phase=v0.2.x"},
	}

	createdName, err := deploy.CreateVolume(runner.Connection, volumeName, opts)
	if err != nil {
		t.Fatalf("Failed to create volume for tests: %v", err)
	}

	if createdName != volumeName {
		t.Errorf("Expected volume name %s, got %s", volumeName, createdName)
	}

	t.Logf("✅ Volume %s created successfully for testing", volumeName)

	t.Run("TestVolumeCreate", func(t *testing.T) {
		// This test now just verifies the creation worked
		t.Logf("✅ Volume creation already verified at function level")
	})

	t.Run("TestVolumeList", func(t *testing.T) {
		opts := &deploy.VolumeListOptions{}

		volumes, err := deploy.ListVolumes(runner.Connection, opts)
		if err != nil {
			t.Fatalf("Failed to list volumes: %v", err)
		}

		found := false
		for _, vol := range volumes {
			if vol.Name == volumeName {
				found = true
				if vol.Driver != "local" {
					t.Errorf("Expected driver 'local', got '%s'", vol.Driver)
				}
				break
			}
		}

		if !found {
			t.Errorf("Expected to find volume %s in list", volumeName)
		}

		// Test quiet mode
		opts.Quiet = true
		volumeNames, err := deploy.ListVolumeNames(runner.Connection, opts)
		if err != nil {
			t.Fatalf("Failed to list volume names: %v", err)
		}

		found = false
		for _, name := range volumeNames {
			if name == volumeName {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected to find volume %s in quiet list", volumeName)
		}

		t.Logf("✅ Volume listing completed successfully, found %d volumes", len(volumes))
	})

	t.Run("TestVolumeInspect", func(t *testing.T) {
		opts := &deploy.VolumeInspectOptions{}

		info, err := deploy.InspectVolume(runner.Connection, volumeName, opts)
		if err != nil {
			t.Fatalf("Failed to inspect volume: %v", err)
		}

		if !strings.Contains(info, volumeName) {
			t.Errorf("Expected volume info to contain %s, got: %s", volumeName, info)
		}

		// Test with format
		opts.Format = "'{{.Name}}'"
		formattedInfo, err := deploy.InspectVolume(runner.Connection, volumeName, opts)
		if err != nil {
			t.Fatalf("Failed to inspect volume with format: %v", err)
		}

		if !strings.Contains(formattedInfo, volumeName) {
			t.Errorf("Expected formatted info to contain %s, got: %s", volumeName, formattedInfo)
		}

		t.Logf("✅ Volume inspection completed successfully")
	})

	t.Run("TestVolumeUsage", func(t *testing.T) {
		// Verify the volume exists first
		exists, err := helpers.VerifyVolumeExists(runner.Connection, volumeName)
		if err != nil {
			t.Fatalf("Failed to check if volume exists: %v", err)
		}
		if !exists {
			t.Fatalf("Volume %s does not exist for usage test", volumeName)
		}

		// Test volume usage by using our Container function with simpler command
		containerName := "test-vol-usage-" + helpers.RandomString(6)
		runner.Cleanup.AddContainer(containerName)

		volOpts := deploy.NewContainerOptions("alpine:latest")
		volOpts.Name = containerName
		volOpts.Volumes = []string{fmt.Sprintf("%s:/data", volumeName)}
		volOpts.Command = []string{"sh", "-c", "df /data && echo 'Volume mount verified' && sleep 10"}

		containerID, err := deploy.Container(runner.Connection, volOpts)
		if err != nil {
			t.Fatalf("Failed to deploy container with volume: %v", err)
		}

		// Wait for container to complete command
		time.Sleep(5 * time.Second)

		// Get container logs to verify volume mount worked
		logs, err := deploy.GetContainerLogs(runner.Connection, containerName, "all", "", false)
		if err != nil {
			t.Fatalf("Failed to get container logs: %v", err)
		}

		// Check that volume mount was successful
		if !strings.Contains(logs, "Volume mount verified") {
			t.Errorf("Expected 'Volume mount verified' in logs, got: %s", logs)
		}

		// Check that df shows the mount point
		if !strings.Contains(logs, "/data") {
			t.Errorf("Expected '/data' mount point in df output, got: %s", logs)
		}

		shortID := containerID
		if len(containerID) > 12 {
			shortID = containerID[:12]
		}
		t.Logf("✅ Volume usage verified successfully (container: %s)", shortID)
	})

	t.Run("TestVolumeRemove", func(t *testing.T) {
		opts := &deploy.VolumeRemoveOptions{
			Force: true,
		}

		if err := deploy.RemoveVolume(runner.Connection, volumeName, opts); err != nil {
			t.Fatalf("Failed to remove volume: %v", err)
		}

		// Verify volume is removed
		volumes, err := deploy.ListVolumes(runner.Connection, &deploy.VolumeListOptions{})
		if err != nil {
			t.Fatalf("Failed to list volumes after removal: %v", err)
		}

		for _, vol := range volumes {
			if vol.Name == volumeName {
				t.Error("Volume should have been removed but still exists")
			}
		}

		t.Logf("✅ Volume %s removed successfully", volumeName)
	})
}

// testNetworkManagement tests network operations
func testNetworkManagement(t *testing.T, runner *TestRunner) {
	networkName := "test-net-" + helpers.RandomString(6)
	runner.Cleanup.AddNetwork(networkName)

	t.Run("TestNetworkCreate", func(t *testing.T) {
		opts := &deploy.NetworkCreateOptions{
			Driver:   "bridge",
			Subnet:   []string{"172.25.0.0/16"},
			Gateway:  []string{"172.25.0.1"},
			Labels:   []string{"test=integration", "phase=network"},
			Internal: false,
		}

		networkID, err := deploy.CreateNetwork(runner.Connection, networkName, opts)
		if err != nil {
			t.Fatalf("Failed to create network: %v", err)
		}

		if networkID == "" {
			t.Fatal("Network ID is empty")
		}

		t.Logf("✅ Network %s created successfully (ID: %s)", networkName, networkID)
	})

	t.Run("TestNetworkList", func(t *testing.T) {
		opts := &deploy.NetworkListOptions{}

		networks, err := deploy.ListNetworks(runner.Connection, opts)
		if err != nil {
			t.Fatalf("Failed to list networks: %v", err)
		}

		found := false
		for _, net := range networks {
			if net.Name == networkName {
				found = true
				if net.Driver != "bridge" {
					t.Errorf("Expected driver 'bridge', got '%s'", net.Driver)
				}
				break
			}
		}

		if !found {
			t.Errorf("Expected to find network %s in list", networkName)
		}

		t.Logf("✅ Network listing completed successfully, found %d networks", len(networks))
	})

	t.Run("TestNetworkInspect", func(t *testing.T) {
		opts := &deploy.NetworkInspectOptions{}

		info, err := deploy.InspectNetwork(runner.Connection, networkName, opts)
		if err != nil {
			t.Fatalf("Failed to inspect network: %v", err)
		}

		if !strings.Contains(info, networkName) {
			t.Errorf("Expected network info to contain %s", networkName)
		}

		if !strings.Contains(info, "172.25.0.0/16") {
			t.Errorf("Expected network info to contain subnet 172.25.0.0/16")
		}

		t.Logf("✅ Network inspection completed successfully")
	})

	t.Run("TestNetworkConnectDisconnect", func(t *testing.T) {
		// Create a test container to connect to the network
		containerName := "test-net-container-" + helpers.RandomString(6)
		runner.Cleanup.AddContainer(containerName)

		opts := deploy.NewContainerOptions("alpine:latest")
		opts.Name = containerName
		opts.Command = []string{"sleep", "60"}

		containerID, err := deploy.Container(runner.Connection, opts)
		if err != nil {
			t.Fatalf("Failed to deploy container for network test: %v", err)
		}

		// Wait for container to be running
		if err := helpers.WaitForContainer(runner.Connection, containerName, "Up", 30*time.Second); err != nil {
			t.Fatalf("Container did not start: %v", err)
		}

		// Test network connect
		connectOpts := &deploy.NetworkConnectOptions{
			Alias: []string{"test-alias"},
			IP:    "172.25.0.100",
		}

		if err := deploy.ConnectContainerToNetwork(runner.Connection, networkName, containerName, connectOpts); err != nil {
			t.Fatalf("Failed to connect container to network: %v", err)
		}

		// Verify connection by inspecting the container
		inspectInfo, err := deploy.InspectObject(runner.Connection, containerName, &deploy.InspectOptions{})
		if err != nil {
			t.Fatalf("Failed to inspect container after network connect: %v", err)
		}

		if !strings.Contains(inspectInfo, networkName) {
			t.Errorf("Expected container to be connected to network %s", networkName)
		}

		// Test network disconnect
		disconnectOpts := &deploy.NetworkDisconnectOptions{
			Force: false,
		}

		if err := deploy.DisconnectContainerFromNetwork(runner.Connection, networkName, containerName, disconnectOpts); err != nil {
			t.Fatalf("Failed to disconnect container from network: %v", err)
		}

		shortID := containerID
		if len(containerID) > 12 {
			shortID = containerID[:12]
		}
		t.Logf("✅ Network connect/disconnect operations completed successfully (container: %s)", shortID)
	})

	t.Run("TestNetworkRemove", func(t *testing.T) {
		if err := deploy.RemoveNetwork(runner.Connection, networkName); err != nil {
			t.Fatalf("Failed to remove network: %v", err)
		}

		// Verify network is removed
		networks, err := deploy.ListNetworks(runner.Connection, &deploy.NetworkListOptions{})
		if err != nil {
			t.Fatalf("Failed to list networks after removal: %v", err)
		}

		for _, net := range networks {
			if net.Name == networkName {
				t.Error("Network should have been removed but still exists")
			}
		}

		t.Logf("✅ Network %s removed successfully", networkName)
	})
}

// testSystemOperations tests system df, info, prune, and inspect
func testSystemOperations(t *testing.T, runner *TestRunner) {
	t.Run("TestSystemDf", func(t *testing.T) {
		opts := &deploy.SystemDfOptions{
			Verbose: false,
		}

		usage, err := deploy.GetSystemDf(runner.Connection, opts)
		if err != nil {
			t.Fatalf("Failed to get system disk usage: %v", err)
		}

		if len(usage) == 0 {
			t.Error("Expected system disk usage data")
		}

		// Check that we have the expected categories
		categories := make(map[string]bool)
		for _, item := range usage {
			categories[item.Type] = true
		}

		// Check for essential categories (Build Cache and Local Volumes might be missing)
		essentialCategories := []string{"Images", "Containers"}
		for _, category := range essentialCategories {
			if !categories[category] {
				t.Errorf("Expected to find %s in system df output", category)
			}
		}

		// Log all found categories for debugging
		t.Logf("Found categories: %v", categories)

		t.Logf("✅ System disk usage retrieved successfully, found %d categories", len(usage))
	})

	t.Run("TestSystemInfo", func(t *testing.T) {
		opts := &deploy.SystemInfoOptions{}

		info, err := deploy.GetSystemInfo(runner.Connection, opts)
		if err != nil {
			t.Fatalf("Failed to get system info: %v", err)
		}

		if !strings.Contains(info, "Docker") {
			t.Errorf("Expected system info to contain 'Docker', got: %s", info)
		}

		if !strings.Contains(info, "Server Version") {
			t.Errorf("Expected system info to contain 'Server Version', got: %s", info)
		}

		t.Logf("✅ System info retrieved successfully")
	})

	t.Run("TestInspectObjects", func(t *testing.T) {
		// Create a test container to inspect
		containerName := "test-inspect-" + helpers.RandomString(6)
		runner.Cleanup.AddContainer(containerName)

		opts := deploy.NewContainerOptions("alpine:latest")
		opts.Name = containerName
		opts.Command = []string{"sleep", "30"}

		containerID, err := deploy.Container(runner.Connection, opts)
		if err != nil {
			t.Fatalf("Failed to deploy container for inspect test: %v", err)
		}

		// Test container inspection
		inspectOpts := &deploy.InspectOptions{
			Type: "container",
		}

		info, err := deploy.InspectObject(runner.Connection, containerName, inspectOpts)
		if err != nil {
			t.Fatalf("Failed to inspect container: %v", err)
		}

		if !strings.Contains(info, containerName) {
			t.Errorf("Expected container info to contain %s", containerName)
		}

		// Test with format template
		inspectOpts.Format = "'{{.State.Status}}'"
		statusInfo, err := deploy.InspectObject(runner.Connection, containerName, inspectOpts)
		if err != nil {
			t.Fatalf("Failed to inspect container with format: %v", err)
		}

		if !strings.Contains(statusInfo, "running") {
			t.Errorf("Expected container to be running, got status: %s", statusInfo)
		}

		shortID := containerID
		if len(containerID) > 12 {
			shortID = containerID[:12]
		}
		t.Logf("✅ Object inspection completed successfully (container: %s)", shortID)
	})
}
