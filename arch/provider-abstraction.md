# Provider Abstraction Layer

> Status: 🟡 In Design

## Overview

The provider abstraction layer is the core of Showbiz's cloud-agnostic design. It defines a common interface that all cloud providers implement, allowing the API to manage **resources** (machines, networks) without knowing provider-specific details.

Providers are **platform-defined and read-only** from the API consumer's perspective. They are registered at startup and listed via `GET /v1/providers`. Users interact with providers through **Connections** — a connection binds a project to a specific provider account with credentials. When a resource is created, the API resolves the connection to find the provider implementation and credentials.

## Provider Interface

```go
// Provider defines the contract every cloud provider must implement.
type Provider interface {
    // Name returns the provider identifier (e.g., "aws", "gcp", "azure").
    Name() string

    // ResourceTypes returns the resource types this provider supports.
    ResourceTypes() []string

    // ValidateCredentials checks that the configured credentials work.
    ValidateCredentials(ctx context.Context) error

    // CreateResource provisions a resource.
    CreateResource(ctx context.Context, input *CreateResourceInput) (*ResourceOutput, error)

    // GetResource returns the current state of a resource.
    GetResource(ctx context.Context, resourceID string) (*ResourceOutput, error)

    // UpdateResource modifies an existing resource.
    UpdateResource(ctx context.Context, input *UpdateResourceInput) (*ResourceOutput, error)

    // DeleteResource tears down a resource.
    DeleteResource(ctx context.Context, resourceID string) error

    // DetectDrifts compares expected state against actual provider state
    // and returns a list of drifted resources. The abstraction layer
    // handles reconciliation decisions — providers only report drifts.
    DetectDrifts(ctx context.Context, resources []ResourceExpectedState) ([]DriftReport, error)
}

// CreateResourceInput is the input for creating a resource.
type CreateResourceInput struct {
    ConnectionID string                 // Resolves to provider + credentials
    ResourceType string                 // "machine", "network"
    Values       map[string]interface{} // Provider-agnostic values
}

// ResourceOutput is the output after a resource operation.
type ResourceOutput struct {
    ID           string
    ResourceType string
    Status       string                 // "creating", "active", "updating", "deleting", "failed"
    Values       map[string]interface{} // Current values
}

// ResourceExpectedState is the expected state of a resource for drift detection.
type ResourceExpectedState struct {
    ID           string
    ResourceType string
    Values       map[string]interface{}
}

// DriftReport describes a single resource that has drifted from expected state.
type DriftReport struct {
    ResourceID     string
    ExpectedValues map[string]interface{}
    ActualValues   map[string]interface{}
    Missing        bool // true if the resource no longer exists at the provider
}
```

## Resource Type Mapping

Each provider maps unified resource types to its own implementation:

| Unified Type | What it represents | Example provider mapping |
|---|---|---|
| `machine` | A compute instance (VM/server) | AWS → EC2, GCP → Compute Engine, Azure → VM |
| `network` | A virtual network / VPC | AWS → VPC, GCP → VPC Network, Azure → VNet |

The **values** for each resource type are provider-agnostic. The provider implementation is responsible for translating them to provider-specific API calls.

## Provider Registration

Providers are compiled-in and registered at startup:

```go
registry := providers.NewRegistry()
registry.Register("stub", provider.NewStubProvider())
registry.Register("fakeprovider", provider.NewFakeProvider(cfg.FakeProviderURL))
// Future: registry.Register("aws", aws.NewProvider())
```

---

## Implemented Providers

### stub

A mock provider for unit testing and development. All operations succeed immediately with hardcoded responses. Supports `machine` and `network` resource types. No real infrastructure is provisioned.

### fakeprovider

A local development provider backed by [KubeVirt](https://kubevirt.io/) — a Kubernetes operator for running virtual machines. Used for end-to-end testing of the full resource lifecycle without requiring cloud accounts. See [ADR-023](./decisions.md#adr-023-fakeprovider-for-local-e2e-testing).

**Architecture:**

```
┌──────────────┐      HTTP       ┌────────────────────┐     client-go     ┌──────────────┐
│  Showbiz API │ ──────────────► │  FakeProvider Svc  │ ────────────────► │   KubeVirt   │
│  (provider)  │                 │  (services/        │                   │   (VMIs on   │
│              │ ◄────────────── │   fakeprovider)    │ ◄──────────────── │   Minikube)  │
└──────────────┘   JSON response └────────────────────┘    watch/poll     └──────────────┘
```

**Resource types:** `machine`

**How it works:**

1. The API service's `FakeProvider` implements the `Provider` interface
2. On `CreateResource`, it calls the fakeprovider microservice via HTTP `POST /v1/machines`
3. The fakeprovider service creates a KubeVirt `VirtualMachineInstance` CR via `client-go`
4. Creation is **asynchronous** — the API immediately receives status `Initialized`
5. The API's resource service spawns a background poller (1-second intervals) that calls `GetResource` on the provider
6. Once the VM reaches `Running` phase and has an IP, status transitions to `Ready` → the API updates the resource to `active`
7. On `DeleteResource`, the provider calls `DELETE /v1/machines/{id}` which removes the VMI

**FakeProvider connection schema:**

```json
{
  "name": "local-kubevirt",
  "provider": "fakeprovider",
  "credentials": {},
  "config": {}
}
```

No credentials are required — the fakeprovider service runs in the same cluster and accesses KubeVirt directly.

**Machine resource values:**

```json
{
  "name": "my-vm",
  "connectionId": "conn_id",
  "resourceType": "machine",
  "values": {
    "cpu": 2,
    "memoryMB": 1024,
    "image": "quay.io/kubevirt/cirros-container-disk-demo",
    "namespace": "vmis"
  }
}
```

**Resource lifecycle:**

| Status | Meaning |
|---|---|
| `Initialized` | Create request sent to fakeprovider, VMI not yet created |
| `Provisioning` | VMI created on KubeVirt, waiting for IP assignment |
| `active` | VM is running, IP available in resource values |
| `failed` | VMI failed to start or timed out (5 min) |

**FakeProvider microservice** (`services/fakeprovider/`):

| Endpoint | Description |
|---|---|
| `POST /v1/machines` | Create machine (async — returns Initialized) |
| `GET /v1/machines` | List all machines |
| `GET /v1/machines/{id}` | Get machine (includes IP when Ready) |
| `PUT /v1/machines/{id}` | Update machine properties |
| `DELETE /v1/machines/{id}` | Delete machine and VMI |

The microservice uses an in-memory store (no database) and `k8s.io/client-go` dynamic client for KubeVirt interactions.

**Infrastructure requirements:**

- KubeVirt operator deployed on the cluster (`infra/modules/local/kubevirt`)
- A `vmis` namespace for VM instances
- The fakeprovider service must have RBAC permissions to create/get/delete VMIs

## Provider Lifecycle

```
Register → Configure (credentials) → Validate → CreateResource / UpdateResource / DeleteResource
```

## Abstraction Boundaries

Showbiz is a **stand-alone product** — anything that cannot be generalized across providers does not belong in the platform. Provider-specific features that don't map to the unified model are excluded.

**What the abstraction layer handles (not the provider):**
- Retry logic and backoff policies for failed operations
- Rollback orchestration (e.g., if a create partially fails)
- Drift reconciliation decisions (based on drift reports from providers)

**What provider implementations handle:**
- Translating unified resource types/values to native cloud API calls
- Credential authentication with their specific cloud
- Reporting drifts between expected and actual state

**What is explicitly excluded from Showbiz:**
- Provider-specific features that cannot be generalized
- Direct provider API passthrough — users never interact with provider APIs through Showbiz

## Open Questions

None — all decisions resolved for initial design.
