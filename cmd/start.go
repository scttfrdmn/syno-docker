package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/syno-docker/pkg/config"
	"github.com/scttfrdmn/syno-docker/pkg/deploy"
	"github.com/scttfrdmn/syno-docker/pkg/synology"
)

var startCmd = &cobra.Command{
	Use:   "start <container> [containers...]",
	Short: "Start one or more stopped containers",
	Long:  `Start one or more stopped containers on your Synology NAS.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  startContainers,
}

func startContainers(cmd *cobra.Command, args []string) error {
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

	// Start each container
	for _, containerNameOrID := range args {
		fmt.Printf("Starting container %s...\n", containerNameOrID)
		if err := deploy.StartContainer(conn, containerNameOrID); err != nil {
			return fmt.Errorf("failed to start container %s: %w", containerNameOrID, err)
		}
		fmt.Printf("✅ Container %s started successfully!\n", containerNameOrID)
	}

	return nil
}
