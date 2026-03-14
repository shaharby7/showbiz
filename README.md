# Showbiz

A multi-cloud deployment platform that abstracts cloud providers behind a unified API. Showbiz provides a provider-agnostic interface for managing infrastructure resources across AWS, GCP, and other cloud providers.

## Architecture

See [architecture/](architecture/) for detailed design documents and ADRs.

## Repository Structure

```
showbiz/
├── services/
│   └── api/              # Main API service (Go)
├── sdk/                  # Go & TypeScript SDKs (planned)
├── cli/                  # CLI tool (planned)
├── terraform/            # Terraform provider plugin (planned)
├── ui/                   # Web UI - Vue.js (planned)
├── e2e/                  # End-to-end tests (planned)
├── infrastructure/       # IaC deployment (planned)
├── examples/             # Example projects (planned)
└── architecture/         # Architecture docs & ADRs
```

## Getting Started

### Prerequisites

Choose one of the following setups:

#### Option A: Devcontainer (Recommended)

All dependencies are pre-installed. Just open the project in VS Code with the [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension, or use GitHub Codespaces.

#### Option B: Local Setup

- **Go** 1.21+
- **MySQL** 8.0+
- **Node.js** 20+ (for UI development, optional)

### Installation

#### 1. Clone the repository

```bash
git clone <repo-url> showbiz
cd showbiz
```

#### 2. Set up the database

**If using the devcontainer**, MySQL starts automatically and env vars are pre-configured. Skip to step 3.

**If running locally**, start MySQL and create the database:

```bash
mysql -u root -p -e "CREATE DATABASE showbiz; CREATE USER 'showbiz'@'localhost' IDENTIFIED BY 'showbiz_dev'; GRANT ALL PRIVILEGES ON showbiz.* TO 'showbiz'@'localhost'; FLUSH PRIVILEGES;"
```

Then configure environment variables. Copy the example file and adjust if needed:

```bash
cp .env.example .env
source .env
export SHOWBIZ_DB_HOST SHOWBIZ_DB_PORT SHOWBIZ_DB_USER SHOWBIZ_DB_PASSWORD SHOWBIZ_DB_NAME SHOWBIZ_JWT_SECRET SHOWBIZ_API_PORT
```

The default values in `.env.example` match the devcontainer MySQL setup:

| Variable | Default | Description |
|---|---|---|
| `SHOWBIZ_DB_HOST` | `localhost` | MySQL host |
| `SHOWBIZ_DB_PORT` | `3306` | MySQL port |
| `SHOWBIZ_DB_USER` | `showbiz` | MySQL user |
| `SHOWBIZ_DB_PASSWORD` | `showbiz_dev` | MySQL password |
| `SHOWBIZ_DB_NAME` | `showbiz` | MySQL database name |
| `SHOWBIZ_JWT_SECRET` | `dev-secret-do-not-use-in-production` | JWT signing key |
| `SHOWBIZ_API_PORT` | `8080` | API server port |

#### 3. Install dependencies

```bash
# Go API
cd services/api
go mod download

# UI
cd ../../ui
npm install
```

> **Devcontainer note:** Both are installed automatically via `postCreateCommand`.

#### 4. Run database migrations

```bash
cd services/api
go run ./cmd/migrate up
```

You can check which migrations have been applied:

```bash
go run ./cmd/migrate status
```

To rollback the last migration:

```bash
go run ./cmd/migrate down
```

#### 5. Start the services

You need two terminals — one for the API and one for the UI:

**Terminal 1 — API server:**

```bash
cd services/api
go run ./cmd/showbiz-api
```

The API is now running at `http://localhost:8080`.

**Terminal 2 — UI dev server:**

```bash
cd ui
npm run dev
```

The UI is now running at `http://localhost:5173` with hot module replacement (HMR). The Vite dev server proxies API requests (`/v1/*`) to the API server at port 8080.

### Verify it works

```bash
# Health check
curl http://localhost:8080/health

# Register a user
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "secret123", "displayName": "Admin"}'

# Login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "secret123"}'
```

The login response returns an `accessToken` (JWT, 30-min expiry) and a `refreshToken`. Use the access token for authenticated requests:

```bash
export TOKEN="<accessToken from login response>"

# Create an organization
curl -X POST http://localhost:8080/v1/organizations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name": "my-org", "displayName": "My Organization"}'

# List providers
curl http://localhost:8080/v1/providers \
  -H "Authorization: Bearer $TOKEN"
```

## API Endpoints

All endpoints (except auth and health) require a valid JWT in the `Authorization: Bearer <token>` header.

### Auth
| Method | Path | Description |
|---|---|---|
| POST | `/v1/auth/register` | Register a new user |
| POST | `/v1/auth/login` | Login, returns JWT + refresh token |
| POST | `/v1/auth/refresh` | Refresh an expired access token |
| GET | `/v1/auth/me` | Get current user info |

### Organizations
| Method | Path | Description |
|---|---|---|
| POST | `/v1/organizations` | Create organization |
| GET | `/v1/organizations` | List organizations |
| GET | `/v1/organizations/{id}` | Get organization |
| PUT | `/v1/organizations/{id}` | Update organization |
| POST | `/v1/organizations/{id}/deactivate` | Deactivate (soft-delete) |
| POST | `/v1/organizations/{id}/activate` | Reactivate |
| GET | `/v1/organizations/{id}/members` | List members |
| POST | `/v1/organizations/{id}/members` | Add member |
| DELETE | `/v1/organizations/{id}/members/{email}` | Remove member |

### Projects
| Method | Path | Description |
|---|---|---|
| POST | `/v1/organizations/{orgId}/projects` | Create project |
| GET | `/v1/organizations/{orgId}/projects` | List projects |
| GET | `/v1/organizations/{orgId}/projects/{projectId}` | Get project |
| PUT | `/v1/organizations/{orgId}/projects/{projectId}` | Update project |
| DELETE | `/v1/organizations/{orgId}/projects/{projectId}` | Delete project (cascade) |

### Connections
| Method | Path | Description |
|---|---|---|
| POST | `/v1/projects/{projectId}/connections` | Create connection |
| GET | `/v1/projects/{projectId}/connections` | List connections |
| GET | `/v1/projects/{projectId}/connections/{connectionId}` | Get connection |
| PUT | `/v1/projects/{projectId}/connections/{connectionId}` | Update connection config |
| DELETE | `/v1/projects/{projectId}/connections/{connectionId}` | Delete connection |

### Resources
| Method | Path | Description |
|---|---|---|
| POST | `/v1/projects/{projectId}/resources` | Create resource |
| GET | `/v1/projects/{projectId}/resources` | List resources |
| GET | `/v1/projects/{projectId}/resources/{resourceId}` | Get resource |
| PUT | `/v1/projects/{projectId}/resources/{resourceId}` | Update resource values |
| DELETE | `/v1/projects/{projectId}/resources/{resourceId}` | Delete resource |

### Providers
| Method | Path | Description |
|---|---|---|
| GET | `/v1/providers` | List available providers |
| GET | `/v1/providers/{id}` | Get provider details |

### IAM
| Method | Path | Description |
|---|---|---|
| GET | `/v1/iam/policies` | List global policies |
| GET | `/v1/iam/policies/{policyId}` | Get policy details |
| GET | `/v1/organizations/{orgId}/policies` | List org policies |
| POST | `/v1/organizations/{orgId}/policies` | Create org policy |
| DELETE | `/v1/organizations/{orgId}/policies/{policyId}` | Delete org policy |
| GET | `/v1/organizations/{orgId}/projects/{projectId}/attachments` | List policy attachments |
| POST | `/v1/organizations/{orgId}/projects/{projectId}/attachments` | Attach policy to user |
| DELETE | `/v1/organizations/{orgId}/projects/{projectId}/attachments` | Detach policy from user |

## License

Proprietary — All rights reserved.
