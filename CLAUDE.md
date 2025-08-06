# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Deeployer is a CLI deployment tool built with Go and Cobra that executes build commands locally and syncs output to remote servers via rsync over SSH. Configuration is managed through TOML files following the XDG Base Directory specification.

## Build and Development Commands

```bash
# Build the application
go build

# Run the application (after building)
./deeployer --help

# Install dependencies
go mod tidy

# Build and test with example config
./deeployer validate  # requires config at ~/.config/deeployer/conf.toml
```

## Architecture Overview

The application follows a clean layered architecture:

**CLI Layer** (`cmd/`): Cobra commands that handle user interaction and orchestrate the deployment flow
- `deploy.go` - Main deployment orchestration, coordinates all components
- `list.go` - Configuration display utilities  
- `validate.go` - Configuration validation workflow

**Core Services** (`internal/`):
- `config/` - TOML configuration loading with XDG directory support and validation
- `executor/` - Local command execution with dry-run support
- `rsync/` - Rsync wrapper for file synchronization to remote servers
- `ssh/` - SSH client for remote command execution with key-based authentication

**Data Flow**: 
1. Config loaded from XDG directories and validated against schema
2. Deploy command coordinates: local build → rsync sync → remote SSH commands → local cleanup
3. All components support dry-run and verbose modes consistently

**Configuration Schema**: Two-section TOML with `[projects.*]` defining project paths, build pipelines and output directories, and `[remotes.*]` defining deployment targets. Projects reference remotes by name as an allowlist - users must explicitly specify which remote to deploy to. Each project must specify a `path` field that sets the working directory for all commands.

## Key Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/BurntSushi/toml` - Configuration parsing
- `golang.org/x/crypto/ssh` - SSH client with agent and key file support

## Configuration Location

The tool expects configuration at `$XDG_CONFIG_HOME/deeployer/conf.toml` (typically `~/.config/deeployer/conf.toml`). See `example-conf.toml` for reference structure.