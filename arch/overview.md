# System Overview

## Vision

Showbiz enables developers to manage cloud resources across any provider through a single, unified API. The platform abstracts away provider-specific details, letting teams focus on their applications rather than cloud infrastructure.

## Architecture Layers

```
┌─────────────────────────────────────────────────────┐
│                   Consumer Layer                     │
│                                                      │
│   ┌─────────┐  ┌────────────────┐  ┌─────────────┐  │
│   │   CLI   │  │ Terraform Plug │  │   Web UI    │  │
│   │  (Go)   │  │     (Go)       │  │ (TypeScript) │  │
│   └────┬────┘  └───────┬────────┘  └──────┬──────┘  │
│        │               │                  │          │
├────────┼───────────────┼──────────────────┼──────────┤
│        │          SDK Layer               │          │
│        │                                  │          │
│   ┌────▼────────────────┐  ┌──────────────▼───────┐  │
│   │      Go SDK         │  │   TypeScript SDK     │  │
│   └─────────┬───────────┘  └──────────┬───────────┘  │
│             │                         │              │
├─────────────┼─────────────────────────┼──────────────┤
│             │       API Layer         │              │
│             │                         │              │
│             └────────┐   ┌────────────┘              │
│                      ▼   ▼                           │
│              ┌───────────────────┐                    │
│              │   Showbiz API     │                    │
│              │   (Go + MySQL)    │                    │
│              └────────┬──────────┘                    │
│                       │                              │
├───────────────────────┼──────────────────────────────┤
│                       │  Provider Abstraction Layer   │
│                       ▼                              │
│   ┌──────────────────────────────────────────────┐   │
│   │          Provider Interface                   │   │
│   │                                               │   │
│   │  ┌──────────┐ ┌──────────┐ ┌──────────┐      │   │
│   │  │Provider A│ │Provider B│ │Provider N│ ...   │   │
│   │  └──────────┘ └──────────┘ └──────────┘      │   │
│   └──────────────────────────────────────────────┘   │
└──────────────────────────────────────────────────────┘
```

## Domain Model

```
Provider (read-only, platform-defined)

ResourceType (platform-defined: machine, network, ...)
 ├── Defines input/output schema and validation
 └── Declares whether a provider connection is required

Organization
 ├── Users (JWT auth, email/password, verified email)
 ├── Billing
 └── Project (fully isolated)
      ├── IAM Policies (RBAC per user)
      ├── Connections (link to a provider account + credentials)
      └── Resources (typed: machine, network, ...)
           ├── Provider-backed (deployed via a Connection)
           └── Showbiz-managed (no connection required)
```

## Components

### 1. Core API (Go + MySQL)
The central backend service. Exposes a REST + JSON API (versioned under `/v1/`). Owns business logic, JWT authentication, RBAC authorization, resource orchestration, and state management. See [api.md](./api.md).

### 2. Go SDK
Client library for Go consumers. Used by:
- **CLI tool** — local developer workflows
- **Terraform provider** — infrastructure-as-code integration

### 3. TypeScript SDK
Client library for TypeScript/JavaScript consumers. Used by:
- **Web UI** — browser-based management dashboard

### 4. CLI Tool (Go)
Command-line interface for developers. Wraps the Go SDK to provide commands for managing organizations, projects, connections, resources, and IAM policies. See [cli.md](./cli.md).

### 5. Terraform Provider Plugin (Go)
A Terraform provider that exposes Showbiz resources (projects, connections, machines, networks, IAM policies) as Terraform HCL, using the Go SDK under the hood. See [terraform.md](./terraform.md).

### 6. Web UI (TypeScript)
A browser-based dashboard for managing organizations, projects, connections, resources, and access policies. Built on the TypeScript SDK. See [ui.md](./ui.md).

### 7. Provider Abstraction Layer
A pluggable interface within the API that normalizes resource operations (create/update/delete machines, networks) across cloud providers. Each provider implements a common interface. See [provider-abstraction.md](./provider-abstraction.md).

### 8. FakeProvider Service (Go)
A standalone microservice (`services/fakeprovider`) that manages virtual machines on KubeVirt. Used as the local development provider — implements the same resource lifecycle as a real cloud provider (async creation, status polling, IP assignment) but runs entirely on Minikube. The API service integrates with it via the `fakeprovider` provider implementation. See [provider-abstraction.md](./provider-abstraction.md#fakeprovider).

## Key Principles

- **Provider-agnostic core** — Cloud-specific logic lives only behind the provider interface
- **SDK-first** — All consumers go through SDKs; no direct API calls from CLI/UI
- **Single API** — One source of truth for all operations
- **Pluggable providers** — Adding a new cloud = implementing one interface
- **Connection-based provisioning** — Resources are deployed via connections, not directly to providers
- **RBAC everywhere** — All resource operations gated by IAM policies
