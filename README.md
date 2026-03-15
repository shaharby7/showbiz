# Showbiz

A multi-cloud deployment platform that abstracts cloud providers behind a unified API. Showbiz provides a provider-agnostic interface for managing infrastructure resources across AWS, GCP, and other cloud providers.

## Repository Structure

```
showbiz/
├── services/
│   ├── api/              # Core API service (Go)
│   └── fakeprovider/     # Local KubeVirt provider (Go)
├── pkg/                  # Shared Go packages
├── sdk/                  # Go & TypeScript SDKs
├── cli/                  # CLI tool
├── ui/                   # Web UI (Vue.js)
├── infra/                # Terraform/Terragrunt infrastructure
├── helm/                 # Helm charts and per-environment values
├── arch/                 # Architecture docs & ADRs
├── docs/                 # User-facing documentation
└── go.mod                # Single Go module
```

## Documentation

| Document | Description |
|---|---|
| [docs/index.md](docs/index.md) | Project overview and core concepts |
| [docs/services.md](docs/services.md) | Service descriptions and how they connect |
| [docs/local-development.md](docs/local-development.md) | Local development setup guide |
| [arch/](arch/) | Architecture design documents and ADRs |

## Quick Start

See the [Local Development Guide](docs/local-development.md) for full setup instructions.

## API Documentation

Both backend services expose interactive Swagger UI:

| Service | Swagger UI |
|---|---|
| API | [http://localhost:8080/swagger/](http://localhost:8080/swagger/) |
| FakeProvider | [http://localhost:8081/swagger/](http://localhost:8081/swagger/) |

## License

Proprietary — All rights reserved.
