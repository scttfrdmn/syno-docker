package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the current version of syno-docker
	Version string
	// Commit is the git commit hash
	Commit string
	// Date is the build date
	Date string
)

var rootCmd = &cobra.Command{
	Use:   "syno-docker",
	Short: "Deploy containers to Synology DSM 7.2+",
	Long: `syno-docker is a CLI tool that simplifies Docker container deployment
to Synology NAS devices running DSM 7.2+. It handles SSH connection management,
Docker client setup, and path resolution issues specific to Synology Container Manager.`,
	Version: getVersion(),
}

func getVersion() string {
	if Version == "" {
		return "dev"
	}
	return fmt.Sprintf("%s (commit %s, built %s)", Version, Commit, Date)
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(psCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(execCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(imagesCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(rmiCmd)
	rootCmd.AddCommand(systemCmd)
	rootCmd.AddCommand(volumeCmd)
	rootCmd.AddCommand(inspectCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
}
