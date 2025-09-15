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

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage networks",
	Long:  `Manage Docker networks on your Synology NAS.`,
}

var networkListCmd = &cobra.Command{
	Use:     "ls [OPTIONS]",
	Short:   "List networks",
	Long:    `List Docker networks on your Synology NAS.`,
	RunE:    listNetworks,
	Aliases: []string{"list"},
}

var networkCreateCmd = &cobra.Command{
	Use:   "create [OPTIONS] NETWORK",
	Short: "Create a network",
	Long:  `Create a Docker network on your Synology NAS.`,
	Args:  cobra.ExactArgs(1),
	RunE:  createNetwork,
}

var networkRemoveCmd = &cobra.Command{
	Use:     "rm NETWORK [NETWORK...]",
	Short:   "Remove one or more networks",
	Long:    `Remove one or more Docker networks from your Synology NAS.`,
	Args:    cobra.MinimumNArgs(1),
	RunE:    removeNetworks,
	Aliases: []string{"remove"},
}

var networkInspectCmd = &cobra.Command{
	Use:   "inspect NETWORK [NETWORK...]",
	Short: "Display detailed information on one or more networks",
	Long:  `Display detailed information on one or more Docker networks.`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  inspectNetworks,
}

var networkConnectCmd = &cobra.Command{
	Use:   "connect [OPTIONS] NETWORK CONTAINER",
	Short: "Connect a container to a network",
	Long:  `Connect a container to a Docker network.`,
	Args:  cobra.ExactArgs(2),
	RunE:  connectToNetwork,
}

var networkDisconnectCmd = &cobra.Command{
	Use:   "disconnect [OPTIONS] NETWORK CONTAINER",
	Short: "Disconnect a container from a network",
	Long:  `Disconnect a container from a Docker network.`,
	Args:  cobra.ExactArgs(2),
	RunE:  disconnectFromNetwork,
}

var networkPruneCmd = &cobra.Command{
	Use:   "prune [OPTIONS]",
	Short: "Remove all unused networks",
	Long:  `Remove all unused Docker networks.`,
	RunE:  pruneNetworks,
}

var (
	// List options
	networkListFormat string
	networkListQuiet  bool
	networkListFilter []string

	// Create options
	networkCreateDriver     string
	networkCreateDriverOpts []string
	networkCreateGateway    []string
	networkCreateIPRange    []string
	networkCreateIPAM       []string
	networkCreateSubnet     []string
	networkCreateLabel      []string
	networkCreateAttachable bool
	networkCreateIngress    bool
	networkCreateInternal   bool
	networkCreateIPv6       bool

	// Connect options
	networkConnectAlias     []string
	networkConnectIP        string
	networkConnectIP6       string
	networkConnectLinkLocal []string

	// Disconnect options
	networkDisconnectForce bool

	// Inspect options
	networkInspectFormat string

	// Prune options
	networkPruneForce  bool
	networkPruneFilter []string
)

func listNetworks(cmd *cobra.Command, args []string) error {
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

	// List networks
	opts := &deploy.NetworkListOptions{
		Format: networkListFormat,
		Quiet:  networkListQuiet,
		Filter: networkListFilter,
	}

	if networkListQuiet {
		networkIDs, err := deploy.ListNetworkIDs(conn, opts)
		if err != nil {
			return fmt.Errorf("failed to list networks: %w", err)
		}
		for _, id := range networkIDs {
			fmt.Println(id)
		}
		return nil
	}

	networks, err := deploy.ListNetworks(conn, opts)
	if err != nil {
		return fmt.Errorf("failed to list networks: %w", err)
	}

	if len(networks) == 0 {
		fmt.Println("No networks found.")
		return nil
	}

	// Display networks in table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NETWORK ID\tNAME\tDRIVER\tSCOPE")

	for _, network := range networks {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			network.ID, network.Name, network.Driver, network.Scope)
	}

	return w.Flush()
}

func createNetwork(cmd *cobra.Command, args []string) error {
	networkName := args[0]

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

	// Create network
	opts := &deploy.NetworkCreateOptions{
		Driver:     networkCreateDriver,
		DriverOpts: networkCreateDriverOpts,
		Gateway:    networkCreateGateway,
		IPRange:    networkCreateIPRange,
		IPAM:       networkCreateIPAM,
		Subnet:     networkCreateSubnet,
		Labels:     networkCreateLabel,
		Attachable: networkCreateAttachable,
		Ingress:    networkCreateIngress,
		Internal:   networkCreateInternal,
		IPv6:       networkCreateIPv6,
	}

	fmt.Printf("Creating network %s...\n", networkName)
	networkID, err := deploy.CreateNetwork(conn, networkName, opts)
	if err != nil {
		return fmt.Errorf("failed to create network: %w", err)
	}

	fmt.Printf("✅ Network %s created successfully! (ID: %s)\n", networkName, networkID)
	return nil
}

func removeNetworks(cmd *cobra.Command, args []string) error {
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

	// Remove each network
	for _, networkName := range args {
		fmt.Printf("Removing network %s...\n", networkName)
		if err := deploy.RemoveNetwork(conn, networkName); err != nil {
			return fmt.Errorf("failed to remove network %s: %w", networkName, err)
		}
		fmt.Printf("✅ Network %s removed successfully!\n", networkName)
	}

	return nil
}

func inspectNetworks(cmd *cobra.Command, args []string) error {
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

	// Inspect networks
	opts := &deploy.NetworkInspectOptions{
		Format: networkInspectFormat,
	}

	for _, networkName := range args {
		info, err := deploy.InspectNetwork(conn, networkName, opts)
		if err != nil {
			return fmt.Errorf("failed to inspect network %s: %w", networkName, err)
		}
		fmt.Print(info)
	}

	return nil
}

func connectToNetwork(cmd *cobra.Command, args []string) error {
	networkName := args[0]
	containerName := args[1]

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

	// Connect container to network
	opts := &deploy.NetworkConnectOptions{
		Alias:     networkConnectAlias,
		IP:        networkConnectIP,
		IPv6:      networkConnectIP6,
		LinkLocal: networkConnectLinkLocal,
	}

	fmt.Printf("Connecting container %s to network %s...\n", containerName, networkName)
	if err := deploy.ConnectContainerToNetwork(conn, networkName, containerName, opts); err != nil {
		return fmt.Errorf("failed to connect container to network: %w", err)
	}

	fmt.Printf("✅ Container %s connected to network %s successfully!\n", containerName, networkName)
	return nil
}

func disconnectFromNetwork(cmd *cobra.Command, args []string) error {
	networkName := args[0]
	containerName := args[1]

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

	// Disconnect container from network
	opts := &deploy.NetworkDisconnectOptions{
		Force: networkDisconnectForce,
	}

	fmt.Printf("Disconnecting container %s from network %s...\n", containerName, networkName)
	if err := deploy.DisconnectContainerFromNetwork(conn, networkName, containerName, opts); err != nil {
		return fmt.Errorf("failed to disconnect container from network: %w", err)
	}

	fmt.Printf("✅ Container %s disconnected from network %s successfully!\n", containerName, networkName)
	return nil
}

func pruneNetworks(cmd *cobra.Command, args []string) error {
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
	if !networkPruneForce {
		fmt.Print("WARNING! This will remove all networks not used by at least one container.\n")
		fmt.Print("Are you sure you want to continue? [y/N] ")

		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Prune networks
	opts := &deploy.NetworkPruneOptions{
		Force:  networkPruneForce,
		Filter: networkPruneFilter,
	}

	result, err := deploy.PruneNetworks(conn, opts)
	if err != nil {
		return fmt.Errorf("failed to prune networks: %w", err)
	}

	fmt.Printf("Deleted Networks: %d\n", result.NetworksDeleted)
	fmt.Printf("Total reclaimed space: %s\n", result.SpaceReclaimed)

	return nil
}

func init() {
	// network list command
	networkListCmd.Flags().StringVar(&networkListFormat, "format", "", "Pretty-print networks using a Go template")
	networkListCmd.Flags().BoolVarP(&networkListQuiet, "quiet", "q", false, "Only display network IDs")
	networkListCmd.Flags().StringSliceVarP(&networkListFilter, "filter", "f", []string{}, "Provide filter values (e.g. 'driver=bridge')")

	// network create command
	networkCreateCmd.Flags().StringVarP(&networkCreateDriver, "driver", "d", "bridge", "Driver to manage the Network")
	networkCreateCmd.Flags().StringSliceVar(&networkCreateDriverOpts, "opt", []string{}, "Set driver specific options")
	networkCreateCmd.Flags().StringSliceVar(&networkCreateGateway, "gateway", []string{}, "IPv4 or IPv6 Gateway for the master subnet")
	networkCreateCmd.Flags().StringSliceVar(&networkCreateIPRange, "ip-range", []string{}, "Allocate container ip from a sub-range")
	networkCreateCmd.Flags().StringSliceVar(&networkCreateIPAM, "ipam-driver", []string{}, "IP Address Management Driver")
	networkCreateCmd.Flags().StringSliceVar(&networkCreateSubnet, "subnet", []string{}, "Subnet in CIDR format that represents a network segment")
	networkCreateCmd.Flags().StringSliceVar(&networkCreateLabel, "label", []string{}, "Set metadata on a network")
	networkCreateCmd.Flags().BoolVar(&networkCreateAttachable, "attachable", false, "Enable manual container attachment")
	networkCreateCmd.Flags().BoolVar(&networkCreateIngress, "ingress", false, "Create swarm routing-mesh network")
	networkCreateCmd.Flags().BoolVar(&networkCreateInternal, "internal", false, "Restrict external access to the network")
	networkCreateCmd.Flags().BoolVar(&networkCreateIPv6, "ipv6", false, "Enable IPv6 networking")

	// network connect command
	networkConnectCmd.Flags().StringSliceVar(&networkConnectAlias, "alias", []string{}, "Add network-scoped alias for the container")
	networkConnectCmd.Flags().StringVar(&networkConnectIP, "ip", "", "IPv4 address (e.g., 172.30.100.104)")
	networkConnectCmd.Flags().StringVar(&networkConnectIP6, "ip6", "", "IPv6 address (e.g., 2001:db8::33)")
	networkConnectCmd.Flags().StringSliceVar(&networkConnectLinkLocal, "link-local", []string{}, "Container IPv4/IPv6 link-local addresses")

	// network disconnect command
	networkDisconnectCmd.Flags().BoolVarP(&networkDisconnectForce, "force", "f", false, "Force the container to disconnect from a network")

	// network inspect command
	networkInspectCmd.Flags().StringVarP(&networkInspectFormat, "format", "f", "", "Format the output using the given Go template")

	// network prune command
	networkPruneCmd.Flags().BoolVarP(&networkPruneForce, "force", "f", false, "Do not prompt for confirmation")
	networkPruneCmd.Flags().StringSliceVar(&networkPruneFilter, "filter", []string{}, "Provide filter values")

	// Add subcommands to network
	networkCmd.AddCommand(networkListCmd)
	networkCmd.AddCommand(networkCreateCmd)
	networkCmd.AddCommand(networkRemoveCmd)
	networkCmd.AddCommand(networkInspectCmd)
	networkCmd.AddCommand(networkConnectCmd)
	networkCmd.AddCommand(networkDisconnectCmd)
	networkCmd.AddCommand(networkPruneCmd)
}
