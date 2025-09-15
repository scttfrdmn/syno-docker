package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/syno-docker/pkg/config"
	"github.com/scttfrdmn/syno-docker/pkg/deploy"
	"github.com/scttfrdmn/syno-docker/pkg/synology"
)

var (
	logsFollow     bool
	logsTail       string
	logsSince      string
	logsTimestamps bool
)

var logsCmd = &cobra.Command{
	Use:   "logs <container>",
	Short: "Fetch the logs of a container",
	Long:  `Fetch the logs of a container from your Synology NAS.`,
	Args:  cobra.ExactArgs(1),
	RunE:  showContainerLogs,
}

func showContainerLogs(cmd *cobra.Command, args []string) error {
	containerNameOrID := args[0]

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

	// Show container logs
	if logsFollow {
		fmt.Printf("Following logs for container %s (press Ctrl+C to stop)...\n", containerNameOrID)
		return deploy.FollowContainerLogs(conn, containerNameOrID, logsTail, logsSince, logsTimestamps, os.Stdout, os.Stderr)
	} else {
		fmt.Printf("Fetching logs for container %s...\n", containerNameOrID)
		logs, err := deploy.GetContainerLogs(conn, containerNameOrID, logsTail, logsSince, logsTimestamps)
		if err != nil {
			return fmt.Errorf("failed to get container logs: %w", err)
		}
		fmt.Print(logs)
		return nil
	}
}

func init() {
	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Follow log output")
	logsCmd.Flags().StringVar(&logsTail, "tail", "all", "Number of lines to show from the end of the logs")
	logsCmd.Flags().StringVar(&logsSince, "since", "", "Show logs since timestamp (e.g. 2013-01-02T13:23:37Z) or relative (e.g. 42m for 42 minutes)")
	logsCmd.Flags().BoolVarP(&logsTimestamps, "timestamps", "t", false, "Show timestamps")
}
