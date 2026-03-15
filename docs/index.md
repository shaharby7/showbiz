# Showbiz

Showbiz is a multi-cloud deployment platform that lets you manage infrastructure resources — virtual machines, networks, and more — across any cloud provider through a single, unified API.

## What Problem Does It Solve?

Managing infrastructure across multiple cloud providers means learning different APIs, SDKs, and consoles for each one. Teams end up writing provider-specific glue code, managing multiple credential systems, and dealing with inconsistent resource models.

Showbiz solves this by providing:

- **One API** for all providers — create a machine on AWS, GCP, or a local KubeVirt cluster using the exact same request format
- **Connection-based provisioning** — link a project to a provider account once, then deploy resources through that connection without touching provider-specific credentials again
- **Consistent resource lifecycle** — every resource follows the same status model (`creating → active → deleting`) regardless of the underlying cloud
- **Built-in access control** — RBAC policies control who can manage resources in each project, across all providers

## How It Works

```
You (CLI / UI / Terraform)
        │
        ▼
   Showbiz API          ← single REST API for everything
        │
        ▼
  Provider Layer         ← translates to cloud-specific calls
   ┌────┼────┐
   ▼    ▼    ▼
  AWS  GCP  KubeVirt     ← actual cloud resources
```

1. **Create an organization** to group your team and projects
2. **Create a project** as an isolated workspace for resources
3. **Create a connection** linking the project to a cloud provider account (with credentials)
4. **Create resources** (machines, networks) through the connection — Showbiz handles the provider-specific translation

## Core Concepts

| Concept | Description |
|---|---|
| **Organization** | Top-level grouping for teams. Contains users, billing, and projects. |
| **Project** | An isolated workspace within an organization. Resources in different projects are completely separate. |
| **Connection** | Links a project to a cloud provider account. Holds the credentials and configuration needed to provision resources. |
| **Resource** | An infrastructure object (e.g., a virtual machine or network) deployed through a connection. |
| **Provider** | A cloud platform (AWS, GCP, Azure, or KubeVirt for local development). Providers are platform-defined — you connect to them, you don't create them. |

## Access Control

Showbiz uses role-based access control (RBAC) at the project level. Organization admins can create IAM policies and attach them to users per project, controlling who can create, read, update, or delete resources.

## Getting Started

See the [Getting Started](../README.md#getting-started) section in the main README, or jump straight to the [Services documentation](services.md) to understand each component.
