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

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Manage Docker system",
	Long:  `Manage Docker system on your Synology NAS.`,
}

var systemDfCmd = &cobra.Command{
	Use:   "df [OPTIONS]",
	Short: "Show docker filesystem usage",
	Long:  `Show Docker filesystem usage on your Synology NAS.`,
	RunE:  showSystemDf,
}

var systemInfoCmd = &cobra.Command{
	Use:   "info [OPTIONS]",
	Short: "Display system-wide information",
	Long:  `Display Docker system information on your Synology NAS.`,
	RunE:  showSystemInfo,
}

var systemPruneCmd = &cobra.Command{
	Use:   "prune [OPTIONS]",
	Short: "Remove unused data",
	Long:  `Remove all unused containers, networks, images (both dangling and unreferenced), and optionally, volumes.`,
	RunE:  systemPrune,
}

var (
	systemDfFormat     string
	systemDfVerbose    bool
	systemInfoFormat   string
	systemPruneAll     bool
	systemPruneForce   bool
	systemPruneVolumes bool
	systemPruneFilter  []string
)

func showSystemDf(cmd *cobra.Command, args []string) error {
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

	// Show disk usage
	opts := &deploy.SystemDfOptions{
		Format:  systemDfFormat,
		Verbose: systemDfVerbose,
	}

	usage, err := deploy.GetSystemDf(conn, opts)
	if err != nil {
		return fmt.Errorf("failed to get system disk usage: %w", err)
	}

	// Display disk usage in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE\tTOTAL\tACTIVE\tSIZE\tRECLAIMABLE")

	for _, item := range usage {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			item.Type, item.Total, item.Active, item.Size, item.Reclaimable)
	}

	return w.Flush()
}

func showSystemInfo(cmd *cobra.Command, args []string) error {
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

	// Show system info
	opts := &deploy.SystemInfoOptions{
		Format: systemInfoFormat,
	}

	info, err := deploy.GetSystemInfo(conn, opts)
	if err != nil {
		return fmt.Errorf("failed to get system info: %w", err)
	}

	fmt.Print(info)
	return nil
}

func systemPrune(cmd *cobra.Command, args []string) error {
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

	// Confirm with user unless --force is specified
	if !systemPruneForce {
		fmt.Print("WARNING! This will remove:\n")
		fmt.Print("  - all stopped containers\n")
		fmt.Print("  - all networks not used by at least one container\n")
		fmt.Print("  - all dangling images\n")
		if systemPruneAll {
			fmt.Print("  - all images without at least one container associated to them\n")
		}
		if systemPruneVolumes {
			fmt.Print("  - all volumes not used by at least one container\n")
		}
		fmt.Print("Are you sure you want to continue? [y/N] ")

		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Perform system prune
	opts := &deploy.SystemPruneOptions{
		All:     systemPruneAll,
		Force:   systemPruneForce,
		Volumes: systemPruneVolumes,
		Filter:  systemPruneFilter,
	}

	result, err := deploy.SystemPrune(conn, opts)
	if err != nil {
		return fmt.Errorf("failed to prune system: %w", err)
	}

	fmt.Printf("Deleted Containers: %d\n", result.ContainersDeleted)
	fmt.Printf("Deleted Images: %d\n", result.ImagesDeleted)
	fmt.Printf("Deleted Networks: %d\n", result.NetworksDeleted)
	if systemPruneVolumes {
		fmt.Printf("Deleted Volumes: %d\n", result.VolumesDeleted)
	}
	fmt.Printf("Total reclaimed space: %s\n", result.SpaceReclaimed)

	return nil
}

func init() {
	// system df command
	systemDfCmd.Flags().StringVar(&systemDfFormat, "format", "", "Pretty-print disk usage using a Go template")
	systemDfCmd.Flags().BoolVarP(&systemDfVerbose, "verbose", "v", false, "Show detailed information on space usage")

	// system info command
	systemInfoCmd.Flags().StringVar(&systemInfoFormat, "format", "", "Format the output using the given Go template")

	// system prune command
	systemPruneCmd.Flags().BoolVarP(&systemPruneAll, "all", "a", false, "Remove all unused images not just dangling ones")
	systemPruneCmd.Flags().BoolVarP(&systemPruneForce, "force", "f", false, "Do not prompt for confirmation")
	systemPruneCmd.Flags().BoolVar(&systemPruneVolumes, "volumes", false, "Prune volumes")
	systemPruneCmd.Flags().StringSliceVar(&systemPruneFilter, "filter", []string{}, "Provide filter values (e.g. 'until=<timestamp>')")

	// Add subcommands to system
	systemCmd.AddCommand(systemDfCmd)
	systemCmd.AddCommand(systemInfoCmd)
	systemCmd.AddCommand(systemPruneCmd)
}
