# Core API Design

> Status: 🟡 In Design

## Overview

The Showbiz API is the central backend service. All platform operations flow through this API. Communication is strictly **JSON over HTTP**.

## Tech Stack

| Component | Technology |
|---|---|
| Language | Go |
| Database | MySQL |
| API format | REST + JSON (only) |
| Auth | JWT (email/password registration with email verification + token refresh) |

## Domain Model

```
Provider (cloud abstraction — read-only, platform-defined)

Organization
 ├── Users (members, authenticated via email/password)
 ├── Billing
 └── Project
      ├── IAM Policies (RBAC — per-project, per-user permissions)
      ├── Connections (link a project to a specific provider account/credentials)
      └── Resources (provider-agnostic abstractions, deployed via a Connection)
```

### Key Relationships

- **Organizations** own **Users**, **Projects**, and billing
- **Projects** are scoped to an organization; resources are **fully isolated** between projects
- **Users** belong to an organization; their permissions are controlled by **IAM policies** at the project level
- **Connections** are scoped to a project; they define how the project authenticates with a specific provider account. A project can have **multiple connections to the same provider** (e.g., two different AWS accounts)
- **Resources** are created within a project **via a Connection** (not a provider directly), subject to IAM permissions
- **Providers** are platform-level abstractions — not owned by any organization

---

## Entities

### Provider

An abstraction of a cloud provider (e.g., AWS, GCP, Azure). Providers are **platform-defined and read-only** from the API consumer's perspective.

| Operation | Endpoint | Description |
|---|---|---|
| List providers | `GET /v1/providers` | Returns all available providers |
| Get provider | `GET /v1/providers/{id}` | Returns implementation details for a provider |

Providers cannot be created or modified via the API — they are registered at the platform level and extended by adding new provider implementations.

---

### Organization

A top-level tenant that groups users, projects, and billing.

| Operation | Endpoint | Description |
|---|---|---|
| Create | `POST /v1/organizations` | Create a new organization |
| Get | `GET /v1/organizations/{id}` | Get organization details |
| Update | `PUT /v1/organizations/{id}` | Update organization |
| Deactivate | `POST /v1/organizations/{id}/deactivate` | Deactivate organization |
| Activate | `POST /v1/organizations/{id}/activate` | Reactivate organization |
| List members | `GET /v1/organizations/{id}/members` | List organization members |
| Add member | `POST /v1/organizations/{id}/members` | Add a user to the organization |
| Remove member | `DELETE /v1/organizations/{id}/members/{userId}` | Remove a user |

---

### User

A user account created under an organization. The user's **email is the unique identifier** — there are no separate user IDs. Email must be **verified during registration** before the account becomes active.

| Operation | Endpoint | Description |
|---|---|---|
| Register | `POST /v1/auth/register` | Register with email/password (triggers email verification) |
| Verify email | `POST /v1/auth/verify` | Verify email with a verification token |
| Login | `POST /v1/auth/login` | Authenticate with email/password, receive a JWT |
| Refresh | `POST /v1/auth/refresh` | Refresh an expired JWT |
| Get current user | `GET /v1/auth/me` | Get the authenticated user's profile |
| Get user | `GET /v1/users/{email}` | Get user details |
| Update user | `PUT /v1/users/{email}` | Update user profile |
| Deactivate user | `POST /v1/users/{email}/deactivate` | Deactivate user account |
| Activate user | `POST /v1/users/{email}/activate` | Reactivate a deactivated user |

---

### Project

A project belongs to an organization. Resources are **completely isolated** between projects — there is no cross-project resource visibility or access.

| Operation | Endpoint | Description |
|---|---|---|
| Create | `POST /v1/organizations/{orgId}/projects` | Create a project |
| List | `GET /v1/organizations/{orgId}/projects` | List projects in an org |
| Get | `GET /v1/projects/{id}` | Get project details |
| Update | `PUT /v1/projects/{id}` | Update project |
| Delete | `DELETE /v1/projects/{id}` | Delete project and all its resources |

---

### Connection

A connection defines how a project connects to a specific cloud provider account. It holds the credentials, account identifiers, and configuration needed to provision resources on that provider. A project can have **multiple connections to the same provider** (e.g., different AWS accounts).

| Operation | Endpoint | Description |
|---|---|---|
| Create | `POST /v1/projects/{projectId}/connections` | Create a connection |
| List | `GET /v1/projects/{projectId}/connections` | List connections in a project |
| Get | `GET /v1/projects/{projectId}/connections/{id}` | Get connection details |
| Update | `PUT /v1/projects/{projectId}/connections/{id}` | Update a connection |
| Delete | `DELETE /v1/projects/{projectId}/connections/{id}` | Delete a connection |

#### Connection Schema

```json
{
  "id": "conn_101",
  "projectId": "proj_456",
  "name": "AWS-1234",
  "provider": "aws",
  "credentials": {
    "accessKeyId": "AKIA...",
    "secretAccessKey": "..."
  },
  "config": {
    "accountId": "123456789012",
    "defaultRegion": "us-east-1"
  }
}
```

- **`name`** — a human-readable identifier for this connection (e.g., `"AWS-1234"`, `"GCP-prod"`)
- **`provider`** — which cloud provider this connection targets (must match a registered provider)
- **`credentials`** — provider-specific authentication credentials
- **`config`** — provider-specific account configuration (account ID, default region, etc.)

> **Note:** Credentials are write-only — they are accepted on create/update but never returned in GET responses.

---

### Resource

A resource is an abstraction for **any object created by the organization within a project**. Resources are managed through the API and subject to IAM permissions.

| Operation | Endpoint | Description |
|---|---|---|
| Create | `POST /v1/projects/{projectId}/resources` | Create a resource |
| List | `GET /v1/projects/{projectId}/resources` | List resources in a project |
| Get | `GET /v1/projects/{projectId}/resources/{id}` | Get resource details |
| Update | `PUT /v1/projects/{projectId}/resources/{id}` | Update a resource |
| Delete | `DELETE /v1/projects/{projectId}/resources/{id}` | Delete a resource |

All resource operations are gated by the user's IAM policy for the given project.

#### Resource Types

Resource types are platform-defined abstractions that describe the schema and behavior of each kind of resource. Each type implements a common interface that provides validation and schema metadata. Not all resource types require a connection to a provider.

| Operation | Endpoint | Description |
|---|---|---|
| List types | `GET /v1/resource-types` | List all registered resource types with their schemas |
| Get type | `GET /v1/resource-types/{name}` | Get a specific resource type's schema |

##### Resource Type Schema

```json
{
  "name": "machine",
  "requiresConnection": true,
  "inputSchema": [
    { "name": "cpu", "type": "number", "required": true, "description": "Number of CPU cores" },
    { "name": "memoryMB", "type": "number", "required": true, "description": "Memory in megabytes" },
    { "name": "image", "type": "string", "required": true, "description": "OS image identifier" },
    { "name": "namespace", "type": "string", "required": false, "description": "Target namespace" }
  ],
  "outputSchema": [
    { "name": "ip", "type": "string", "description": "Assigned IP address" },
    { "name": "providerResourceId", "type": "string", "description": "Provider-side resource ID" }
  ]
}
```

##### Registered Resource Types

| Type | Requires Connection | Description |
|---|---|---|
| `machine` | Yes | A compute instance (VM/server). Provisioned via a provider connection. |
| `network` | No | A virtual network managed entirely by Showbiz. No provider required. |

#### Resource Schema

Every resource has a **name** (unique within the project), an optional **connection** (required only for provider-backed types), a **resource type**, and **values** (validated against the type's input schema):

```json
{
  "id": "sbz:machine:proj_456:AWS-1234:web-server-1",
  "name": "web-server-1",
  "projectId": "proj_456",
  "connectionId": "conn_101",
  "resourceType": "machine",
  "values": {
    "cpu": 2,
    "memoryMB": 1024,
    "image": "ubuntu-22.04"
  }
}
```

A Showbiz-managed resource (no connection):

```json
{
  "id": "sbz:network:proj_456:my-network",
  "name": "my-network",
  "projectId": "proj_456",
  "connectionId": null,
  "resourceType": "network",
  "values": {
    "cidr": "10.0.0.0/16"
  }
}
```

#### Resource ID Format

Resource IDs are deterministic and follow two structures depending on whether a connection is used:

```
sbz:<resource-type>:<project-id>:<connection-name>:<resource-name>   (provider-backed)
sbz:<resource-type>:<project-id>:<resource-name>                     (Showbiz-managed)
```

#### Resource Fields

- **`name`** — user-provided name, **unique within the project**
- **`connectionId`** — the connection used to provision this resource. **Required** for provider-backed types (e.g., `machine`), **null** for Showbiz-managed types (e.g., `network`)
- **`resourceType`** — must match a registered resource type
- **`values`** — key-value pairs **validated against the resource type's input schema**

---

### IAM (Identity & Access Management)

The IAM model is **RBAC (Role-Based Access Control)**. Policies are **standalone resources** that define a set of permissions. They exist at two levels:

- **Global policies** — platform-defined by Showbiz, inherited by all organizations and projects. Immutable via API.
- **Organization policies** — defined at the organization level, available to all projects within that organization.

Policies are then **attached** to users in the scope of a specific project.

#### Global Policies (read-only)

| Operation | Endpoint | Description |
|---|---|---|
| List | `GET /v1/iam/policies` | List all global policies |
| Get | `GET /v1/iam/policies/{id}` | Get a global policy |

**Built-in global policies:**

| Policy | Permissions | Description |
|---|---|---|
| `Administrator` | `*:*` (all permissions) | Full access to everything in the project |

#### Organization Policies

| Operation | Endpoint | Description |
|---|---|---|
| Create | `POST /v1/organizations/{orgId}/iam/policies` | Create an org-level policy |
| List | `GET /v1/organizations/{orgId}/iam/policies` | List org policies |
| Get | `GET /v1/organizations/{orgId}/iam/policies/{id}` | Get an org policy |
| Update | `PUT /v1/organizations/{orgId}/iam/policies/{id}` | Update policy permissions |
| Delete | `DELETE /v1/organizations/{orgId}/iam/policies/{id}` | Delete a policy |

##### Policy Schema

```json
{
  "id": "policy_101",
  "name": "DevOps",
  "scope": "organization",
  "organizationId": "org_123",
  "permissions": [
    "resource:create", "resource:read", "resource:update", "resource:delete",
    "connection:create", "connection:read"
  ]
}
```

#### Policy Attachments

Attaching a policy to a user grants them the policy's permissions **on a specific project**.

| Operation | Endpoint | Description |
|---|---|---|
| Attach | `POST /v1/projects/{projectId}/iam/attachments` | Attach a policy to a user on this project |
| List | `GET /v1/projects/{projectId}/iam/attachments` | List all attachments in a project |
| Detach | `DELETE /v1/projects/{projectId}/iam/attachments/{id}` | Remove a policy attachment |

##### Attachment Schema

```json
{
  "id": "att_201",
  "projectId": "proj_456",
  "userId": "user@example.com",
  "policyId": "policy_101"
}
```

A user can have **multiple policies** attached on the same project. Their effective permissions are the **union** of all attached policies.

#### Permission Model

Permissions follow the format `{entity}:{action}`. Available permissions:
- `resource:create`, `resource:read`, `resource:update`, `resource:delete`
- `connection:create`, `connection:read`, `connection:update`, `connection:delete`
- `*:*` — wildcard, grants all permissions (used by the `Administrator` global policy)

---

## API Conventions

### Immutability

All entity **names** and **IDs** are immutable once created. They cannot be changed via update operations. This applies to all entities: organizations, users, projects, connections, resources, and IAM policies. Update endpoints may modify other fields (e.g., values, config, permissions) but must reject attempts to change `id` or `name`.

### Request/Response Format

All requests and responses use `Content-Type: application/json`. No other content types are supported.

### Error Envelope

```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "User does not have 'resource:create' permission on project proj_456"
  }
}
```

### Pagination

List endpoints use cursor-based pagination:

```json
{
  "data": [...],
  "pagination": {
    "next_cursor": "abc123",
    "has_more": true
  }
}
```

### Standard HTTP Status Codes

| Code | Usage |
|---|---|
| `200` | Success |
| `201` | Created |
| `400` | Bad request / validation error |
| `401` | Not authenticated |
| `403` | Forbidden (IAM denial) |
| `404` | Not found |
| `409` | Conflict (e.g., duplicate) |
| `500` | Internal server error |

---

## Decided

- ✅ **JWT** for authentication — 30-minute access token expiry, refresh tokens for renewal, static signing key
- ✅ **API versioning** via URL path prefix (`/v1/`)
- ✅ **No rate limiting** for now
- ✅ **Soft-delete** (deactivate) for organizations and users; hard-delete for projects and resources
- ✅ **Immutable names and IDs** — all entity names and IDs are immutable once created
- ✅ **Resource schema** — typed with optional `connectionId`, `resourceType`, and `values` validated against type schema
- ✅ **Initial resource types** — `machine` (provider-backed) and `network` (Showbiz-managed)
- ✅ **Resource Type interface** — each type implements validation and schema methods; registered in a ResourceTypeRegistry
- ✅ **Optional connection** — resource types declare whether they require a connection; Showbiz-managed types (e.g., network) have `connectionId: null`
- ✅ **Async resource creation** — creating a resource calls the provider, which may be asynchronous. The resource is returned immediately with its initial status (e.g., `Initialized`). The API polls the provider in the background (1-second intervals) until the resource reaches `active` or `failed` status.

## Registered Providers

| Provider | Resource Types | Backend | Environment |
|---|---|---|---|
| `fakeprovider` | `machine` | KubeVirt VMIs via `services/fakeprovider` | Local (Minikube) |

See [provider-abstraction.md](./provider-abstraction.md) for full provider documentation.

## Open Questions

None — all decisions resolved for initial design.
