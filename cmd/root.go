/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "deeployer",
	Short: "Deploy projects to remote servers via rsync",
	Long: `Deeployer is a deployment tool that executes build commands locally 
and syncs the output to remote servers using rsync over SSH. 

Configuration is managed through a TOML file at $XDG_CONFIG_HOME/deeployer/conf.toml
which defines projects, their build steps, and deployment targets.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here if needed
}


