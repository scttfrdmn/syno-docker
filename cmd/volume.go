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

var volumeCmd = &cobra.Command{
	Use:   "volume",
	Short: "Manage volumes",
	Long:  `Manage Docker volumes on your Synology NAS.`,
}

var volumeListCmd = &cobra.Command{
	Use:     "ls [OPTIONS]",
	Short:   "List volumes",
	Long:    `List Docker volumes on your Synology NAS.`,
	RunE:    listVolumes,
	Aliases: []string{"list"},
}

var volumeCreateCmd = &cobra.Command{
	Use:   "create [OPTIONS] [VOLUME]",
	Short: "Create a volume",
	Long:  `Create a Docker volume on your Synology NAS.`,
	RunE:  createVolume,
}

var volumeRemoveCmd = &cobra.Command{
	Use:     "rm VOLUME [VOLUME...]",
	Short:   "Remove one or more volumes",
	Long:    `Remove one or more Docker volumes from your Synology NAS.`,
	Args:    cobra.MinimumNArgs(1),
	RunE:    removeVolumes,
	Aliases: []string{"remove"},
}

var volumeInspectCmd = &cobra.Command{
	Use:   "inspect VOLUME [VOLUME...]",
	Short: "Display detailed information on one or more volumes",
	Long:  `Display detailed information on one or more Docker volumes.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  inspectVolumes,
}

var volumePruneCmd = &cobra.Command{
	Use:   "prune [OPTIONS]",
	Short: "Remove all unused local volumes",
	Long:  `Remove all unused local Docker volumes.`,
	RunE:  pruneVolumes,
}

var (
	volumeListFormat    string
	volumeListQuiet     bool
	volumeCreateDriver  string
	volumeCreateLabel   []string
	volumeCreateOptions []string
	volumeRemoveForce   bool
	volumeInspectFormat string
	volumePruneForce    bool
	volumePruneFilter   []string
)

func listVolumes(cmd *cobra.Command, args []string) error {
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

	// List volumes
	opts := &deploy.VolumeListOptions{
		Format: volumeListFormat,
		Quiet:  volumeListQuiet,
	}

	if volumeListQuiet {
		volumeNames, err := deploy.ListVolumeNames(conn, opts)
		if err != nil {
			return fmt.Errorf("failed to list volumes: %w", err)
		}
		for _, name := range volumeNames {
			fmt.Println(name)
		}
		return nil
	}

	volumes, err := deploy.ListVolumes(conn, opts)
	if err != nil {
		return fmt.Errorf("failed to list volumes: %w", err)
	}

	if len(volumes) == 0 {
		fmt.Println("No volumes found.")
		return nil
	}

	// Display volumes in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DRIVER\tVOLUME NAME")

	for _, volume := range volumes {
		fmt.Fprintf(w, "%s\t%s\n", volume.Driver, volume.Name)
	}

	return w.Flush()
}

func createVolume(cmd *cobra.Command, args []string) error {
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

	// Create volume
	var volumeName string
	if len(args) > 0 {
		volumeName = args[0]
	}

	opts := &deploy.VolumeCreateOptions{
		Driver:  volumeCreateDriver,
		Labels:  volumeCreateLabel,
		Options: volumeCreateOptions,
	}

	name, err := deploy.CreateVolume(conn, volumeName, opts)
	if err != nil {
		return fmt.Errorf("failed to create volume: %w", err)
	}

	fmt.Printf("✅ Volume %s created successfully!\n", name)
	return nil
}

func removeVolumes(cmd *cobra.Command, args []string) error {
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

	// Remove each volume
	opts := &deploy.VolumeRemoveOptions{
		Force: volumeRemoveForce,
	}

	for _, volumeName := range args {
		fmt.Printf("Removing volume %s...\n", volumeName)
		if err := deploy.RemoveVolume(conn, volumeName, opts); err != nil {
			return fmt.Errorf("failed to remove volume %s: %w", volumeName, err)
		}
		fmt.Printf("✅ Volume %s removed successfully!\n", volumeName)
	}

	return nil
}

func inspectVolumes(cmd *cobra.Command, args []string) error {
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

	// Inspect volumes
	opts := &deploy.VolumeInspectOptions{
		Format: volumeInspectFormat,
	}

	for _, volumeName := range args {
		info, err := deploy.InspectVolume(conn, volumeName, opts)
		if err != nil {
			return fmt.Errorf("failed to inspect volume %s: %w", volumeName, err)
		}
		fmt.Print(info)
	}

	return nil
}

func pruneVolumes(cmd *cobra.Command, args []string) error {
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
	if !volumePruneForce {
		fmt.Print("WARNING! This will remove all local volumes not used by at least one container.\n")
		fmt.Print("Are you sure you want to continue? [y/N] ")

		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Prune volumes
	opts := &deploy.VolumePruneOptions{
		Force:  volumePruneForce,
		Filter: volumePruneFilter,
	}

	result, err := deploy.PruneVolumes(conn, opts)
	if err != nil {
		return fmt.Errorf("failed to prune volumes: %w", err)
	}

	fmt.Printf("Deleted Volumes: %d\n", result.VolumesDeleted)
	fmt.Printf("Total reclaimed space: %s\n", result.SpaceReclaimed)

	return nil
}

func init() {
	// volume list command
	volumeListCmd.Flags().StringVar(&volumeListFormat, "format", "", "Pretty-print volumes using a Go template")
	volumeListCmd.Flags().BoolVarP(&volumeListQuiet, "quiet", "q", false, "Only display volume names")

	// volume create command
	volumeCreateCmd.Flags().StringVarP(&volumeCreateDriver, "driver", "d", "local", "Specify volume driver name")
	volumeCreateCmd.Flags().StringSliceVar(&volumeCreateLabel, "label", []string{}, "Set metadata for a volume")
	volumeCreateCmd.Flags().StringSliceVarP(&volumeCreateOptions, "opt", "o", []string{}, "Set driver specific options")

	// volume remove command
	volumeRemoveCmd.Flags().BoolVarP(&volumeRemoveForce, "force", "f", false, "Force the removal of one or more volumes")

	// volume inspect command
	volumeInspectCmd.Flags().StringVarP(&volumeInspectFormat, "format", "f", "", "Format the output using the given Go template")

	// volume prune command
	volumePruneCmd.Flags().BoolVarP(&volumePruneForce, "force", "f", false, "Do not prompt for confirmation")
	volumePruneCmd.Flags().StringSliceVar(&volumePruneFilter, "filter", []string{}, "Provide filter values")

	// Add subcommands to volume
	volumeCmd.AddCommand(volumeListCmd)
	volumeCmd.AddCommand(volumeCreateCmd)
	volumeCmd.AddCommand(volumeRemoveCmd)
	volumeCmd.AddCommand(volumeInspectCmd)
	volumeCmd.AddCommand(volumePruneCmd)
}