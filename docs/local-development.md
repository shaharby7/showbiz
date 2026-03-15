# Local Development

This guide walks through setting up a full local Showbiz environment — FakeProvider on Minikube, API + MySQL + UI via devcontainer, all connected end-to-end.

## Prerequisites

- **Docker** — required for both Minikube and the devcontainer
- **Minikube** — local Kubernetes cluster
- **Terragrunt** and **Terraform** — for deploying infrastructure
- **kubectl** and **Helm** — for Kubernetes management
- **VS Code** with the [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension (recommended)

## Overview

The local setup has two parts:

1. **Minikube** runs the FakeProvider + KubeVirt (manages virtual machines)
2. **Devcontainer** runs the API, MySQL, and UI (via Docker Compose)

The API in the devcontainer connects to the FakeProvider on Minikube through a port-forward on the host.

```
┌─ Devcontainer (Docker Compose) ──────────────────┐
│  UI (:5173) → API (:8080) → MySQL (:3306)        │
│                  │                                │
│         SHOWBIZ_FAKEPROVIDER_URL                  │
│          = host.docker.internal:8081              │
└──────────────────┼────────────────────────────────┘
                   │ (port-forward)
┌─ Minikube ───────┼────────────────────────────────┐
│  FakeProvider (:8081) → KubeVirt → VMs            │
└───────────────────────────────────────────────────┘
```

## Step 1: Deploy FakeProvider Infrastructure

Use the Terragrunt configuration in `infra/env/local/` to create the Minikube cluster, install KubeVirt, and deploy the FakeProvider.

```bash
cd infra/env/local
terragrunt run-all apply
```

This creates:
- A Minikube cluster (`showbiz` profile) with 4 CPUs, 8GB RAM
- KubeVirt operator (v1.2.0) for running virtual machines
- A `vmis` namespace for VM instances
- ArgoCD for GitOps (optional, deploys services via Helm)

### Expose the FakeProvider Locally

Port-forward the FakeProvider service so the API container can reach it:

```bash
kubectl --context showbiz -n showbiz port-forward svc/fakeprovider-showbiz-app --address 0.0.0.0 8081:80
```

Keep this running in a terminal. Verify it works:

```bash
curl http://localhost:8081/health
# {"status":"ok"}
```

## Step 2: Start the Devcontainer

Open the project in VS Code and reopen in the devcontainer. When prompted, choose **Showbiz API**.

This starts three containers via Docker Compose:

| Service | URL | Description |
|---|---|---|
| API | `http://localhost:8080` | Go API server |
| UI | `http://localhost:5173` | Vue.js with hot reload |
| MySQL | `localhost:13306` | Database |

The API container is preconfigured with `SHOWBIZ_FAKEPROVIDER_URL=http://host.docker.internal:8081`, which routes through the host to the Minikube port-forward.

## Step 3: Run Database Migrations

Inside the devcontainer terminal:

```bash
cd services/api
go run ./cmd/migrate up
```

Then start the API server:

```bash
go run ./cmd/showbiz-api
```

## Step 4: Verify End-to-End

The quickest way to verify everything is connected is to create a machine through the UI:

1. Open `http://localhost:5173` in your browser
2. Register a user and log in
3. Create an organization and a project
4. Create a connection with provider `fakeprovider` (no credentials needed)
5. Create a resource of type `machine` with values like:
   ```json
   {
     "cpu": 1,
     "memoryMB": 128,
     "image": "quay.io/kubevirt/cirros-container-disk-demo",
     "namespace": "vmis"
   }
   ```
6. Watch the resource status go from `Initialized` → `Provisioning` → `active`

Once the resource reaches `active`, the full pipeline is working: UI → API → FakeProvider → KubeVirt → running VM.

### Verify via CLI

You can also verify with curl:

```bash
# Health checks
curl http://localhost:8080/health    # API
curl http://localhost:8081/health    # FakeProvider

# Register and login
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "dev@example.com", "password": "secret123", "displayName": "Dev"}'

curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "dev@example.com", "password": "secret123"}'
# → save the accessToken

# Create org, project, connection, then resource...
```

## API Documentation

Both backend services expose interactive Swagger UI:

- **API**: [http://localhost:8080/swagger/](http://localhost:8080/swagger/)
- **FakeProvider**: [http://localhost:8081/swagger/](http://localhost:8081/swagger/)

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `SHOWBIZ_DB_HOST` | `localhost` | MySQL host |
| `SHOWBIZ_DB_PORT` | `13306` | MySQL port |
| `SHOWBIZ_DB_USER` | `showbiz` | MySQL user |
| `SHOWBIZ_DB_PASSWORD` | `showbiz_dev` | MySQL password |
| `SHOWBIZ_DB_NAME` | `showbiz` | MySQL database name |
| `SHOWBIZ_JWT_SECRET` | `dev-secret-...` | JWT signing key |
| `SHOWBIZ_API_PORT` | `8080` | API server port |
| `SHOWBIZ_FAKEPROVIDER_URL` | `http://localhost:8081` | FakeProvider URL (set to `http://host.docker.internal:8081` in devcontainer) |

## Useful Commands

```bash
# Check KubeVirt VMs
kubectl --context showbiz -n vmis get virtualmachineinstances

# Console into a VM (install virtctl first)
virtctl console <vm-name> -n vmis
# Default cirros credentials: cirros / gocubsgo

# Restart the FakeProvider after a code change
cd /path/to/showbiz
docker build -t showbiz-fakeprovider:local -f services/fakeprovider/Dockerfile .
# (build inside minikube's docker: eval $(minikube -p showbiz docker-env))
kubectl --context showbiz -n showbiz rollout restart deployment fakeprovider-showbiz-app

# Run database migrations
cd services/api && go run ./cmd/migrate up
cd services/api && go run ./cmd/migrate status
cd services/api && go run ./cmd/migrate down
```
