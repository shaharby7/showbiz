# Development Workflow

> Status: 🟡 In Design

## Overview

Showbiz is developed as a **monorepo** with all components in a single repository. CI/CD uses **GitHub Actions**. Local development uses a **devcontainer** (VS Code / GitHub Codespaces). Infrastructure is documented separately in [infra.md](./infra.md).

---

## Directory Structure

```
showbiz/
├── .devcontainer/
│   ├── devcontainer.json          # Devcontainer config
│   ├── docker-compose.yml         # MySQL + supporting services
│   └── Dockerfile                 # Dev image (Go, Node.js, Terraform, tools)
│
├── .github/
│   └── workflows/
│       ├── ci.yml                 # PR checks (lint, test, build all components)
│       ├── release-major-minor.yml # Major/minor release (all components)
│       └── release-patch.yml      # Patch release (single component)
│
├── architecture/                  # Architecture docs (this directory)
│
├── services/                      # Backend microservices
│   └── api/                       # Core API — main entry point for clients
│       ├── cmd/
│       │   └── showbiz-api/
│       │       └── main.go        # Entrypoint
│       ├── internal/
│       │   ├── auth/              # JWT, email verification
│       │   ├── handler/           # HTTP handlers (per entity)
│       │   ├── middleware/        # Auth, logging, error handling
│       │   ├── model/             # Domain models
│       │   ├── repository/        # MySQL data access
│       │   ├── service/           # Business logic
│       │   └── provider/          # Provider abstraction layer
│       ├── migrations/            # SQL migration files
│       ├── go.mod
│       └── Makefile
│   # Future microservices: services/scheduler/, services/worker/, etc.
│
├── sdk/
│   ├── go/                        # Go SDK
│   └── typescript/                # TypeScript SDK
│
├── cli/                           # CLI tool (Go, Cobra)
│
├── terraform/                     # Terraform provider (Go)
│
├── ui/                            # Web UI (Vue.js + Vite)
│
├── e2e/                           # End-to-end tests
│
├── infra/                         # Infrastructure-as-code (see infra.md)
├── helm/                          # Helm charts and values (see infra.md)
│   ├── charts/
│   │   ├── app-of-apps/           # ArgoCD app-of-apps bootstrap chart
│   │   └── showbiz-app/           # Generic chart for Showbiz services
│   └── values/                    # Per-environment values
│       └── local/                 # Local env values (api/, ui/)
│
├── examples/                      # Example projects
├── VERSION                        # Current major.minor version
└── README.md
```

---

## Versioning

### Strategy

Showbiz uses a **shared major/minor, independent patch** versioning model.

```
<major>.<minor>.<patch>
   │       │       │
   │       │       └── Per-component, backward-compatible fixes
   │       └────────── Shared across all components (new features)
   └────────────────── Shared across all components (breaking changes)
```

### Rules

1. **Major/minor versions** are shared across all components. When the project bumps from `1.2.x` to `1.3.0`, **all components** release `1.3.0` simultaneously.
2. **Patch versions** are independent per component. The Go SDK may be at `1.3.2` while the CLI is at `1.3.5`. Patch releases must maintain backward compatibility with the current major/minor version.
3. A `VERSION` file at the repo root holds the current `major.minor` (e.g., `1.3`). Each component tracks its own patch version via git tags.

### Git Tags

```
v1.3.0                    # Major/minor release (all components)
api/v1.3.1                # Patch release for API
sdk-go/v1.3.2             # Patch release for Go SDK
sdk-ts/v1.3.1             # Patch release for TypeScript SDK
cli/v1.3.5                # Patch release for CLI
terraform/v1.3.0          # Terraform provider (still on .0)
ui/v1.3.3                 # Patch release for UI
```

### Compatibility Guarantee

All components with the same **major.minor** version are guaranteed to work together. A patch release of any component is backward-compatible with all other components at the same major.minor.

---

## CI/CD

### GitHub Actions Pipelines

#### 1. CI — Pull Request Checks (`ci.yml`)

Runs on every PR. Validates all components.

```
Trigger: pull_request → main

Jobs:
  ┌─────────────────────────────────────────────┐
  │  detect-changes                              │
  │  (determines which components were modified) │
  └──────────────┬──────────────────────────────┘
                 │
    ┌────────────┼────────────┬──────────────┬──────────────┐
    ▼            ▼            ▼              ▼              ▼
  api          sdk-go      sdk-ts         cli            ui
  ┌──────┐   ┌──────┐    ┌──────┐      ┌──────┐      ┌──────┐
  │ lint │   │ lint │    │ lint │      │ lint │      │ lint │
  │ test │   │ test │    │ test │      │ test │      │ test │
  │ build│   │ build│    │ build│      │ build│      │ build│
  └──────┘   └──────┘    └──────┘      └──────┘      └──────┘
                                                         │
                                          terraform      │
                                          ┌──────┐       │
                                          │ lint │       │
                                          │ test │       │
                                          │ build│       │
                                          └──────┘
                 │
                 ▼
          ┌──────────────────────────────────┐
          │        e2e tests                 │
          │  (API + MySQL + UI, full suite)  │
          │  - API E2E via SDK               │
          │  - CLI E2E                       │
          │  - Terraform acceptance          │
          │  - UI E2E via Playwright         │
          └──────────────────────────────────┘
```

- **Change detection:** Only runs jobs for components that have changes (or their dependencies)
- **E2E tests:** Spin up full environment (API + MySQL + UI), run the `e2e/` test suite covering all components
- **All jobs must pass** before a PR can merge

#### 2. Major/Minor Release (`release-major-minor.yml`)

Triggered manually. Releases **all components** at the new version.

```
Trigger: workflow_dispatch (input: version bump type — major or minor)

Steps:
  1. Bump VERSION file (e.g., 1.2 → 1.3)
  2. Run full test suite for all components
  3. Build all artifacts:
     - API: Docker image → container registry
     - Go SDK: tag + publish Go module
     - TypeScript SDK: npm publish
     - CLI: build binaries (linux/mac/windows) → GitHub Release
     - Terraform: build + publish to Terraform Registry
     - UI: vite build → deploy to CDN
  4. Create git tag: v1.3.0
  5. Create GitHub Release with changelog
```

#### 3. Patch Release (`release-patch.yml`)

Triggered manually for a single component.

```
Trigger: workflow_dispatch (input: component name)

Steps:
  1. Determine current patch version for the component from git tags
  2. Increment patch (e.g., api/v1.3.1 → api/v1.3.2)
  3. Run tests for the component + integration tests
  4. Build and publish only that component's artifact
  5. Create git tag: <component>/v1.3.2
  6. Update GitHub Release notes
```

### Artifacts

| Component | Artifact | Published To |
|---|---|---|
| API | Docker image | Container registry (GHCR / ECR) |
| Go SDK | Go module | GitHub (Go proxy via git tags) |
| TypeScript SDK | npm package | npm registry (`@showbiz/sdk`) |
| CLI | Binaries (linux/mac/win) | GitHub Releases |
| Terraform | Provider binary | Terraform Registry |
| UI | Static assets (HTML/JS/CSS) | CDN |

---

## Testing

### Testing Pyramid

```
         ╱  E2E Tests  ╲           Few, slow, high confidence
        ╱───────────────╲
       ╱ Integration Tests╲        API + MySQL, SDK → API
      ╱─────────────────────╲
     ╱     Unit Tests        ╲     Fast, isolated, many
    ╱─────────────────────────╲
```

### Per-Component Testing

| Component | Unit Tests | Integration Tests |
|---|---|---|
| **API** | Handler, service, repository logic (mocked DB) | API + MySQL (real DB, migrations applied) |
| **Go SDK** | Client construction, request building, error parsing | SDK → running API |
| **TypeScript SDK** | Same as Go SDK | SDK → running API |
| **CLI** | Command parsing, output formatting | CLI → running API |
| **Terraform** | Resource schema validation | Terraform acceptance tests (→ running API) |
| **UI** | Component tests (Vitest + Vue Test Utils) | — |

### E2E Tests (`e2e/`)

E2E tests are an **independent component** in the `e2e/` directory. They test the full system end-to-end and can cover any combination of components:

| Suite | What it tests |
|---|---|
| `e2e/tests/api/` | Full API workflows via Go SDK or raw HTTP |
| `e2e/tests/sdk-go/` | Go SDK scenarios against a running API |
| `e2e/tests/sdk-ts/` | TypeScript SDK scenarios against a running API |
| `e2e/tests/cli/` | CLI commands against a running API |
| `e2e/tests/terraform/` | Terraform apply/plan/destroy against a running API |
| `e2e/tests/ui/` | Playwright browser tests (UI → API) |

### Test Commands

```bash
# API
cd services/api && go test ./...                    # Unit tests
cd services/api && go test -tags=integration ./...  # Integration tests (needs MySQL)

# Go SDK
cd sdk/go && go test ./...
cd sdk/go && go test -tags=integration ./...

# TypeScript SDK
cd sdk/typescript && npm test

# CLI
cd cli && go test ./...

# Terraform
cd terraform && TF_ACC=1 go test ./...    # Acceptance tests

# UI
cd ui && npm run test                      # Vitest unit/component tests

# E2E (requires running API + MySQL + UI)
cd e2e && go test ./tests/api/...          # API E2E
cd e2e && go test ./tests/cli/...          # CLI E2E
cd e2e && npm run test:ui                  # Playwright UI E2E
```

### Integration Test Environment

Integration tests require a running API + MySQL. E2E tests additionally need the UI running. In CI, these are provided by GitHub Actions service containers. Locally, the devcontainer provides everything automatically.

---

## Local Development

### Devcontainer

Local development uses a **devcontainer** that provides all required tools and services.

#### What's Included

| Tool/Service | Purpose |
|---|---|
| Go (latest) | Services, CLI, SDKs, Terraform provider |
| Node.js (LTS) | TypeScript SDK, UI |
| MySQL 8 | Database (started as a service container) |
| Terraform | Terraform provider development/testing |
| golangci-lint | Go linting |
| Air | Go hot-reload for API development |
| Vite | UI dev server with HMR |

#### Getting Started

```bash
# 1. Open in VS Code with devcontainer (or use GitHub Codespaces)
#    → Devcontainer auto-starts MySQL and applies migrations

# 2. Start the API (with hot-reload)
cd services/api && air

# 3. In another terminal, start the UI dev server
cd ui && npm run dev

# 4. The API is available at http://localhost:8080
#    The UI is available at http://localhost:5173 (Vite HMR)
```

#### Environment Variables

The devcontainer pre-configures environment variables:

```bash
SHOWBIZ_DB_HOST=localhost
SHOWBIZ_DB_PORT=3306
SHOWBIZ_DB_USER=showbiz
SHOWBIZ_DB_PASSWORD=showbiz_dev
SHOWBIZ_DB_NAME=showbiz
SHOWBIZ_JWT_SECRET=dev-secret-do-not-use-in-production
SHOWBIZ_API_URL=http://localhost:8080
```

#### Database Migrations

Migrations are applied automatically when the API starts in development mode. They can also be run manually:

```bash
cd services/api && go run ./cmd/migrate up      # Apply all pending migrations
cd services/api && go run ./cmd/migrate down    # Rollback last migration
cd services/api && go run ./cmd/migrate status  # Show migration status
```

#### Workflow: Making Changes

```
1. Create a feature branch
2. Make changes to one or more components
3. Run tests locally:
   - cd services/api && go test ./...
   - cd sdk/go && go test ./...
   - cd ui && npm test
   - cd e2e && go test ./...   (optional, full E2E)
4. Push and open a PR
5. CI validates all affected components
6. Merge to main
7. Release when ready (manual trigger)
```

---

## Branch Strategy

| Branch | Purpose |
|---|---|
| `main` | Stable, always releasable |
| `feature/*` | Feature development |
| `fix/*` | Bug fixes |
| `release/<major>.<minor>` | Release branch (created for major/minor releases) |

- All development happens on feature branches, merged to `main` via PR
- Release branches are cut from `main` for major/minor releases
- Patch releases are cut from the release branch
