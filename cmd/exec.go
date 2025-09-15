package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/syno-docker/pkg/config"
	"github.com/scttfrdmn/syno-docker/pkg/deploy"
	"github.com/scttfrdmn/syno-docker/pkg/synology"
)

var (
	execInteractive bool
	execTTY         bool
	execUser        string
	execWorkdir     string
	execEnv         []string
)

var execCmd = &cobra.Command{
	Use:   "exec [OPTIONS] <container> <command> [args...]",
	Short: "Execute a command in a running container",
	Long:  `Execute a command in a running container on your Synology NAS.`,
	Args:  cobra.MinimumNArgs(2),
	RunE:  executeCommand,
}

func executeCommand(cmd *cobra.Command, args []string) error {
	containerNameOrID := args[0]
	command := args[1:]

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

	// Execute command in container
	opts := &deploy.ExecOptions{
		Interactive: execInteractive,
		TTY:         execTTY,
		User:        execUser,
		WorkingDir:  execWorkdir,
		Env:         execEnv,
	}

	if execInteractive {
		fmt.Printf("Executing interactive command in container %s...\n", containerNameOrID)
		return deploy.ExecInteractive(conn, containerNameOrID, command, opts)
	} else {
		output, err := deploy.ExecCommand(conn, containerNameOrID, command, opts)
		if err != nil {
			return fmt.Errorf("failed to execute command: %w", err)
		}
		fmt.Print(output)
		return nil
	}
}

func init() {
	execCmd.Flags().BoolVarP(&execInteractive, "interactive", "i", false, "Keep STDIN open even if not attached")
	execCmd.Flags().BoolVarP(&execTTY, "tty", "t", false, "Allocate a pseudo-TTY")
	execCmd.Flags().StringVarP(&execUser, "user", "u", "", "Username or UID (format: <name|uid>[:<group|gid>])")
	execCmd.Flags().StringVarP(&execWorkdir, "workdir", "w", "", "Working directory inside the container")
	execCmd.Flags().StringSliceVarP(&execEnv, "env", "e", []string{}, "Set environment variables")
}
