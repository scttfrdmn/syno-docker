package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/syno-docker/pkg/config"
	"github.com/scttfrdmn/syno-docker/pkg/deploy"
	"github.com/scttfrdmn/syno-docker/pkg/synology"
)

var (
	pullAllTags      bool
	pullPlatform     string
	pullQuiet        bool
	pullDisableContentTrust bool
)

var pullCmd = &cobra.Command{
	Use:   "pull [OPTIONS] NAME[:TAG|@DIGEST]",
	Short: "Pull an image or a repository from a registry",
	Long:  `Pull an image or repository from a registry to your Synology NAS.`,
	Args:  cobra.ExactArgs(1),
	RunE:  pullImage,
}

func pullImage(cmd *cobra.Command, args []string) error {
	imageName := args[0]

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

	// Pull image
	opts := &deploy.PullOptions{
		AllTags:              pullAllTags,
		Platform:             pullPlatform,
		Quiet:                pullQuiet,
		DisableContentTrust:  pullDisableContentTrust,
	}

	fmt.Printf("Pulling image %s...\n", imageName)
	if err := deploy.PullImage(conn, imageName, opts); err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	fmt.Printf("✅ Image %s pulled successfully!\n", imageName)
	return nil
}

func init() {
	pullCmd.Flags().BoolVarP(&pullAllTags, "all-tags", "a", false, "Download all tagged images in the repository")
	pullCmd.Flags().StringVar(&pullPlatform, "platform", "", "Set platform if server is multi-platform capable")
	pullCmd.Flags().BoolVarP(&pullQuiet, "quiet", "q", false, "Suppress verbose output")
	pullCmd.Flags().BoolVar(&pullDisableContentTrust, "disable-content-trust", true, "Skip image verification")
}