package cmd

import (
	"fmt"
	"sort"
	"strings"

	"deeployer/internal/config"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured projects and remotes",
	Long:  `Display all configured projects and their associated remotes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		listProjects(cfg)
		fmt.Println()
		listRemotes(cfg)

		return nil
	},
}

func listProjects(cfg *config.Config) {
	fmt.Println("Projects:")
	fmt.Println("=========")

	if len(cfg.Projects) == 0 {
		fmt.Println("No projects configured")
		return
	}

	var projectNames []string
	for name := range cfg.Projects {
		projectNames = append(projectNames, name)
	}
	sort.Strings(projectNames)

	for _, name := range projectNames {
		project := cfg.Projects[name]
		fmt.Printf("\n%s:\n", name)
		fmt.Printf("  Path: %s\n", project.Path)
		fmt.Printf("  Output Directory: %s\n", project.OutputDir)
		fmt.Printf("  Build Commands: %s\n", formatCommands(project.BuildCommands))
		if len(project.PostCommands) > 0 {
			fmt.Printf("  Post Commands: %s\n", formatCommands(project.PostCommands))
		}
		fmt.Printf("  Remotes: %s\n", strings.Join(project.Remotes, ", "))
	}
}

func listRemotes(cfg *config.Config) {
	fmt.Println("Remotes:")
	fmt.Println("========")

	if len(cfg.Remotes) == 0 {
		fmt.Println("No remotes configured")
		return
	}

	var remoteNames []string
	for name := range cfg.Remotes {
		remoteNames = append(remoteNames, name)
	}
	sort.Strings(remoteNames)

	for _, name := range remoteNames {
		remote := cfg.Remotes[name]
		fmt.Printf("\n%s:\n", name)
		fmt.Printf("  Host: %s@%s\n", remote.User, remote.Host)
		fmt.Printf("  Path: %s\n", remote.Path)
		fmt.Printf("  Rsync Options: %s\n", strings.Join(remote.RsyncOptions, " "))
		if len(remote.PostCommands) > 0 {
			fmt.Printf("  Post Commands: %s\n", formatCommands(remote.PostCommands))
		}
	}
}

func formatCommands(commands []string) string {
	if len(commands) == 0 {
		return "(none)"
	}
	if len(commands) == 1 {
		return commands[0]
	}
	return fmt.Sprintf("[%s]", strings.Join(commands, ", "))
}

func init() {
	rootCmd.AddCommand(listCmd)
}