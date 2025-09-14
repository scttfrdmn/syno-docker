package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/synodeploy/pkg/config"
	"github.com/scttfrdmn/synodeploy/pkg/deploy"
	"github.com/scttfrdmn/synodeploy/pkg/synology"
)

var rmForce bool

var rmCmd = &cobra.Command{
	Use:   "rm <container>",
	Short: "Remove container",
	Long:  `Remove a container from your Synology NAS by name or ID.`,
	Args:  cobra.ExactArgs(1),
	RunE:  removeContainer,
}

func removeContainer(cmd *cobra.Command, args []string) error {
	containerNameOrID := args[0]

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Connect to Synology NAS
	fmt.Printf("Connecting to %s@%s:%d...\n", cfg.User, cfg.Host, cfg.Port)
	conn := synology.NewConnection(cfg)
	if err := conn.Connect(); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()

	// Remove container
	fmt.Printf("Removing container %s...\n", containerNameOrID)
	if err := deploy.RemoveContainer(conn, containerNameOrID, rmForce); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	fmt.Printf("✅ Container %s removed successfully!\n", containerNameOrID)
	return nil
}

func init() {
	rmCmd.Flags().BoolVarP(&rmForce, "force", "f", false, "Force removal of running container")
}
