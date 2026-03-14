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
registry.Register("aws", aws.NewProvider)
registry.Register("gcp", gcp.NewProvider)
```

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
