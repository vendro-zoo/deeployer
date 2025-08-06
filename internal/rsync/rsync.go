package rsync

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Client struct {
	DryRun  bool
	Verbose bool
}

func New(dryRun, verbose bool) *Client {
	return &Client{
		DryRun:  dryRun,
		Verbose: verbose,
	}
}

func (c *Client) Sync(localPath, remoteUser, remoteHost, remotePath string, options []string) error {
	if err := c.validatePaths(localPath); err != nil {
		return err
	}

	args := c.buildRsyncArgs(localPath, remoteUser, remoteHost, remotePath, options)

	if c.Verbose || c.DryRun {
		fmt.Printf("Executing: rsync %s\n", strings.Join(args, " "))
	}

	if c.DryRun {
		return nil
	}

	cmd := exec.Command("rsync", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (c *Client) buildRsyncArgs(localPath, remoteUser, remoteHost, remotePath string, options []string) []string {
	args := make([]string, 0, len(options)+3)
	
	args = append(args, options...)

	if c.DryRun {
		args = append(args, "--dry-run")
	}

	if c.Verbose {
		args = append(args, "--verbose")
	}

	localPath = c.ensureTrailingSlash(localPath)
	args = append(args, localPath)

	remote := fmt.Sprintf("%s@%s:%s", remoteUser, remoteHost, remotePath)
	args = append(args, remote)

	return args
}

func (c *Client) validatePaths(localPath string) error {
	absPath, err := filepath.Abs(localPath)
	if err != nil {
		return fmt.Errorf("invalid local path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("local path does not exist: %s", absPath)
	}

	return nil
}

func (c *Client) ensureTrailingSlash(path string) string {
	if !strings.HasSuffix(path, "/") {
		return path + "/"
	}
	return path
}

func (c *Client) CheckRsyncAvailable() error {
	if c.DryRun {
		return nil
	}

	_, err := exec.LookPath("rsync")
	if err != nil {
		return fmt.Errorf("rsync not found in PATH: %w", err)
	}

	return nil
}