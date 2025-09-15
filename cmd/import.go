package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/syno-docker/pkg/config"
	"github.com/scttfrdmn/syno-docker/pkg/deploy"
	"github.com/scttfrdmn/syno-docker/pkg/synology"
)

var importCmd = &cobra.Command{
	Use:   "import [OPTIONS] file|URL|- [REPOSITORY[:TAG]]",
	Short: "Import the contents from a tarball to create a filesystem image",
	Long:  `Import the contents from a tarball to create a filesystem image on your Synology NAS.`,
	Args:  cobra.RangeArgs(1, 2),
	RunE:  importImage,
}

var (
	importChange   []string
	importMessage  string
	importPlatform string
)

func importImage(cmd *cobra.Command, args []string) error {
	source := args[0]
	var repository string
	if len(args) > 1 {
		repository = args[1]
	}

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

	// Import image
	opts := &deploy.ImportOptions{
		Change:   importChange,
		Message:  importMessage,
		Platform: importPlatform,
	}

	fmt.Printf("Importing image from %s...\n", source)
	imageID, err := deploy.ImportImage(conn, source, repository, opts)
	if err != nil {
		return fmt.Errorf("failed to import image: %w", err)
	}

	if repository != "" {
		fmt.Printf("✅ Image imported as %s (ID: %s) successfully!\n", repository, imageID)
	} else {
		fmt.Printf("✅ Image imported with ID %s successfully!\n", imageID)
	}
	return nil
}

func init() {
	importCmd.Flags().StringSliceVarP(&importChange, "change", "c", []string{}, "Apply Dockerfile instruction to the created image")
	importCmd.Flags().StringVarP(&importMessage, "message", "m", "", "Set commit message for imported image")
	importCmd.Flags().StringVar(&importPlatform, "platform", "", "Set platform if server is multi-platform capable")
}