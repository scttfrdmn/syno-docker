package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/scttfrdmn/synodeploy/pkg/config"
	"github.com/scttfrdmn/synodeploy/pkg/deploy"
	"github.com/scttfrdmn/synodeploy/pkg/synology"
)

var (
	deployProject string
	deployEnvFile string
)

var deployCmd = &cobra.Command{
	Use:   "deploy <compose-file>",
	Short: "Deploy from docker-compose.yml",
	Long: `Deploy containers from a docker-compose.yml file to your Synology NAS.
This command parses the compose file and creates individual containers for each service.`,
	Args: cobra.ExactArgs(1),
	RunE: deployCompose,
}

func deployCompose(cmd *cobra.Command, args []string) error {
	composeFile := args[0]

	// Check if compose file exists
	absPath, err := filepath.Abs(composeFile)
	if err != nil {
		return fmt.Errorf("failed to resolve compose file path: %w", err)
	}

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

	// Generate project name if not specified
	projectName := deployProject
	if projectName == "" {
		projectName = deploy.GenerateProjectName(absPath)
	}

	// Prepare deploy options
	opts := &deploy.ComposeOptions{
		ComposeFile: absPath,
		ProjectName: projectName,
		EnvFile:     deployEnvFile,
	}

	// Deploy compose
	if err := deploy.Compose(conn, opts); err != nil {
		return fmt.Errorf("deployment failed: %w", err)
	}

	fmt.Printf("\nYou can check the status with: synodeploy ps\n")
	return nil
}

func init() {
	deployCmd.Flags().StringVarP(&deployProject, "project", "p", "", "Project name (auto-generated from directory if not specified)")
	deployCmd.Flags().StringVar(&deployEnvFile, "env-file", "", "Environment file path")
}
