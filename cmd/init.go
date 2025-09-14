package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/synodeploy/pkg/config"
	"github.com/scttfrdmn/synodeploy/pkg/synology"
)

var (
	initUser       string
	initPort       int
	initSSHKey     string
	initVolumePath string
)

var initCmd = &cobra.Command{
	Use:   "init <host>",
	Short: "Setup connection to Synology NAS",
	Long: `Initialize SynoDeploy configuration for connecting to your Synology NAS.
This command sets up SSH connection details and tests the connection.`,
	Args: cobra.ExactArgs(1),
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	host := args[0]

	// Create new config
	cfg := config.New()
	cfg.Host = host
	cfg.User = initUser
	cfg.Port = initPort
	cfg.SSHKeyPath = initSSHKey
	cfg.Defaults.VolumePath = initVolumePath

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Test connection
	fmt.Printf("Testing connection to %s@%s:%d...\n", cfg.User, cfg.Host, cfg.Port)
	conn := synology.NewConnection(cfg)
	if err := conn.Connect(); err != nil {
		return fmt.Errorf("connection test failed: %w\n\nTry:\n  1. Verify host is reachable: ping %s\n  2. Check SSH service is enabled on your NAS\n  3. Verify username and SSH key path\n  4. Ensure your user has admin privileges", err, cfg.Host)
	}
	defer conn.Close()

	// Additional connection tests
	if err := conn.TestConnection(); err != nil {
		return fmt.Errorf("docker connection test failed: %w\n\nTry:\n  1. Ensure Container Manager is installed and running\n  2. Verify your user is in the docker group\n  3. Check if Docker service is running: systemctl status pkg-ContainerManager-dockerd", err)
	}

	// Save configuration
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	configPath, _ := config.GetConfigPath()
	fmt.Printf("✅ Connection successful!\nConfiguration saved to %s\n", configPath)
	fmt.Printf("You can now deploy containers using 'synodeploy run' or 'synodeploy deploy'\n")

	return nil
}

func init() {
	// Default SSH key path
	homeDir, _ := os.UserHomeDir()
	defaultSSHKey := filepath.Join(homeDir, ".ssh", "id_rsa")

	initCmd.Flags().StringVarP(&initUser, "user", "u", config.DefaultUser, "SSH username")
	initCmd.Flags().IntVarP(&initPort, "port", "p", config.DefaultPort, "SSH port")
	initCmd.Flags().StringVarP(&initSSHKey, "key", "k", defaultSSHKey, "SSH private key path")
	initCmd.Flags().StringVar(&initVolumePath, "volume-path", config.DefaultVolumePath, "Default volume path on NAS")
}
