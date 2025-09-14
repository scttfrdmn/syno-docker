package integration

import (
	"flag"
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

var (
	integrationTest = flag.Bool("integration", false, "Run integration tests")
	nasHost         = flag.String("nas-host", "", "Synology NAS IP address or hostname")
	nasUser         = flag.String("nas-user", "admin", "SSH username (admin or administrator)")
	nasPort         = flag.Int("nas-port", 22, "SSH port")
	nasKeyPath      = flag.String("nas-key", "", "SSH private key path")
	cleanup         = flag.Bool("cleanup", true, "Cleanup test resources after tests")
)

// TestRunner manages the integration test environment
type TestRunner struct {
	Config     *config.Config
	Connection *synology.Connection
	Cleanup    *helpers.CleanupManager
}

func TestMain(m *testing.M) {
	flag.Parse()

	// Only run TestMain setup for tests that actually need it (TestIntegration)
	// Other tests like TestConnectionToChubChub can run independently

	os.Exit(m.Run())
}

func setupTestEnvironment() (*TestRunner, error) {
	// Set default SSH key path if not provided
	sshKeyPath := *nasKeyPath
	if sshKeyPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		sshKeyPath = filepath.Join(homeDir, ".ssh", "id_rsa")
	}

	// Create test configuration
	cfg := &config.Config{
		Host:       *nasHost,
		User:       *nasUser,
		Port:       *nasPort,
		SSHKeyPath: sshKeyPath,
		Defaults: struct {
			VolumePath string `yaml:"volume_path"`
			Network    string `yaml:"network,omitempty"`
		}{
			VolumePath: "/volume1/synodeploy-test",
			Network:    "bridge",
		},
	}

	fmt.Printf("🔌 Testing connection to %s@%s:%d using key %s\n", cfg.User, cfg.Host, cfg.Port, cfg.SSHKeyPath)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Test connection
	conn := synology.NewConnection(cfg)
	if err := conn.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to NAS: %w", err)
	}

	// Test Docker availability
	if err := conn.TestConnection(); err != nil {
		return nil, fmt.Errorf("Docker not available on NAS: %w", err)
	}

	// Initialize cleanup manager
	cleanupMgr := helpers.NewCleanupManager(conn)

	// Create test volume directory
	if err := helpers.CreateTestDirectory(conn, cfg.Defaults.VolumePath); err != nil {
		return nil, fmt.Errorf("failed to create test directory: %w", err)
	}
	cleanupMgr.AddDirectory(cfg.Defaults.VolumePath)

	return &TestRunner{
		Config:     cfg,
		Connection: conn,
		Cleanup:    cleanupMgr,
	}, nil
}

// Global test runner instance
var testRunner *TestRunner

func TestIntegration(t *testing.T) {
	if !*integrationTest {
		t.Skip("Integration tests not enabled. Use -integration flag.")
	}

	if *nasHost == "" {
		t.Skip("NAS host not specified. Use -nas-host flag.")
	}

	var err error
	testRunner, err = setupTestEnvironment()
	if err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}

	t.Run("BasicDeployment", testBasicDeployment)
	// Additional test functions will be implemented in future versions
}

func testBasicDeployment(t *testing.T) {
	// Test single container deployment
	containerName := "test-nginx-" + helpers.RandomString(8)
	testRunner.Cleanup.AddContainer(containerName)

	// Deploy nginx container
	opts := &deploy.ContainerOptions{
		Image:       "nginx:alpine",
		Name:        containerName,
		Ports:       []string{"8080:80"},
		Volumes:     []string{fmt.Sprintf("%s/html:/usr/share/nginx/html", testRunner.Config.Defaults.VolumePath)},
		Restart:     "unless-stopped",
		NetworkMode: "bridge",
	}

	// Create test HTML file
	htmlContent := fmt.Sprintf("<h1>SynoDeploy Test - %s</h1>", time.Now().Format("2006-01-02 15:04:05"))
	if err := helpers.CreateTestFile(testRunner.Connection,
		fmt.Sprintf("%s/html/index.html", testRunner.Config.Defaults.VolumePath),
		htmlContent); err != nil {
		t.Fatalf("Failed to create test HTML file: %v", err)
	}

	// Deploy container
	containerID, err := deploy.DeployContainer(testRunner.Connection, opts)
	if err != nil {
		t.Fatalf("Failed to deploy container: %v", err)
	}

	t.Logf("Deployed container: %s (ID: %s)", containerName, containerID)

	// Wait for container to be ready
	time.Sleep(5 * time.Second)

	// Verify container is running
	containers, err := deploy.ListContainers(testRunner.Connection, false)
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
			break
		}
	}

	if !found {
		t.Errorf("Container %s not found in running containers", containerName)
	}

	// Test HTTP connectivity
	if err := helpers.TestHTTPEndpoint(fmt.Sprintf("http://%s:8080", testRunner.Config.Host),
		"SynoDeploy Test", 30*time.Second); err != nil {
		t.Errorf("HTTP connectivity test failed: %v", err)
	}
}
