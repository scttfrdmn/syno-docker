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

var psAll bool

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List containers",
	Long:  `List containers running on your Synology NAS.`,
	RunE:  listContainers,
}

func listContainers(cmd *cobra.Command, args []string) error {
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

	// List containers
	containers, err := deploy.ListContainers(conn, psAll)
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	if len(containers) == 0 {
		fmt.Println("No containers found.")
		return nil
	}

	// Display containers in a table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "CONTAINER ID\tNAME\tIMAGE\tSTATUS\tPORTS")

	for _, container := range containers {
		ports := "-"
		if len(container.Ports) > 0 {
			ports = fmt.Sprintf("%v", container.Ports)
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			container.ID,
			container.Name,
			container.Image,
			container.Status,
			ports,
		)
	}

	return w.Flush()
}

func init() {
	psCmd.Flags().BoolVarP(&psAll, "all", "a", false, "Show all containers (default: running only)")
}
