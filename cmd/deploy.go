package cmd

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"deeployer/internal/config"
	"deeployer/internal/executor"
	"deeployer/internal/rsync"
	"deeployer/internal/ssh"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	dryRun  bool
	verbose bool
)

var deployCmd = &cobra.Command{
	Use:   "deploy [project] [remote]",
	Short: "Deploy a project to a specific remote",
	Long: `Deploy a project by executing its build commands, syncing the output directory 
to the specified remote server via rsync, and running post-deployment commands.

The remote must be in the project's allowed remotes list.`,
	Args: cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		var projectName string
		if len(args) <= 0 {
			options := make([]huh.Option[string], 0, len(cfg.Projects))
			for name := range cfg.Projects {
				options = append(options, huh.NewOption(name, name))
			}

			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Project").
						Options(options...).
						Value(&projectName),
				),
			)

			err := form.Run()
			if err != nil {
				return fmt.Errorf("failed to select project: %w", err)
			}

			if projectName == "" {
				return fmt.Errorf("no project selected")
			}
		} else {
			projectName = args[0]
		}

		project, exists := cfg.Projects[projectName]
		if !exists {
			return fmt.Errorf("project '%s' not found in configuration", projectName)
		}

		var remoteName string
		if len(args) <= 1 {
			options := make([]huh.Option[string], 0, len(project.Remotes))
			for _, remote := range project.Remotes {
				options = append(options, huh.NewOption(remote, remote))
			}

			form := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[string]().
						Title("Project").
						Options(options...).
						Value(&remoteName),
				),
			)

			err := form.Run()
			if err != nil {
				return fmt.Errorf("failed to select remote: %w", err)
			}

			if remoteName == "" {
				return fmt.Errorf("no remote selected")
			}
		} else {
			remoteName = args[1]
		}

		return deployProject(cfg, projectName, project, remoteName)
	},
}

func deployProject(cfg *config.Config, projectName string, project config.Project, remoteName string) error {
	// Validate that the remote is allowed for this project
	remoteAllowed := slices.Contains(project.Remotes, remoteName)

	if !remoteAllowed {
		return fmt.Errorf("remote '%s' is not allowed for project '%s'. Available remotes: %s",
			remoteName, projectName, strings.Join(project.Remotes, ", "))
	}

	// Check if remote exists in configuration
	remote, exists := cfg.Remotes[remoteName]
	if !exists {
		return fmt.Errorf("remote '%s' not found in configuration", remoteName)
	}

	exec := executor.New(dryRun, verbose)
	rsyncClient := rsync.New(dryRun, verbose)
	sshClient := ssh.New(dryRun, verbose)

	if verbose {
		fmt.Printf("Deploying project: %s (path: %s) to remote: %s\n", projectName, project.Path, remoteName)
	}

	if err := rsyncClient.CheckRsyncAvailable(); err != nil {
		return fmt.Errorf("rsync check failed: %w", err)
	}

	if verbose {
		fmt.Println("Executing build commands...")
	}
	if err := exec.ExecuteCommands(project.BuildCommands, project.Path); err != nil {
		return fmt.Errorf("build commands failed: %w", err)
	}

	// Validate output directory doesn't contain directory traversal
	if strings.Contains(project.OutputDir, "..") {
		return fmt.Errorf("output directory contains directory traversal: %s", project.OutputDir)
	}

	outputPath := filepath.Join(project.Path, project.OutputDir)
	// Clean the path to resolve any remaining . or .. elements
	outputPath = filepath.Clean(outputPath)

	// Ensure the cleaned path is still within the project directory
	projectAbsPath, err := filepath.Abs(project.Path)
	if err != nil {
		return fmt.Errorf("failed to get absolute project path: %w", err)
	}

	outputAbsPath, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute output path: %w", err)
	}

	if !strings.HasPrefix(outputAbsPath, projectAbsPath) {
		return fmt.Errorf("output directory is outside project directory: %s", outputAbsPath)
	}

	if err := exec.CheckOutputDir(outputPath); err != nil {
		return fmt.Errorf("output directory check failed: %w", err)
	}

	if verbose {
		fmt.Printf("Syncing to remote: %s\n", remoteName)
	}
	if err := rsyncClient.Sync(outputPath, remote.User, remote.Host, remote.Path, remote.RsyncOptions); err != nil {
		return fmt.Errorf("rsync to %s failed: %w", remoteName, err)
	}

	if len(remote.PostCommands) > 0 {
		if verbose {
			fmt.Printf("Executing post commands on remote: %s\n", remoteName)
		}
		if err := sshClient.ExecuteCommands(remote.Host, remote.User, remote.PostCommands); err != nil {
			return fmt.Errorf("remote post commands failed on %s: %w", remoteName, err)
		}
	}

	if len(project.PostCommands) > 0 {
		if verbose {
			fmt.Println("Executing local post commands...")
		}
		if err := exec.ExecuteCommands(project.PostCommands, project.Path); err != nil {
			return fmt.Errorf("local post commands failed: %w", err)
		}
	}

	fmt.Printf("Successfully deployed %s to %s\n", projectName, remoteName)
	return nil
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be executed without making changes")
	deployCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}
