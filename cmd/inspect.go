package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/syno-docker/pkg/config"
	"github.com/scttfrdmn/syno-docker/pkg/deploy"
	"github.com/scttfrdmn/syno-docker/pkg/synology"
)

var (
	inspectFormat string
	inspectSize   bool
	inspectType   string
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [OPTIONS] NAME|ID [NAME|ID...]",
	Short: "Return low-level information on Docker objects",
	Long:  `Return low-level information on Docker containers, images, volumes, and networks.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  inspectObjects,
}

func inspectObjects(cmd *cobra.Command, args []string) error {
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

	// Inspect objects
	opts := &deploy.InspectOptions{
		Format: inspectFormat,
		Size:   inspectSize,
		Type:   inspectType,
	}

	for _, objectName := range args {
		info, err := deploy.InspectObject(conn, objectName, opts)
		if err != nil {
			return fmt.Errorf("failed to inspect object %s: %w", objectName, err)
		}
		fmt.Print(info)
	}

	return nil
}

func init() {
	inspectCmd.Flags().StringVarP(&inspectFormat, "format", "f", "", "Format the output using the given Go template")
	inspectCmd.Flags().BoolVarP(&inspectSize, "size", "s", false, "Display total file sizes if the type is container")
	inspectCmd.Flags().StringVar(&inspectType, "type", "", "Return JSON for specified type (container, image, volume, network)")
}
