# CLI Tool Design

> Status: 🟡 In Design

## Overview

`showbiz` is a command-line tool for developers to manage organizations, projects, and cloud resources from their terminal. Built in Go using **Cobra**, with the Go SDK.

## Features

- Built on [Cobra](https://github.com/spf13/cobra) CLI framework
- **Shell completion** generation for Bash, Zsh, Fish, and PowerShell (`showbiz completion <shell>`)
- No plugin system — all functionality is built-in

## Command Structure

```
showbiz
├── auth
│   ├── login          # Login with username/password, receive JWT
│   ├── register       # Register a new account
│   ├── logout         # Clear stored credentials
│   └── status         # Show current auth state
├── org
│   ├── list           # List organizations
│   ├── create         # Create new org
│   ├── get            # Show org details
│   ├── update         # Update org
│   ├── deactivate     # Deactivate org (deletes all projects/resources)
│   ├── activate       # Reactivate org
│   ├── members list   # List org members
│   ├── members add    # Add member to org
│   └── members remove # Remove member from org
├── project
│   ├── list           # List projects in an org
│   ├── create         # Create a project
│   ├── get            # Show project details
│   ├── update         # Update project
│   └── delete         # Delete project and all its resources
├── resource
│   ├── list           # List resources in a project
│   ├── create         # Create a resource (machine, network) via a connection
│   ├── get            # Show resource details
│   ├── update         # Update a resource
│   └── delete         # Delete a resource
├── connection
│   ├── list           # List connections in a project
│   ├── create         # Create a connection to a provider account
│   ├── get            # Show connection details
│   ├── update         # Update connection (credentials, config)
│   └── delete         # Delete a connection
├── iam
│   ├── policy list     # List policies (global + org)
│   ├── policy get      # Show policy details
│   ├── policy create   # Create an org-level policy
│   ├── policy update   # Update policy permissions
│   ├── policy delete   # Delete an org policy
│   ├── attach          # Attach a policy to a user on a project
│   ├── attachments     # List policy attachments on a project
│   └── detach          # Remove a policy attachment
├── provider
│   ├── list           # List available providers (read-only)
│   └── get            # Show provider details (read-only)
└── config
    ├── set            # Set a config value
    └── get            # Get a config value
```

## Configuration

- Config file: `~/.showbiz/config.yaml` (stores API URL, active org)
- Credentials: `~/.showbiz/credentials.json` (stores JWT + refresh token)
- Environment variables: `SHOWBIZ_API_URL`, `SHOWBIZ_USERNAME`, `SHOWBIZ_PASSWORD`

## Usage Examples

```bash
# Login
showbiz auth login --username user@example.com --password ****

# Create a project
showbiz project create --org org_123 --name "my-app"

# Create a connection to an AWS account
showbiz connection create --project proj_123 \
  --name "AWS-1234" \
  --provider aws \
  --credentials '{"accessKeyId":"AKIA...","secretAccessKey":"..."}' \
  --config '{"accountId":"123456789012","defaultRegion":"us-east-1"}'

# Create a machine resource via the connection
showbiz resource create --project proj_123 \
  --connection conn_101 \
  --type machine \
  --name web-1 \
  --values '{"size":"medium","region":"us-east","image":"ubuntu-22.04"}'

# Set IAM policy
showbiz iam policy create --org org_123 \
  --name "DevOps" \
  --permissions "resource:create,resource:read,resource:update,resource:delete,connection:create,connection:read"

# Attach policy to user on a project
showbiz iam attach --project proj_123 --user dev@example.com --policy policy_101

# List attachments on a project
showbiz iam attachments --project proj_123

# List resources
showbiz resource list --project proj_123 --output json

# List connections
showbiz connection list --project proj_123
```

## UX Principles

- Colorized output with `--no-color` flag
- JSON output mode with `--output json` for scripting
- Interactive prompts where it makes sense, with `--yes` flag to skip
- Helpful error messages with suggestions

## Open Questions

None — all decisions resolved for initial design.
