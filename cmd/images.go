package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/syno-docker/pkg/config"
	"github.com/scttfrdmn/syno-docker/pkg/deploy"
	"github.com/scttfrdmn/syno-docker/pkg/synology"
)

var (
	imagesAll      bool
	imagesDangling bool
	imagesDigests  bool
	imagesFormat   string
	imagesNoTrunc  bool
	imagesQuiet    bool
)

var imagesCmd = &cobra.Command{
	Use:     "images [OPTIONS] [REPOSITORY[:TAG]]",
	Short:   "List images",
	Long:    `List Docker images on your Synology NAS.`,
	RunE:    listImages,
	Aliases: []string{"image", "img"},
}

func listImages(cmd *cobra.Command, args []string) error {
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

	// List images
	opts := &deploy.ImagesOptions{
		All:      imagesAll,
		Dangling: imagesDangling,
		Digests:  imagesDigests,
		Format:   imagesFormat,
		NoTrunc:  imagesNoTrunc,
		Quiet:    imagesQuiet,
	}

	var repository string
	if len(args) > 0 {
		repository = args[0]
	}

	if imagesQuiet {
		imageIDs, err := deploy.ListImageIDs(conn, repository, opts)
		if err != nil {
			return fmt.Errorf("failed to list images: %w", err)
		}
		for _, id := range imageIDs {
			fmt.Println(id)
		}
		return nil
	}

	images, err := deploy.ListImages(conn, repository, opts)
	if err != nil {
		return fmt.Errorf("failed to list images: %w", err)
	}

	if len(images) == 0 {
		fmt.Println("No images found.")
		return nil
	}

	// Display images in a table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if imagesDigests {
		fmt.Fprintln(w, "REPOSITORY\tTAG\tDIGEST\tIMAGE ID\tCREATED\tSIZE")
	} else {
		fmt.Fprintln(w, "REPOSITORY\tTAG\tIMAGE ID\tCREATED\tSIZE")
	}

	for _, image := range images {
		if imagesDigests {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				image.Repository, image.Tag, image.Digest, image.ID, image.Created, image.Size)
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				image.Repository, image.Tag, image.ID, image.Created, image.Size)
		}
	}

	return w.Flush()
}

func init() {
	imagesCmd.Flags().BoolVarP(&imagesAll, "all", "a", false, "Show all images (default hides intermediate images)")
	imagesCmd.Flags().BoolVar(&imagesDangling, "dangling", false, "Show only dangling images")
	imagesCmd.Flags().BoolVar(&imagesDigests, "digests", false, "Show digests")
	imagesCmd.Flags().StringVar(&imagesFormat, "format", "", "Pretty-print images using a Go template")
	imagesCmd.Flags().BoolVar(&imagesNoTrunc, "no-trunc", false, "Don't truncate output")
	imagesCmd.Flags().BoolVarP(&imagesQuiet, "quiet", "q", false, "Only show image IDs")
}
