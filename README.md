# Deeployer

A deployment tool that executes commands locally and syncs output to remote servers via rsync over SSH.

## Overview

Deeployer allows you to configure build commands, deployment targets, and cleanup operations through a simple TOML configuration file. It supports multiple projects and multiple deployment targets per project.

## Configuration

Configuration is stored in `$XDG_CONFIG_HOME/deeployer/conf.toml` (typically `~/.config/deeployer/conf.toml`).

### Example Configuration

```toml
[projects.webapp]
path = "/home/user/projects/my-webapp"
build_commands = ["npm run build", "npm run test"]
output_dir = "./dist"
post_commands = ["rm -rf ./temp", "npm run clean"]
remotes = ["production", "staging"]

[projects.api]
path = "/home/user/projects/my-api"
build_commands = ["go build -o ./bin/api"]
output_dir = "./bin"
post_commands = ["go clean", "rm -f ./debug.log"]
remotes = ["production"]

[remotes.production]
host = "prod.example.com"
path = "/var/www/app"
user = "deploy"
rsync_options = ["-avz", "--delete"]
post_commands = ["sudo systemctl restart nginx"]

[remotes.staging]
host = "staging.example.com" 
path = "/var/www/app"
user = "deploy"
rsync_options = ["-avz"]
post_commands = []
```

## Deployment Flow

1. Change to the project's `path` directory
2. Execute project `build_commands` locally in that directory
3. Rsync `output_dir` from the project path to remote `path` 
4. Execute remote `post_commands` on the remote server via SSH
5. Execute project `post_commands` locally in the project directory for cleanup

## Usage

```bash
# Deploy a project to a specific remote
deeployer deploy webapp production

# Deploy to staging
deeployer deploy webapp staging

# Show available remotes for a project (when remote is omitted)
deeployer deploy webapp

# List configured projects and remotes
deeployer list

# Validate configuration
deeployer validate

# Dry run (show what would be executed)
deeployer deploy webapp production --dry-run

# Verbose output
deeployer deploy webapp production --verbose
```

## Implementation Plan

### Core Components

1. **Configuration System**
   - TOML parsing with validation
   - XDG config directory support
   - Project and remote configuration structs

2. **Execution Engine**
   - Local command execution
   - Rsync integration with SSH
   - Remote SSH command execution
   - Error handling and logging

3. **CLI Interface**
   - Deploy command with project selection
   - List command for viewing configuration
   - Validate command for config verification
   - Dry-run and verbose modes

### File Structure

```
cmd/
├── root.go           # Updated root command
├── deploy.go         # Deploy command implementation  
├── list.go          # List projects/remotes command
└── validate.go      # Config validation command

internal/
├── config/          # Configuration loading and validation
├── executor/        # Command execution logic
├── rsync/          # Rsync wrapper
└── ssh/            # SSH client for remote commands
```

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/BurntSushi/toml` - TOML configuration parsing
- `golang.org/x/crypto/ssh` - SSH client functionality