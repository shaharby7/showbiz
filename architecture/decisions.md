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

> Add new ADRs below as decisions are made during architecture design.
