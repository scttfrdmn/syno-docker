package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/syno-docker/pkg/config"
	"github.com/scttfrdmn/syno-docker/pkg/deploy"
	"github.com/scttfrdmn/syno-docker/pkg/synology"
)

var (
	statsAll      bool
	statsNoStream bool
	statsFormat   string
)

var statsCmd = &cobra.Command{
	Use:   "stats [container...]",
	Short: "Display a live stream of container(s) resource usage statistics",
	Long:  `Display a live stream of container resource usage statistics from your Synology NAS.`,
	RunE:  showContainerStats,
}

func showContainerStats(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Connect to Synology NAS
	conn := synology.NewConnection(cfg)
	if err := conn.Connect(); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()

	// Show container statistics
	opts := &deploy.StatsOptions{
		All:      statsAll,
		NoStream: statsNoStream,
		Format:   statsFormat,
	}

	if len(args) == 0 {
		// Show stats for all containers
		fmt.Println("Showing statistics for all containers...")
		return deploy.ShowContainerStats(conn, nil, opts)
	} else {
		// Show stats for specific containers
		fmt.Printf("Showing statistics for containers: %v...\n", args)
		return deploy.ShowContainerStats(conn, args, opts)
	}
}

func init() {
	statsCmd.Flags().BoolVarP(&statsAll, "all", "a", false, "Show all containers (default shows just running)")
	statsCmd.Flags().BoolVar(&statsNoStream, "no-stream", false, "Disable streaming stats and only pull the first result")
	statsCmd.Flags().StringVar(&statsFormat, "format", "", "Pretty-print images using a Go template")
}