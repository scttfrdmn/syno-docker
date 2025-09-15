package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/syno-docker/pkg/config"
	"github.com/scttfrdmn/syno-docker/pkg/deploy"
	"github.com/scttfrdmn/syno-docker/pkg/synology"
)

var stopTimeout int

var stopCmd = &cobra.Command{
	Use:   "stop <container> [containers...]",
	Short: "Stop one or more running containers",
	Long:  `Stop one or more running containers on your Synology NAS.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  stopContainers,
}

func stopContainers(cmd *cobra.Command, args []string) error {
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

	// Stop each container
	for _, containerNameOrID := range args {
		fmt.Printf("Stopping container %s...\n", containerNameOrID)
		if err := deploy.StopContainer(conn, containerNameOrID, stopTimeout); err != nil {
			return fmt.Errorf("failed to stop container %s: %w", containerNameOrID, err)
		}
		fmt.Printf("✅ Container %s stopped successfully!\n", containerNameOrID)
	}

	return nil
}

func init() {
	stopCmd.Flags().IntVarP(&stopTimeout, "time", "t", 10, "Seconds to wait for stop before killing container")
}
