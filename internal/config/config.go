package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Projects map[string]Project `toml:"projects"`
	Remotes  map[string]Remote  `toml:"remotes"`
}

type Project struct {
	Path          string   `toml:"path"`
	BuildCommands []string `toml:"build_commands"`
	OutputDir     string   `toml:"output_dir"`
	PostCommands  []string `toml:"post_commands"`
	Remotes       []string `toml:"remotes"`
}

type Remote struct {
	Host         string   `toml:"host"`
	Path         string   `toml:"path"`
	User         string   `toml:"user"`
	RsyncOptions []string `toml:"rsync_options"`
	PostCommands []string `toml:"post_commands"`
}

func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found at %s", configPath)
	}

	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if len(c.Projects) == 0 {
		return fmt.Errorf("no projects defined")
	}

	for name, project := range c.Projects {
		if err := project.Validate(); err != nil {
			return fmt.Errorf("project %s: %w", name, err)
		}

		for _, remoteName := range project.Remotes {
			if _, exists := c.Remotes[remoteName]; !exists {
				return fmt.Errorf("project %s references unknown remote: %s", name, remoteName)
			}
		}
	}

	for name, remote := range c.Remotes {
		if err := remote.Validate(); err != nil {
			return fmt.Errorf("remote %s: %w", name, err)
		}
	}

	return nil
}

func (p *Project) Validate() error {
	if p.Path == "" {
		return fmt.Errorf("project path not specified")
	}

	// Validate project path exists and is accessible
	absPath, err := filepath.Abs(p.Path)
	if err != nil {
		return fmt.Errorf("invalid project path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", absPath)
	} else if err != nil {
		return fmt.Errorf("cannot access project path: %w", err)
	}

	if len(p.BuildCommands) == 0 {
		return fmt.Errorf("no build commands defined")
	}

	if p.OutputDir == "" {
		return fmt.Errorf("output directory not specified")
	}

	if len(p.Remotes) == 0 {
		return fmt.Errorf("no remotes specified")
	}

	return nil
}

func (r *Remote) Validate() error {
	if r.Host == "" {
		return fmt.Errorf("host not specified")
	}

	if r.Path == "" {
		return fmt.Errorf("path not specified")
	}

	if r.User == "" {
		return fmt.Errorf("user not specified")
	}

	if len(r.RsyncOptions) == 0 {
		r.RsyncOptions = []string{"-avz"}
	}

	return nil
}

func getConfigPath() (string, error) {
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		xdgConfig = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(xdgConfig, "deeployer", "conf.toml"), nil
}

func GetConfigPath() (string, error) {
	return getConfigPath()
}