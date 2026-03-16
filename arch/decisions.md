# Architecture Decision Records

## ADR-001: Go for the Core API

**Status:** Accepted  
**Context:** Need a language for the core API backend.  
**Decision:** Go — aligns with the Go SDK, CLI, and Terraform provider. Strong concurrency support, fast compilation, good cloud ecosystem.  
**Consequences:** Team needs Go expertise. TypeScript UI will be the only non-Go component.

---

## ADR-002: SDK-First Client Architecture

**Status:** Accepted  
**Context:** CLI, Terraform provider, and UI all need to talk to the API.  
**Decision:** All consumers must use an official SDK rather than calling the API directly. This ensures consistent behavior, auth handling, and error mapping.  
**Consequences:** SDKs become critical path. API changes must flow through SDK updates.

---

## ADR-003: Provider-Agnostic Core

**Status:** Accepted  
**Context:** Platform must support multiple cloud providers.  
**Decision:** Cloud-specific logic lives only behind a provider interface. The core API, SDKs, and consumers are provider-agnostic.  
**Consequences:** Some provider-specific features may be harder to expose. Need a strategy for provider-specific config.

---

## ADR-004: JWT Authentication

**Status:** Accepted
**Context:** Need an auth mechanism for the API.
**Decision:** JWT with username/password registration. 30-minute access token expiry. Refresh tokens for renewal. No API keys for now.
**Consequences:** Need to manage JWT signing secret (static for now). Stateless auth simplifies horizontal scaling.

---

## ADR-005: API Versioning via URL Path

**Status:** Accepted
**Context:** Need a strategy for evolving the API without breaking clients.
**Decision:** Version prefix in URL path (`/v1/`). New major versions get a new prefix.
**Consequences:** Clean separation between versions. SDKs target a specific version.

---

## ADR-006: Soft-Delete for Organizations and Users

**Status:** Accepted
**Context:** Deleting orgs/users could cascade destructively and lose audit trails.
**Decision:** Organizations and users are deactivated (soft-delete), never hard-deleted. Projects and resources can be hard-deleted. **When an organization is deactivated, all of its projects and resources are deleted.**
**Consequences:** Need `active/inactive` status on org and user records. Deactivated entities should be excluded from normal queries. Org deactivation triggers cascading deletion of all projects and resources.

---

## ADR-007: Resource Schema — Connection + Type + Values

**Status:** Accepted
**Context:** Resources must be provider-agnostic at the API level.
**Decision:** Every resource has `connectionId` (which determines the provider and account), `resourceType` (unified across providers), and `values` (type-specific but provider-agnostic key-value pairs). Initial resource types: **machine** (compute instance) and **network** (virtual network/VPC).
**Consequences:** Showbiz must maintain a mapping from unified types/values to provider-specific configurations. Resources reference a connection, which resolves to a provider. Extensibility depends on good type definitions.

---

## ADR-008: Connection Entity — Project-to-Provider Link

**Status:** Accepted
**Context:** Resources need to target a specific provider account, and a project may use multiple accounts on the same provider.
**Decision:** Introduce a **Connection** entity at the project level. A connection binds a project to a specific provider account with credentials and configuration. Resources reference a `connectionId` instead of a `provider` directly. A project can have multiple connections to the same provider.
**Consequences:** Connections manage credentials (write-only, never returned in GET). The provider abstraction layer resolves connections to find the correct provider implementation and credentials at resource operation time.

---

## ADR-009: Immutable Names and IDs

**Status:** Accepted
**Context:** Entity names are used in deterministic ID generation (e.g., resource IDs: `sbz:<type>:<project>:<connection>:<name>`). Allowing renames would invalidate IDs and break references.
**Decision:** All entity names and IDs are **immutable once created**. Update operations may modify other fields but must reject changes to `id` or `name`. This applies to all entities: organizations, users, projects, connections, resources, and IAM policies.
**Consequences:** To "rename" an entity, users must delete and recreate it. This is a deliberate trade-off for ID stability and referential integrity.

---

## ADR-010: Vue.js + Vite for Web UI

**Status:** Accepted
**Context:** Need a frontend framework and build tool for the Web UI.
**Decision:** Vue.js 3 with Vite. Production builds are static (deployed to CDN, hosted-only). Development uses Vite dev server with HMR.
**Consequences:** UI is a static SPA — no SSR. All data fetching is client-side via the TypeScript SDK. CDN hosting simplifies deployment and scaling.

---

## ADR-011: No Provider-Specific Passthrough

**Status:** Accepted
**Context:** Some cloud features don't generalize across providers.
**Decision:** Showbiz is a stand-alone product. Anything that cannot be generalized across providers is excluded. There are no direct provider API calls outside the abstraction layer — users never interact with provider APIs through Showbiz.
**Consequences:** Some advanced provider-specific features will not be available. This is a deliberate trade-off for a clean, unified experience.

---

## ADR-012: Retry and Rollback at Abstraction Level

**Status:** Accepted
**Context:** Resource operations can fail due to transient errors or partial provisioning.
**Decision:** Retry logic, backoff policies, and rollback orchestration are handled by the **abstraction layer**, not by individual provider implementations. Providers return errors; the abstraction decides how to retry or roll back.
**Consequences:** Consistent retry/rollback behavior across all providers. Provider implementations stay simple — they only need to execute single operations and report success/failure.

---

## ADR-013: Drift Detection via Provider, Reconciliation via Abstraction

**Status:** Accepted
**Context:** Resources provisioned by Showbiz may drift from expected state (manual changes, external tools).
**Decision:** Providers expose a `DetectDrifts` method that compares expected state against actual provider state and returns a list of drifts. Providers do **not** reconcile — the abstraction layer decides how to handle drifts (alert, auto-fix, etc.).
**Consequences:** Clean separation of concerns. Drift detection is provider-specific (each cloud has different APIs), but reconciliation policy is centralized.

---

## ADR-014: Email as User Identity

**Status:** Accepted
**Context:** Need a unique identifier for users across the platform.
**Decision:** The user's **email address** is the unique identifier (primary key). There are no separate user IDs. Email must be **verified during registration** before the account becomes active. Email is immutable once set.
**Consequences:** All references to users (IAM attachments, org members, JWT claims) use email. Email verification is required in the registration flow. Simpler model — no mapping between opaque IDs and emails.

---

## ADR-015: Cobra for CLI Framework

**Status:** Accepted
**Context:** Need a CLI framework for the `showbiz` command-line tool.
**Decision:** Use **Cobra**. It's the de facto standard for Go CLIs (used by kubectl, docker, gh). Provides subcommand structure, flag parsing, shell completions, and help generation out of the box.
**Consequences:** Well-understood by Go developers. Shell completions for Bash/Zsh/Fish/PowerShell come for free.

---

## ADR-016: No CLI Plugin System

**Status:** Accepted
**Context:** Should the CLI support third-party plugins for extensibility?
**Decision:** No plugin system. All CLI functionality is built-in. The CLI is a thin wrapper around the Go SDK — new features are added to the SDK/API first, then exposed in the CLI.
**Consequences:** Simpler maintenance and distribution. If extensibility is needed later, it can be added without breaking existing users.

---

## ADR-017: No SDK-Level Retry

**Status:** Accepted
**Context:** Should SDKs automatically retry failed API requests?
**Decision:** No. SDKs do not implement retry or backoff logic. Failed requests return errors to the caller. Consumers (CLI, Terraform, UI) decide their own retry strategy if needed.
**Consequences:** SDKs stay thin and predictable. No hidden retry behavior that could mask issues or cause unexpected delays.

---

## ADR-018: SDK Code Generation from OpenAPI

**Status:** Accepted
**Context:** SDKs need to stay in sync with the API as it evolves.
**Decision:** Generate SDK client code from an **OpenAPI 3.x spec** maintained with the API. Go SDK uses oapi-codegen, TypeScript SDK uses openapi-typescript-codegen. Generated code is committed. Hand-written ergonomic wrappers sit on top.
**Consequences:** API changes automatically flow to SDKs via spec regeneration. The OpenAPI spec becomes a first-class artifact that must be kept up to date.

---

## ADR-019: Infrastructure and Helm Directory Structure

**Status:** Accepted  
**Context:** Need a clear separation between reusable Terraform modules, per-environment deployment configuration, and Helm charts.  
**Decision:** `infra/modules/<provider>/<module-name>` organizes modules by provider: `local/` for local-dev-only modules (minikube, mysql via Helm), `k8s/` for Kubernetes modules usable in any env (argocd, logs), `aws/` for AWS-specific modules (RDS). `infra/env/` contains Terragrunt HCL files per environment (local, staging, production). `helm/charts/` holds Showbiz-specific charts (`showbiz-app` is a generic chart reusable across services). `helm/local/` contains per-service values overrides deployed by ArgoCD in the local environment.  
**Consequences:** Terragrunt provides DRY configuration across environments. Module paths clearly indicate which provider/context they belong to. Local dev mirrors production topology via Minikube + ArgoCD.

---

## ADR-020: Terraform Module File Conventions

**Status:** Accepted  
**Context:** Terraform modules had all code (variables, outputs, resources, provider config) in a single `main.tf`, making modules hard to navigate as they grow.  
**Decision:** Every Terraform module must have separate files: `main.tf` (resources only), `variables.tf` (input variables), `outputs.tf` (output values), `versions.tf` (required providers and Terraform version), and `provider.tf` (provider configuration).  
**Consequences:** Consistent, predictable module structure across all modules. Easier code review and navigation.

---

## ADR-021: Minikube Terraform Provider

**Status:** Accepted  
**Context:** The Minikube module used `null_resource` with `local-exec` provisioners to run `minikube start/delete`. This is fragile — no state tracking, no plan/apply lifecycle, and the destroy provisioner hard-codes the profile name.  
**Decision:** Use the `scott-the-programmer/minikube` Terraform provider (`minikube_cluster` resource) which provides a proper Terraform lifecycle for Minikube clusters, including plan, apply, and destroy with full state management.  
**Consequences:** Minikube cluster is a first-class Terraform resource with proper state. Cluster attributes (certificates, host) are available as outputs for downstream modules.

---

## ADR-022: ArgoCD App-of-Apps Pattern

**Status:** Accepted  
**Context:** ArgoCD needs to know which applications to deploy in each environment. Manually creating ArgoCD Application CRs is error-prone and doesn't scale across environments.  
**Decision:** Use the app-of-apps pattern: a dedicated Helm chart (`helm/charts/app-of-apps`) generates ArgoCD `Application` CRs for each service. The `k8s/argocd` Terraform module deploys both ArgoCD and the app-of-apps chart, passing the environment name so the chart loads values from `helm/values/<environment>/<service>/values.yaml`.  
**Consequences:** Adding a new service or environment only requires adding a values file and an entry in the app-of-apps chart. The Terraform module is the single entry point for bootstrapping the entire deployment pipeline.

---

## ADR-023: FakeProvider for Local E2E Testing

**Status:** Accepted  
**Context:** The platform needs a real provider implementation to test the full resource lifecycle end-to-end (create, poll, ready, delete). Cloud providers (AWS, GCP) are not available in local development.  
**Decision:** Implement a "fakeprovider" using KubeVirt on Minikube. A dedicated microservice (`services/fakeprovider`) exposes a CRUD API for virtual machines, backed by KubeVirt VirtualMachineInstance CRDs via `client-go`. The API service implements the `Provider` interface for fakeprovider, which calls the microservice and polls asynchronously (1s intervals) until the machine is Ready. KubeVirt is deployed via a Terraform module (`infra/modules/local/kubevirt`).  
**Consequences:** Full resource lifecycle can be tested locally without cloud accounts. The fakeprovider serves as a reference implementation for future real providers. KubeVirt requires nested virtualization or a compatible driver on the host.

---

## ADR-024: Single Go Module with Shared Packages

**Status:** Accepted  
**Context:** The repository contains multiple Go services (`services/api`, `services/fakeprovider`) and Go-based components (CLI, SDK, Terraform provider) that share common patterns — HTTP helpers, Swagger UI serving, middleware, error types, and more. Maintaining separate `go.mod` files per service leads to duplicated code and inconsistent implementations (e.g., each service implementing its own Swagger handler).  
**Decision:** The entire repository uses a **single Go module** (`github.com/shaharby7/showbiz`). Shared code lives in a top-level `pkg/` directory. Services import shared packages directly. Each service still has its own `cmd/` entrypoint and `internal/` for service-specific logic.  
**Consequences:** Simpler dependency management — one `go.mod`, one `go.sum`. Shared packages are versioned together with the rest of the repo. Cross-service refactoring is easier. Trade-off: the module includes all dependencies for all services, but this is acceptable for a monorepo.

---

## ADR-025: Resource Type Interface with Schema-Driven Validation

**Status:** Accepted  
**Context:** Resources use untyped `map[string]interface{}` values with no per-type validation or schema. All resource types are treated identically in the API, SDK, and UI — the create form is a raw JSON textarea. Additionally, some resource types (e.g., network) are fully managed by Showbiz and should not require a connection to a provider.  
**Decision:** Introduce a **ResourceType interface** that each type (machine, network) must implement. The interface enforces:
- `Name()` — the type identifier
- `RequiresConnection()` — whether this type needs a provider connection
- `ValidateCreate(values)` — validate input values before creation
- `ValidateUpdate(currentValues, newValues)` — validate values before update
- `InputSchema()` / `OutputSchema()` — describe the expected input and output fields with types, required flags, and descriptions

Resource types are registered in a **ResourceTypeRegistry** at startup (parallel to the provider registry). The API exposes `GET /v1/resource-types` so consumers (UI, CLI) can discover types and render dynamic forms. `connectionId` becomes optional — only required when `RequiresConnection()` is true. For Showbiz-managed types (e.g., network), resources are stored directly with no provider call.

The resource ID format adapts: `sbz:<type>:<project>:<connection-name>:<name>` when a connection is used, `sbz:<type>:<project>:<name>` when not.

**Consequences:** Each new resource type must implement the interface and be registered. Providers continue to declare which type names they support via `ResourceTypes() []string`. The UI can render per-type tabs with type-specific input forms and output columns. Validation is enforced consistently across API, SDK, and all consumers.

---

> Add new ADRs below as decisions are made during architecture design.
