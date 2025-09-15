package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/syno-docker/pkg/config"
	"github.com/scttfrdmn/syno-docker/pkg/deploy"
	"github.com/scttfrdmn/syno-docker/pkg/synology"
)

var exportCmd = &cobra.Command{
	Use:   "export [OPTIONS] CONTAINER",
	Short: "Export a container's filesystem as a tar archive",
	Long:  `Export a container's filesystem as a tar archive to your local machine.`,
	Args:  cobra.ExactArgs(1),
	RunE:  exportContainer,
}

var (
	exportOutput string
)

func exportContainer(cmd *cobra.Command, args []string) error {
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

	// Export container
	opts := &deploy.ExportOptions{
		Output: exportOutput,
	}

	fmt.Printf("Exporting container %s...\n", containerNameOrID)
	if err := deploy.ExportContainer(conn, containerNameOrID, opts); err != nil {
		return fmt.Errorf("failed to export container: %w", err)
	}

	if exportOutput != "" {
		fmt.Printf("✅ Container %s exported to %s successfully!\n", containerNameOrID, exportOutput)
	} else {
		fmt.Printf("✅ Container %s exported successfully!\n", containerNameOrID)
	}
	return nil
}

func init() {
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Write to a file, instead of STDOUT")
}
