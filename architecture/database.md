# MySQL Database Schema

> Status: 🟡 In Design — Review needed

## Overview

All persistent state is stored in MySQL. This document defines the table schemas, relationships, constraints, and indexes.

## ER Diagram

```
┌──────────────────┐       ┌──────────────────────┐
│  organizations   │       │       users           │
├──────────────────┤       ├──────────────────────┤
│ id (PK)          │◄──┐   │ email (PK)           │
│ name (UNIQUE)    │   │   │ organization_id (FK) │
│ display_name     │   ├───│ password_hash        │
│ active           │   │   │ email_verified       │
│ created_at       │   │   │ active               │
│ updated_at       │   │   │ created_at           │
└──────┬───────────┘   │   │ updated_at            │
       │               │   └───────────┬──────────┘
       │ 1:N           │               │
       │               │               │
       ├───────────────┤               │
       │               │               │
       ▼               │               │
┌──────────────────┐   │               │
│    projects      │   │               │
├──────────────────┤   │               │
│ id (PK)          │   │               │
│ name             │   │               │
│ organization_id  │───┘               │
│ description      │                   │
│ created_at       │                   │
│ updated_at       │                   │
└──┬───────────────┘                   │
   │                                   │
   │ 1:N            1:N                │
   ├────────────┐ ┌────────────┐       │
   ▼            ▼ ▼            ▼       │
┌────────────┐ ┌──────────────────┐    │
│connections │ │     resources    │    │
├────────────┤ ├──────────────────┤    │
│ id (PK)    │ │ id (PK)          │    │
│ project_id │ │ name             │    │
│ name       │ │ project_id (FK)  │    │
│ provider   │ │ connection_id(FK)│    │
│ credentials│ │ resource_type    │    │
│ config     │ │ values (JSON)    │    │
│ created_at │ │ status           │    │
│ updated_at │ │ created_at       │    │
└────────────┘ │ updated_at       │    │
               └──────────────────┘    │
                                       │
┌──────────────────────────────────────┼──────────────────┐
│              IAM                     │                   │
│                                      │                   │
│  ┌──────────────────┐                │                   │
│  │    policies       │                │                   │
│  ├──────────────────┤                │                   │
│  │ id (PK)          │                │                   │
│  │ name             │                │                   │
│  │ scope            │  (global or    │                   │
│  │ organization_id  │   org-scoped)  │                   │
│  └────────┬─────────┘                │                   │
│           │                          │                   │
│     1:N   │   1:N                    │                   │
│     ▼     │   ▼                      ▼                   │
│  ┌────────────────┐   ┌───────────────────────────┐      │
│  │policy_permis-  │   │   policy_attachments      │      │
│  │  sions         │   ├───────────────────────────┤      │
│  ├────────────────┤   │ id (PK)                   │      │
│  │ policy_id (FK) │   │ project_id (FK→projects)  │      │
│  │ permission     │   │ user_id (FK→users) ◄──────┼──────┘
│  └────────────────┘   │ policy_id (FK→policies)   │
│                       └───────────────────────────┘
└──────────────────────────────────────────────────────────┘
```

---

## Tables

### `organizations`

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | `VARCHAR(64)` | `PRIMARY KEY` | Immutable unique identifier |
| `name` | `VARCHAR(128)` | `NOT NULL, UNIQUE` | Immutable human-readable name |
| `display_name` | `VARCHAR(256)` | `NOT NULL` | Display name (mutable) |
| `active` | `BOOLEAN` | `NOT NULL DEFAULT TRUE` | Soft-delete flag |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | Creation time |
| `updated_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` | Last update time |

```sql
CREATE TABLE organizations (
    id           VARCHAR(64)  PRIMARY KEY,
    name         VARCHAR(128) NOT NULL UNIQUE,
    display_name VARCHAR(256) NOT NULL,
    active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

---

### `users`

The user's email is the primary key and unique identifier.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `email` | `VARCHAR(256)` | `PRIMARY KEY` | Immutable unique identifier (verified email) |
| `password_hash` | `VARCHAR(256)` | `NOT NULL` | bcrypt-hashed password |
| `organization_id` | `VARCHAR(64)` | `NOT NULL, FK → organizations.id` | Parent organization |
| `display_name` | `VARCHAR(256)` | `NOT NULL` | Display name (mutable) |
| `email_verified` | `BOOLEAN` | `NOT NULL DEFAULT FALSE` | Whether email has been verified |
| `active` | `BOOLEAN` | `NOT NULL DEFAULT TRUE` | Soft-delete flag |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | Creation time |
| `updated_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` | Last update time |

```sql
CREATE TABLE users (
    email           VARCHAR(256) PRIMARY KEY,
    password_hash   VARCHAR(256) NOT NULL,
    organization_id VARCHAR(64)  NOT NULL,
    display_name    VARCHAR(256) NOT NULL,
    email_verified  BOOLEAN      NOT NULL DEFAULT FALSE,
    active          BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_users_org FOREIGN KEY (organization_id) REFERENCES organizations(id)
);

CREATE INDEX idx_users_org ON users(organization_id);
```

---

### `refresh_tokens`

Stores active refresh tokens for JWT authentication.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | `VARCHAR(64)` | `PRIMARY KEY` | Token identifier |
| `user_id` | `VARCHAR(64)` | `NOT NULL, FK → users.email` | Token owner |
| `token_hash` | `VARCHAR(256)` | `NOT NULL, UNIQUE` | Hashed refresh token |
| `expires_at` | `TIMESTAMP` | `NOT NULL` | Expiration time |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | Creation time |

```sql
CREATE TABLE refresh_tokens (
    id          VARCHAR(64)  PRIMARY KEY,
    user_id     VARCHAR(64)  NOT NULL,
    token_hash  VARCHAR(256) NOT NULL UNIQUE,
    expires_at  TIMESTAMP    NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users(email) ON DELETE CASCADE
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at);
```

---

### `projects`

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | `VARCHAR(64)` | `PRIMARY KEY` | Immutable unique identifier |
| `name` | `VARCHAR(128)` | `NOT NULL` | Immutable name (unique per org) |
| `organization_id` | `VARCHAR(64)` | `NOT NULL, FK → organizations.id` | Parent organization |
| `description` | `TEXT` | | Optional description (mutable) |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | Creation time |
| `updated_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` | Last update time |

```sql
CREATE TABLE projects (
    id              VARCHAR(64)  PRIMARY KEY,
    name            VARCHAR(128) NOT NULL,
    organization_id VARCHAR(64)  NOT NULL,
    description     TEXT,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_projects_org FOREIGN KEY (organization_id) REFERENCES organizations(id),
    CONSTRAINT uq_projects_org_name UNIQUE (organization_id, name)
);

CREATE INDEX idx_projects_org ON projects(organization_id);
```

---

### `connections`

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | `VARCHAR(64)` | `PRIMARY KEY` | Immutable unique identifier |
| `name` | `VARCHAR(128)` | `NOT NULL` | Immutable name (unique per project) |
| `project_id` | `VARCHAR(64)` | `NOT NULL, FK → projects.id` | Parent project |
| `provider` | `VARCHAR(64)` | `NOT NULL` | Provider identifier (e.g., "aws", "gcp") |
| `credentials` | `BLOB` | `NOT NULL` | Encrypted provider credentials (JSON) |
| `config` | `JSON` | | Provider-specific configuration |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | Creation time |
| `updated_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` | Last update time |

```sql
CREATE TABLE connections (
    id          VARCHAR(64)  PRIMARY KEY,
    name        VARCHAR(128) NOT NULL,
    project_id  VARCHAR(64)  NOT NULL,
    provider    VARCHAR(64)  NOT NULL,
    credentials BLOB         NOT NULL,
    config      JSON,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_connections_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT uq_connections_project_name UNIQUE (project_id, name)
);

CREATE INDEX idx_connections_project ON connections(project_id);
CREATE INDEX idx_connections_provider ON connections(project_id, provider);
```

> **Note:** `credentials` is stored as an encrypted `BLOB`, not plain JSON. The application layer encrypts before writing and decrypts on read. Credentials are never returned in API GET responses.

---

### `resources`

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | `VARCHAR(512)` | `PRIMARY KEY` | Deterministic ID: `sbz:<type>:<project_id>:<conn_name>:<name>` |
| `name` | `VARCHAR(128)` | `NOT NULL` | Immutable name (unique per project) |
| `project_id` | `VARCHAR(64)` | `NOT NULL, FK → projects.id` | Parent project |
| `connection_id` | `VARCHAR(64)` | `NOT NULL, FK → connections.id` | Connection used for provisioning |
| `resource_type` | `VARCHAR(64)` | `NOT NULL` | Unified type ("machine", "network") |
| `values` | `JSON` | `NOT NULL` | Provider-agnostic resource values |
| `status` | `VARCHAR(32)` | `NOT NULL DEFAULT 'creating'` | Current state |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | Creation time |
| `updated_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` | Last update time |

```sql
CREATE TABLE resources (
    id             VARCHAR(512) PRIMARY KEY,
    name           VARCHAR(128) NOT NULL,
    project_id     VARCHAR(64)  NOT NULL,
    connection_id  VARCHAR(64)  NOT NULL,
    resource_type  VARCHAR(64)  NOT NULL,
    `values`       JSON         NOT NULL,
    status         VARCHAR(32)  NOT NULL DEFAULT 'creating',
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_resources_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_resources_connection FOREIGN KEY (connection_id) REFERENCES connections(id),
    CONSTRAINT uq_resources_project_name UNIQUE (project_id, name)
);

CREATE INDEX idx_resources_project ON resources(project_id);
CREATE INDEX idx_resources_connection ON resources(connection_id);
CREATE INDEX idx_resources_type ON resources(project_id, resource_type);
CREATE INDEX idx_resources_status ON resources(project_id, status);
```

**Status values:** `creating`, `active`, `updating`, `deleting`, `failed`

---

### `policies`

Standalone policy definitions. Can be global (platform-defined) or organization-scoped.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | `VARCHAR(64)` | `PRIMARY KEY` | Immutable unique identifier |
| `name` | `VARCHAR(128)` | `NOT NULL` | Immutable policy name (unique per scope) |
| `scope` | `VARCHAR(32)` | `NOT NULL` | `"global"` or `"organization"` |
| `organization_id` | `VARCHAR(64)` | `FK → organizations.id` | NULL for global policies, set for org policies |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | Creation time |
| `updated_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` | Last update time |

```sql
CREATE TABLE policies (
    id              VARCHAR(64)  PRIMARY KEY,
    name            VARCHAR(128) NOT NULL,
    scope           VARCHAR(32)  NOT NULL,
    organization_id VARCHAR(64),
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_policies_org FOREIGN KEY (organization_id) REFERENCES organizations(id),
    CONSTRAINT chk_policies_scope CHECK (scope IN ('global', 'organization')),
    CONSTRAINT uq_policies_global_name UNIQUE (name, scope, organization_id)
);

CREATE INDEX idx_policies_org ON policies(organization_id);
CREATE INDEX idx_policies_scope ON policies(scope);
```

> **Note:** Global policies have `organization_id = NULL` and `scope = 'global'`. They are seeded at startup and cannot be created/modified via the API.

---

### `policy_permissions`

Stores individual permissions for each policy (normalized).

| Column | Type | Constraints | Description |
|---|---|---|---|
| `policy_id` | `VARCHAR(64)` | `NOT NULL, FK → policies.id` | Parent policy |
| `permission` | `VARCHAR(64)` | `NOT NULL` | Permission string (e.g., `resource:create`, `*:*`) |

```sql
CREATE TABLE policy_permissions (
    policy_id  VARCHAR(64) NOT NULL,
    permission VARCHAR(64) NOT NULL,

    PRIMARY KEY (policy_id, permission),
    CONSTRAINT fk_policy_perms_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE CASCADE
);
```

**Valid permissions:**
- `resource:create`, `resource:read`, `resource:update`, `resource:delete`
- `connection:create`, `connection:read`, `connection:update`, `connection:delete`
- `*:*` — wildcard, grants all permissions

---

### `policy_attachments`

Binds a user to a policy in the scope of a specific project.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | `VARCHAR(64)` | `PRIMARY KEY` | Immutable unique identifier |
| `project_id` | `VARCHAR(64)` | `NOT NULL, FK → projects.id` | Project scope |
| `user_id` | `VARCHAR(256)` | `NOT NULL, FK → users.email` | User receiving permissions |
| `policy_id` | `VARCHAR(64)` | `NOT NULL, FK → policies.id` | Policy being attached |
| `created_at` | `TIMESTAMP` | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | Creation time |

```sql
CREATE TABLE policy_attachments (
    id          VARCHAR(64) PRIMARY KEY,
    project_id  VARCHAR(64) NOT NULL,
    user_id     VARCHAR(256) NOT NULL,
    policy_id   VARCHAR(64) NOT NULL,
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_attachments_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_attachments_user FOREIGN KEY (user_id) REFERENCES users(email),
    CONSTRAINT fk_attachments_policy FOREIGN KEY (policy_id) REFERENCES policies(id),
    CONSTRAINT uq_attachments_project_user_policy UNIQUE (project_id, user_id, policy_id)
);

CREATE INDEX idx_attachments_project ON policy_attachments(project_id);
CREATE INDEX idx_attachments_user ON policy_attachments(user_id);
CREATE INDEX idx_attachments_policy ON policy_attachments(policy_id);
```

> A user's effective permissions on a project = **union** of all permissions from all attached policies.

---

### Seed Data — Global Policies

```sql
-- Seeded at application startup, not via API
INSERT INTO policies (id, name, scope, organization_id) VALUES
    ('policy_global_admin', 'Administrator', 'global', NULL);

INSERT INTO policy_permissions (policy_id, permission) VALUES
    ('policy_global_admin', '*:*');
```

---

## Cascade Behavior

| Action | Cascade |
|---|---|
| Delete project | Deletes all connections, resources, and policy attachments in that project |
| Deactivate organization | Application layer deletes all projects (which cascades to connections, resources, policy attachments). Org-level policies are **not** deleted (org is deactivated, not deleted) |
| Delete connection | **Does not cascade** — must fail if resources still reference it |
| Delete user | **Does not happen** — users are deactivated only. Policy attachments remain but are ineffective for inactive users |
| Delete policy | Cascades to all policy attachments referencing it, and deletes its permissions |

---

## Design Notes

- **User identity:** The `users` table uses `email` as the primary key. There are no separate user IDs — email is the sole identifier, verified during registration.
- **Encrypted credentials:** Connection credentials are encrypted at the application layer before storage. The `BLOB` type is used instead of `JSON` to store ciphertext.
- **Resource IDs:** The `resources.id` column stores the full deterministic ID (`sbz:<type>:<project>:<conn>:<name>`). It is `VARCHAR(512)` to accommodate the composite format.
- **JSON columns:** `resources.values` and `connections.config` use MySQL's native `JSON` type for structured data with potential query support.
- **Soft-delete:** Only `organizations` and `users` have an `active` flag. All other entities are hard-deleted.
- **No `ON DELETE CASCADE` from connections to resources:** Deleting a connection with active resources should fail with an error, not silently delete resources.
