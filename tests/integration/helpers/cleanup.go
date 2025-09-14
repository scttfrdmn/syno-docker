package helpers

import (
	"fmt"
	"strings"
	"sync"

	"github.com/scttfrdmn/synodeploy/pkg/deploy"
	"github.com/scttfrdmn/synodeploy/pkg/synology"
)

// CleanupManager manages cleanup of test resources
type CleanupManager struct {
	conn       *synology.Connection
	containers []string
	volumes    []string
	files      []string
	dirs       []string
	networks   []string
	mutex      sync.Mutex
}

// NewCleanupManager creates a new cleanup manager
func NewCleanupManager(conn *synology.Connection) *CleanupManager {
	return &CleanupManager{
		conn:       conn,
		containers: make([]string, 0),
		volumes:    make([]string, 0),
		files:      make([]string, 0),
		dirs:       make([]string, 0),
		networks:   make([]string, 0),
	}
}

// AddContainer adds a container to cleanup list
func (cm *CleanupManager) AddContainer(name string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.containers = append(cm.containers, name)
}

// AddVolume adds a volume to cleanup list
func (cm *CleanupManager) AddVolume(name string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.volumes = append(cm.volumes, name)
}

// AddFile adds a file to cleanup list
func (cm *CleanupManager) AddFile(path string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.files = append(cm.files, path)
}

// AddDirectory adds a directory to cleanup list
func (cm *CleanupManager) AddDirectory(path string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.dirs = append(cm.dirs, path)
}

// AddNetwork adds a network to cleanup list
func (cm *CleanupManager) AddNetwork(name string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.networks = append(cm.networks, name)
}

// CleanupContainers removes all tracked containers
func (cm *CleanupManager) CleanupContainers() error {
	cm.mutex.Lock()
	containers := make([]string, len(cm.containers))
	copy(containers, cm.containers)
	cm.mutex.Unlock()

	var errors []string

	for _, containerName := range containers {
		fmt.Printf("Cleaning up container: %s\n", containerName)

		// Stop and remove container
		if err := deploy.RemoveContainer(cm.conn, containerName, true); err != nil {
			errors = append(errors, fmt.Sprintf("failed to remove container %s: %v", containerName, err))
			continue
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("container cleanup errors: %s", strings.Join(errors, "; "))
	}

	// Clear the list
	cm.mutex.Lock()
	cm.containers = cm.containers[:0]
	cm.mutex.Unlock()

	return nil
}

// CleanupVolumes removes all tracked volumes
func (cm *CleanupManager) CleanupVolumes() error {
	cm.mutex.Lock()
	volumes := make([]string, len(cm.volumes))
	copy(volumes, cm.volumes)
	cm.mutex.Unlock()

	var errors []string

	for _, volumeName := range volumes {
		fmt.Printf("Cleaning up volume: %s\n", volumeName)

		_, err := cm.conn.ExecuteDockerCommand([]string{"volume", "rm", "-f", volumeName})
		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to remove volume %s: %v", volumeName, err))
			continue
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("volume cleanup errors: %s", strings.Join(errors, "; "))
	}

	// Clear the list
	cm.mutex.Lock()
	cm.volumes = cm.volumes[:0]
	cm.mutex.Unlock()

	return nil
}

// CleanupFiles removes all tracked files
func (cm *CleanupManager) CleanupFiles() error {
	cm.mutex.Lock()
	files := make([]string, len(cm.files))
	copy(files, cm.files)
	cm.mutex.Unlock()

	var errors []string

	for _, filePath := range files {
		fmt.Printf("Cleaning up file: %s\n", filePath)

		_, err := cm.conn.ExecuteCommand(fmt.Sprintf("rm -f %s", filePath))
		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to remove file %s: %v", filePath, err))
			continue
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("file cleanup errors: %s", strings.Join(errors, "; "))
	}

	// Clear the list
	cm.mutex.Lock()
	cm.files = cm.files[:0]
	cm.mutex.Unlock()

	return nil
}

// CleanupDirectories removes all tracked directories
func (cm *CleanupManager) CleanupDirectories() error {
	cm.mutex.Lock()
	dirs := make([]string, len(cm.dirs))
	copy(dirs, cm.dirs)
	cm.mutex.Unlock()

	var errors []string

	for _, dirPath := range dirs {
		fmt.Printf("Cleaning up directory: %s\n", dirPath)

		_, err := cm.conn.ExecuteCommand(fmt.Sprintf("rm -rf %s", dirPath))
		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to remove directory %s: %v", dirPath, err))
			continue
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("directory cleanup errors: %s", strings.Join(errors, "; "))
	}

	// Clear the list
	cm.mutex.Lock()
	cm.dirs = cm.dirs[:0]
	cm.mutex.Unlock()

	return nil
}

// CleanupNetworks removes all tracked networks
func (cm *CleanupManager) CleanupNetworks() error {
	cm.mutex.Lock()
	networks := make([]string, len(cm.networks))
	copy(networks, cm.networks)
	cm.mutex.Unlock()

	var errors []string

	for _, networkName := range networks {
		fmt.Printf("Cleaning up network: %s\n", networkName)

		_, err := cm.conn.ExecuteDockerCommand([]string{"network", "rm", networkName})
		if err != nil {
			// Network might not exist or be in use, so warn but continue
			fmt.Printf("Warning: failed to remove network %s: %v\n", networkName, err)
			continue
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("network cleanup errors: %s", strings.Join(errors, "; "))
	}

	// Clear the list
	cm.mutex.Lock()
	cm.networks = cm.networks[:0]
	cm.mutex.Unlock()

	return nil
}

// CleanupAll removes all tracked resources
func (cm *CleanupManager) CleanupAll() error {
	var allErrors []string

	// Cleanup in reverse order of dependency
	if err := cm.CleanupContainers(); err != nil {
		allErrors = append(allErrors, err.Error())
	}

	if err := cm.CleanupNetworks(); err != nil {
		allErrors = append(allErrors, err.Error())
	}

	if err := cm.CleanupVolumes(); err != nil {
		allErrors = append(allErrors, err.Error())
	}

	if err := cm.CleanupFiles(); err != nil {
		allErrors = append(allErrors, err.Error())
	}

	if err := cm.CleanupDirectories(); err != nil {
		allErrors = append(allErrors, err.Error())
	}

	if len(allErrors) > 0 {
		return fmt.Errorf("cleanup errors: %s", strings.Join(allErrors, "; "))
	}

	fmt.Println("All test resources cleaned up successfully")
	return nil
}

// ForceCleanupContainers removes all containers with synodeploy test prefix
func (cm *CleanupManager) ForceCleanupContainers() error {
	fmt.Println("Force cleaning up all synodeploy test containers...")

	// List all containers
	containers, err := deploy.ListContainers(cm.conn, true)
	if err != nil {
		return fmt.Errorf("failed to list containers for force cleanup: %w", err)
	}

	var errors []string
	for _, container := range containers {
		// Only clean up test containers (those starting with test- prefix)
		if strings.HasPrefix(container.Name, "test-") {
			fmt.Printf("Force removing container: %s\n", container.Name)
			if err := deploy.RemoveContainer(cm.conn, container.Name, true); err != nil {
				errors = append(errors, fmt.Sprintf("failed to force remove container %s: %v", container.Name, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("force cleanup errors: %s", strings.Join(errors, "; "))
	}

	return nil
}