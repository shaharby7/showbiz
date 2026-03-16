# Services

Showbiz is composed of several services that work together. This document describes what each one does and how they fit together.

## API Service

**Location:** `services/api/`
**Language:** Go
**Port:** 8080

The API service is the central backend of Showbiz. All operations — creating organizations, managing projects, provisioning resources, authenticating users — go through this service.

**What it does:**

- Exposes a versioned REST API (`/v1/`) for all platform operations
- Handles user authentication (JWT-based with email/password login)
- Enforces RBAC authorization on every resource operation
- Manages the full resource lifecycle: validates inputs, calls the appropriate cloud provider, tracks status, and detects drift
- Stores all state in MySQL (organizations, projects, connections, resources, users, policies)

**Key endpoints:**

| Area | Example | Description |
|---|---|---|
| Auth | `POST /v1/auth/login` | Login and receive a JWT access token |
| Organizations | `POST /v1/organizations` | Create a team workspace |
| Projects | `POST /v1/organizations/{orgId}/projects` | Create an isolated project |
| Connections | `POST /v1/projects/{projectId}/connections` | Link a project to a cloud provider |
| Resources | `POST /v1/projects/{projectId}/resources` | Provision infrastructure (VMs, networks) |
| Resource Types | `GET /v1/resource-types` | List resource types with input/output schemas |
| Providers | `GET /v1/providers` | List available cloud providers |
| IAM | `POST /v1/organizations/{orgId}/policies` | Create access control policies |

**How it connects to providers:**

The API service contains a provider registry and a resource type registry. When you create a resource, the API:
1. Looks up the resource type to validate input values against the type's schema
2. If the type requires a connection, looks up the connection to find which provider to use and calls that provider's `CreateResource` method
3. If the type is Showbiz-managed (e.g., network), creates it directly without a provider
4. Tracks provider-backed resource status asynchronously (polling the provider until ready)

---

## FakeProvider Service

**Location:** `services/fakeprovider/`
**Language:** Go
**Port:** 8081

The FakeProvider is a standalone microservice that acts as a local cloud provider for development and testing. Instead of provisioning real cloud VMs, it creates virtual machines on [KubeVirt](https://kubevirt.io/) — a Kubernetes operator that runs VMs as pods.

**What it does:**

- Provides a REST API for machine lifecycle management (create, get, list, update, delete)
- Creates KubeVirt `VirtualMachineInstance` resources on the Kubernetes cluster
- Tracks machine status asynchronously — polls KubeVirt until the VM is running and has an IP address
- Stores machine state in memory (no database required)

**Why it exists:**

Real cloud providers (AWS, GCP) cost money and require accounts. The FakeProvider lets you test the entire Showbiz resource lifecycle — from clicking "Create Machine" in the UI to having a running VM with an IP — entirely on your local machine using Minikube.

**API:**

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/v1/machines` | Create a new virtual machine |
| `GET` | `/v1/machines` | List all machines |
| `GET` | `/v1/machines/{id}` | Get a machine (includes IP when ready) |
| `PUT` | `/v1/machines/{id}` | Update machine properties |
| `DELETE` | `/v1/machines/{id}` | Delete a machine and its KubeVirt VM |
| `GET` | `/health` | Health check |

**Machine lifecycle:**

```
Initialized → Provisioning → Ready
                            ↘ Failed (timeout or error)
```

**Example — create a machine:**

```bash
curl -X POST http://localhost:8081/v1/machines \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-vm",
    "namespace": "vmis",
    "cpu": 2,
    "memoryMB": 1024,
    "image": "quay.io/kubevirt/cirros-container-disk-demo"
  }'
```

**Infrastructure requirements:**

- Kubernetes cluster (Minikube for local development)
- KubeVirt operator installed on the cluster
- A `vmis` namespace for virtual machine instances
- RBAC permissions to manage `virtualmachineinstances` resources

---

## Web UI

**Location:** `ui/`
**Language:** TypeScript (Vue.js)
**Port:** 5173 (dev server)

The Web UI is a browser-based dashboard for managing all Showbiz resources.

**What it does:**

- Provides a visual interface for all platform operations: organizations, projects, connections, resources, and IAM
- Connects to the API service via the TypeScript SDK
- Runs as a Vite dev server with hot module replacement during development

---

## How the Services Connect

```
┌────────────┐     ┌────────────┐     ┌──────────────────┐     ┌──────────┐
│   Web UI   │────►│  Showbiz   │────►│  FakeProvider    │────►│ KubeVirt │
│  (Vue.js)  │     │    API     │     │  Service         │     │  (VMs)   │
│  :5173     │     │  :8080     │     │  :8081           │     │          │
└────────────┘     └──────┬─────┘     └──────────────────┘     └──────────┘
                          │
                     ┌────▼────┐
                     │  MySQL  │
                     │  :3306  │
                     └─────────┘
```

- The **UI** talks to the **API** over HTTP (proxied through Vite in development)
- The **API** talks to the **FakeProvider** over HTTP when provisioning resources via the `fakeprovider` provider
- The **FakeProvider** talks to **KubeVirt** via the Kubernetes API to manage virtual machines
- The **API** stores all persistent state in **MySQL**
