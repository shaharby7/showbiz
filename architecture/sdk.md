# SDK Strategy

> Status: 🟡 In Design

## Overview

Showbiz exposes two SDKs that wrap the core API:

| SDK | Language | Consumers |
|---|---|---|
| `showbiz-go` | Go | CLI, Terraform provider |
| `showbiz-ts` | TypeScript | Web UI |

## Design Principles

1. **Thin wrapper** — SDKs translate API calls into idiomatic language constructs; business logic stays in the API
2. **Generated where possible** — If we use OpenAPI, SDK clients can be partially generated
3. **Idiomatic** — Each SDK follows the conventions of its language (error handling, naming, patterns)
4. **Versioned** — SDK versions map to API versions (`/v1/`)

## Authentication

Both SDKs handle JWT authentication:
- Login with username/password → receive JWT (30-min expiry) + refresh token
- Automatic token refresh when access token expires
- Token storage is consumer-specific (CLI stores in config, UI stores in browser)

## Go SDK

### Package Structure
```
showbiz-go/
├── client.go          # Client constructor, config, auth
├── auth.go            # Login, register, token refresh
├── organizations.go   # Org CRUD, member management
├── users.go           # User operations
├── projects.go        # Project operations
├── connections.go     # Connection CRUD (provider account links)
├── resources.go       # Resource CRUD (machine, network)
├── iam.go             # IAM policies + policy attachments
├── providers.go       # List/get providers (read-only)
├── errors.go          # Error types
└── types.go           # Shared types (Resource, Connection, IAMPolicy, etc.)
```

### Usage Example
```go
client, err := showbiz.NewClient(
    showbiz.WithBaseURL("https://api.showbiz.dev"),
)

// Authenticate
token, err := client.Auth.Login(ctx, "user@example.com", "password")

// Create a connection to AWS
conn, err := client.Connections.Create(ctx, "proj_123", &showbiz.CreateConnectionInput{
    Name:     "AWS-1234",
    Provider: "aws",
    Credentials: map[string]interface{}{
        "accessKeyId":     "AKIA...",
        "secretAccessKey": "...",
    },
    Config: map[string]interface{}{
        "accountId":     "123456789012",
        "defaultRegion": "us-east-1",
    },
})

// Create a resource via that connection
resource, err := client.Resources.Create(ctx, "proj_123", &showbiz.CreateResourceInput{
    Name:         "web-server-1",
    ConnectionID: conn.ID,
    ResourceType: "machine",
    Values: map[string]interface{}{
        "size":   "medium",
        "region": "us-east",
        "image":  "ubuntu-22.04",
    },
})
```

## TypeScript SDK

### Package Structure
```
showbiz-ts/
├── src/
│   ├── client.ts
│   ├── resources/
│   │   ├── auth.ts
│   │   ├── organizations.ts
│   │   ├── users.ts
│   │   ├── projects.ts
│   │   ├── connections.ts
│   │   ├── resources.ts
│   │   ├── iam.ts
│   │   └── providers.ts
│   ├── errors.ts
│   └── types.ts
├── package.json
└── tsconfig.json
```

### Usage Example
```typescript
import { Showbiz } from '@showbiz/sdk';

const client = new Showbiz({ baseUrl: 'https://api.showbiz.dev' });

// Authenticate
await client.auth.login('user@example.com', 'password');

// Create a connection
const conn = await client.connections.create('proj_123', {
  name: 'AWS-1234',
  provider: 'aws',
  credentials: {
    accessKeyId: 'AKIA...',
    secretAccessKey: '...',
  },
  config: {
    accountId: '123456789012',
    defaultRegion: 'us-east-1',
  },
});

// Create a resource via that connection
const resource = await client.resources.create('proj_123', {
  name: 'web-server-1',
  connectionId: conn.id,
  resourceType: 'machine',
  values: {
    name: 'web-server-1',
    size: 'medium',
    region: 'us-east',
    image: 'ubuntu-22.04',
  },
});
```

## Code Generation

SDKs are **generated from an OpenAPI spec** maintained alongside the API. This ensures SDKs stay in sync with the API automatically.

- The API maintains an OpenAPI 3.x spec at `services/api/openapi.yaml`
- SDK code is generated using [oapi-codegen](https://github.com/deepmap/oapi-codegen) (Go) and [openapi-typescript-codegen](https://github.com/ferdikoomen/openapi-typescript-codegen) (TypeScript)
- Generated code is committed to the repo (not generated at build time) so consumers can inspect it
- Ergonomic wrappers (auth helpers, client constructor) are hand-written on top of the generated code

## Open Questions

None — all decisions resolved for initial design.
