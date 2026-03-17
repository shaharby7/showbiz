# Terraform Provider Design

> Status: 🟢 Implemented

## Overview

The Showbiz Terraform provider allows teams to manage Showbiz resources declaratively as infrastructure-as-code. Built in Go using the **Terraform Plugin Framework** and the Showbiz Go SDK.

## Authentication

```hcl
provider "showbiz" {
  api_url  = "https://api.showbiz.dev"
  username = var.showbiz_username
  password = var.showbiz_password
}
```

The provider authenticates via username/password, obtains a JWT internally, and handles token refresh automatically.

## Resources

```hcl
resource "showbiz_project" "my_app" {
  organization_id = "org_123"
  name            = "my-app"
}

resource "showbiz_connection" "aws_prod" {
  project_id = showbiz_project.my_app.id
  name       = "AWS-1234"
  provider   = "aws"

  credentials = {
    access_key_id     = var.aws_access_key
    secret_access_key = var.aws_secret_key
  }

  config = {
    account_id     = "123456789012"
    default_region = "us-east-1"
  }
}

resource "showbiz_resource" "web_server" {
  project_id    = showbiz_project.my_app.id
  connection_id = showbiz_connection.aws_prod.id
  name          = "web-server-1"
  resource_type = "machine"

  values = {
    size   = "medium"
    region = "us-east"
    image  = "ubuntu-22.04"
  }
}

resource "showbiz_resource" "app_network" {
  project_id    = showbiz_project.my_app.id
  connection_id = showbiz_connection.aws_prod.id
  name          = "app-vpc"
  resource_type = "network"

  values = {
    name = "app-vpc"
    cidr = "10.0.0.0/16"
  }
}

resource "showbiz_iam_policy" "devops" {
  organization_id = "org_123"
  name            = "DevOps"
  permissions     = [
    "resource:create", "resource:read", "resource:update", "resource:delete",
    "connection:create", "connection:read"
  ]
}

resource "showbiz_policy_attachment" "dev_access" {
  project_id = showbiz_project.my_app.id
  user_id    = "dev@example.com"
  policy_id  = showbiz_iam_policy.devops.id
}
```

## Planned Resources

| Resource | Description |
|---|---|
| `showbiz_project` | Manage projects within an organization |
| `showbiz_connection` | Manage connections to provider accounts within a project |
| `showbiz_resource` | Manage resources (machine, network) via a connection |
| `showbiz_iam_policy` | Manage IAM policies at the organization level |
| `showbiz_policy_attachment` | Attach a policy to a user on a project |

## Data Sources

| Data Source | Description |
|---|---|
| `showbiz_project` | Look up an existing project |
| `showbiz_connection` | Look up an existing connection |
| `showbiz_resource` | Look up an existing resource |
| `showbiz_provider` | Look up available cloud providers |

## Import Support

All resources support `terraform import` for adopting existing Showbiz resources into Terraform state. Import uses the entity's ID:

```bash
terraform import showbiz_project.my_app proj_123
terraform import showbiz_connection.aws_prod conn_101
terraform import showbiz_resource.web_server "sbz:machine:proj_456:AWS-1234:web-server-1"
terraform import showbiz_iam_policy.dev_access policy_789
```

## Acceptance Testing

Acceptance tests run against a **real Showbiz API instance** (not mocked). Strategy:

- Tests use a dedicated test organization and project, created in `TestMain` setup and torn down after
- Each test creates its own resources (connections, resources, IAM policies) and verifies CRUD lifecycle
- Tests run via `TF_ACC=1 go test ./...`
- CI runs acceptance tests against a staging Showbiz API
- Tests verify:
  - Resource creation, read, update, delete
  - Import of existing resources
  - Plan correctness (no diff after apply)
  - Error handling (invalid inputs, permission denied)

## Open Questions

None — all decisions resolved for initial design.
