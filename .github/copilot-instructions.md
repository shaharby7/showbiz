# Copilot Instructions for Showbiz

## Architecture First

All implementation work must be guided by the architecture documents in `arch/`. Before writing any code:

1. Read the relevant architecture file(s) to understand the intended design
2. Implement strictly according to the architecture — do not introduce patterns, components, or behaviors not described in the architecture docs
3. If the architecture docs don't cover a scenario, flag it and propose an update to the architecture before implementing

Key architecture files:
- `arch/overview.md` — system layers, domain model, and key principles
- `arch/api.md` — API service design, endpoints, and conventions
- `arch/provider-abstraction.md` — provider interface, resource lifecycle, and implemented providers
- `arch/ui.md` — Web UI design
- `arch/sdk.md` — SDK design
- `arch/cli.md` — CLI design
- `arch/infra.md` — infrastructure and deployment
- `arch/decisions.md` — architecture decision records (ADRs)

## Swagger UI for All Backend Services

Every backend service (API, FakeProvider, and any future services) must expose a Swagger UI for interactive API documentation. When implementing or modifying a backend service:

1. Maintain an OpenAPI/Swagger spec that accurately reflects all endpoints
2. Serve Swagger UI at `/swagger/` so developers can explore and test the API from a browser
3. Keep the spec in sync with the code — when endpoints are added, changed, or removed, update the spec in the same commit

## Documentation After Implementation

After implementing any feature or change, update the user-facing documentation in `docs/` to reflect what was actually built:

1. Update `docs/index.md` if the change affects core concepts, the project description, or how users interact with the platform
2. Update `docs/services.md` if the change affects a service's API, behavior, configuration, or how services connect to each other
3. Create new files in `docs/` if the change introduces a new user-facing topic not covered by existing docs
4. Documentation must describe the implemented behavior, not aspirational design — keep `arch/` for design intent and `docs/` for what actually works today
