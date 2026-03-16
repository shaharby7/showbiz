# Web UI Design

> Status: 🟡 In Design

## Overview

The Showbiz Web UI is a browser-based dashboard for managing organizations, projects, connections, resources, and access policies. Built with TypeScript and Vue.js, consuming the TypeScript SDK.

## Tech Stack

| Layer | Technology |
|---|---|
| Framework | Vue.js 3 + Vite |
| Language | TypeScript |
| Styling | Tailwind CSS |
| State management | Pinia (+ TanStack Query for server state) |
| Component library | PrimeVue |

## Hosting & Development

| Mode | Setup |
|---|---|
| **Production** | Static build (`vite build`), deployed to CDN (hosted-only) |
| **Development** | Vite dev server with hot module replacement (HMR) |

- The UI is a **static SPA** — no server-side rendering. The production build outputs static assets (HTML, JS, CSS) served by a CDN.
- All data fetching happens client-side via the TypeScript SDK against the Showbiz API.

## Authentication

- Login page with email/password
- JWT stored in browser (httpOnly cookie or secure storage)
- Automatic token refresh via the TypeScript SDK
- Redirect to login on 401

## Navigation

- **Org switcher** — persistent dropdown in the top navigation bar for users who belong to multiple organizations. Switching orgs updates all views to that org's context.
- **Project switcher** — secondary dropdown (visible after selecting an org) to navigate between projects within the active organization. Switching projects updates resource, connection, and IAM views.

## Key Pages

- **Login / Register** — Authentication flow (email verification required)
- **Dashboard** — Overview of organizations and projects
- **Organization Detail** — Members, billing, projects list
- **Project Detail** — Connections, resources list, IAM policies, project settings
- **Connection Management** — Create/edit/delete connections to provider accounts
- **Resources** — Tabbed view with a **separate tab per resource type** (Machines, Networks, ...). Each tab renders type-specific columns based on the resource type's input/output schema (e.g., Machines shows CPU, memory, IP; Networks shows CIDR, gateway)
- **Resource Create** — Dynamic form that adapts to the selected resource type's input schema. Connection dropdown is shown only for resource types that require a connection.
- **IAM Management** — Browse global policies, manage org policies, attach/detach policies to users per project
- **Providers** — Browse available provider types and their capabilities (read-only)
- **User Settings** — Profile, password change

## Design Principles

- Responsive design
- Dark mode support
- Accessible (WCAG 2.1 AA)

## Open Questions

None — all decisions resolved for initial design.
