package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/syno-docker/pkg/config"
	"github.com/scttfrdmn/syno-docker/pkg/deploy"
	"github.com/scttfrdmn/syno-docker/pkg/synology"
)

var (
	rmiForce   bool
	rmiNoPrune bool
)

var rmiCmd = &cobra.Command{
	Use:   "rmi [OPTIONS] IMAGE [IMAGE...]",
	Short: "Remove one or more images",
	Long:  `Remove one or more Docker images from your Synology NAS.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  removeImages,
}

func removeImages(cmd *cobra.Command, args []string) error {
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

	// Remove each image
	opts := &deploy.RmiOptions{
		Force:   rmiForce,
		NoPrune: rmiNoPrune,
	}

	for _, imageName := range args {
		fmt.Printf("Removing image %s...\n", imageName)
		if err := deploy.RemoveImage(conn, imageName, opts); err != nil {
			return fmt.Errorf("failed to remove image %s: %w", imageName, err)
		}
		fmt.Printf("✅ Image %s removed successfully!\n", imageName)
	}

	return nil
}

func init() {
	rmiCmd.Flags().BoolVarP(&rmiForce, "force", "f", false, "Force removal of the image")
	rmiCmd.Flags().BoolVar(&rmiNoPrune, "no-prune", false, "Do not delete untagged parents")
}
