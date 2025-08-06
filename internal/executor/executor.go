package executor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Executor struct {
	DryRun  bool
	Verbose bool
}

func New(dryRun, verbose bool) *Executor {
	return &Executor{
		DryRun:  dryRun,
		Verbose: verbose,
	}
}

func (e *Executor) ExecuteCommands(commands []string, workDir string) error {
	if len(commands) == 0 {
		return nil
	}

	for _, command := range commands {
		if err := e.executeCommand(command, workDir); err != nil {
			return fmt.Errorf("command failed: %s: %w", command, err)
		}
	}

	return nil
}

func (e *Executor) executeCommand(command, workDir string) error {
	if e.Verbose || e.DryRun {
		fmt.Printf("Executing: %s\n", command)
	}

	if e.DryRun {
		return nil
	}

	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = workDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (e *Executor) CheckOutputDir(outputDir string) error {
	if e.Verbose {
		fmt.Printf("Checking output directory: %s\n", outputDir)
	}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return fmt.Errorf("output directory does not exist: %s", outputDir)
	}

	return nil
}