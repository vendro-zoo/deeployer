package cmd

import (
	"fmt"

	"deeployer/internal/config"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration file",
	Long:  `Check the configuration file for syntax errors and validate all projects and remotes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := config.GetConfigPath()
		if err != nil {
			return fmt.Errorf("failed to get config path: %w", err)
		}

		fmt.Printf("Validating configuration file: %s\n", configPath)

		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("✗ Configuration validation failed: %v\n", err)
			return err
		}

		fmt.Println("✓ Configuration is valid")
		
		fmt.Printf("✓ Found %d project(s)\n", len(cfg.Projects))
		fmt.Printf("✓ Found %d remote(s)\n", len(cfg.Remotes))

		for projectName, project := range cfg.Projects {
			fmt.Printf("✓ Project '%s': %d build command(s), %d post command(s), %d remote(s)\n", 
				projectName, len(project.BuildCommands), len(project.PostCommands), len(project.Remotes))
		}

		for remoteName := range cfg.Remotes {
			fmt.Printf("✓ Remote '%s' configured\n", remoteName)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}