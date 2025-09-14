package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version string
	Commit  string
	Date    string
)

var rootCmd = &cobra.Command{
	Use:   "synodeploy",
	Short: "Deploy containers to Synology DSM 7.2+",
	Long: `SynoDeploy is a CLI tool that simplifies Docker container deployment
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

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(psCmd)
	rootCmd.AddCommand(rmCmd)
}
