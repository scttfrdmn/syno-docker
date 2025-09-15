package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/syno-docker/pkg/config"
	"github.com/scttfrdmn/syno-docker/pkg/deploy"
	"github.com/scttfrdmn/syno-docker/pkg/synology"
)

var restartTimeout int

var restartCmd = &cobra.Command{
	Use:   "restart <container> [containers...]",
	Short: "Restart one or more containers",
	Long:  `Restart one or more containers on your Synology NAS.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  restartContainers,
}

func restartContainers(cmd *cobra.Command, args []string) error {
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

	// Restart each container
	for _, containerNameOrID := range args {
		fmt.Printf("Restarting container %s...\n", containerNameOrID)
		if err := deploy.RestartContainer(conn, containerNameOrID, restartTimeout); err != nil {
			return fmt.Errorf("failed to restart container %s: %w", containerNameOrID, err)
		}
		fmt.Printf("✅ Container %s restarted successfully!\n", containerNameOrID)
	}

	return nil
}

func init() {
	restartCmd.Flags().IntVarP(&restartTimeout, "time", "t", 10, "Seconds to wait for stop before killing container")
}
