package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/synodeploy/pkg/config"
	"github.com/scttfrdmn/synodeploy/pkg/deploy"
	"github.com/scttfrdmn/synodeploy/pkg/synology"
)

var (
	runName       string
	runPorts      []string
	runVolumes    []string
	runEnv        []string
	runRestart    string
	runNetwork    string
	runWorkingDir string
	runUser       string
	runCommand    []string
)

var runCmd = &cobra.Command{
	Use:   "run <image>",
	Short: "Deploy a single container",
	Long: `Deploy a single Docker container to your Synology NAS.
This command pulls the specified image and creates a new container with the given configuration.`,
	Args: cobra.ExactArgs(1),
	RunE: runContainer,
}

func runContainer(cmd *cobra.Command, args []string) error {
	image := args[0]

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

	// Prepare container options
	opts := deploy.NewContainerOptions(image)
	opts.Name = runName
	opts.Ports = runPorts
	opts.Volumes = processVolumes(runVolumes, cfg.Defaults.VolumePath)
	opts.Env = runEnv
	opts.Restart = runRestart
	opts.NetworkMode = runNetwork
	opts.WorkingDir = runWorkingDir
	opts.User = runUser
	opts.Command = runCommand

	// Generate name if not provided
	if opts.Name == "" {
		opts.Name = generateContainerName(image)
	}

	// Deploy container
	containerID, err := deploy.Container(conn, opts)
	if err != nil {
		return fmt.Errorf("deployment failed: %w", err)
	}

	fmt.Printf("✅ Container deployed successfully!\n")
	fmt.Printf("Container ID: %s\n", containerID[:12])
	fmt.Printf("Container Name: %s\n", opts.Name)
	fmt.Printf("\nYou can check the status with: synodeploy ps\n")

	return nil
}

func processVolumes(volumes []string, defaultVolumePath string) []string {
	var processed []string
	for _, volume := range volumes {
		// If volume mapping starts with ./ or doesn't start with /, prepend default path
		if strings.HasPrefix(volume, "./") || (!strings.HasPrefix(volume, "/") && strings.Contains(volume, ":")) {
			parts := strings.SplitN(volume, ":", 2)
			if len(parts) == 2 {
				hostPath := parts[0]
				hostPath = strings.TrimPrefix(hostPath, "./")
				if !strings.HasPrefix(hostPath, "/") {
					hostPath = fmt.Sprintf("%s/%s", defaultVolumePath, hostPath)
				}
				volume = fmt.Sprintf("%s:%s", hostPath, parts[1])
			}
		}
		processed = append(processed, volume)
	}
	return processed
}

func generateContainerName(image string) string {
	// Extract image name without registry and tag
	parts := strings.Split(image, "/")
	name := parts[len(parts)-1]

	// Remove tag if present
	if idx := strings.Index(name, ":"); idx != -1 {
		name = name[:idx]
	}

	return name
}

func init() {
	runCmd.Flags().StringVarP(&runName, "name", "n", "", "Container name (auto-generated if not specified)")
	runCmd.Flags().StringSliceVarP(&runPorts, "port", "p", []string{}, "Port mappings (format: host:container)")
	runCmd.Flags().StringSliceVarP(&runVolumes, "volume", "v", []string{}, "Volume mappings (format: host:container)")
	runCmd.Flags().StringSliceVarP(&runEnv, "env", "e", []string{}, "Environment variables (format: KEY=value)")
	runCmd.Flags().StringVar(&runRestart, "restart", synology.DefaultRestartPolicy, "Restart policy (no, always, unless-stopped, on-failure)")
	runCmd.Flags().StringVar(&runNetwork, "network", synology.DefaultNetwork, "Network mode")
	runCmd.Flags().StringVarP(&runWorkingDir, "workdir", "w", "", "Working directory inside container")
	runCmd.Flags().StringVarP(&runUser, "user", "u", "", "User to run container as (format: uid:gid)")
	runCmd.Flags().StringSliceVar(&runCommand, "command", []string{}, "Command to run in container")
}
