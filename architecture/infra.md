# Infrastructure

> Status: 🟡 In Design

## Overview

Showbiz infrastructure is managed as code using **Terraform** (modules), **Terragrunt** (environment orchestration), and **Helm** (Kubernetes deployments). ArgoCD provides GitOps-driven continuous deployment.

---

## Directory Structure

```
infra/                              # Infrastructure-as-code
├── modules/                        # Terraform modules (organized by provider)
│   ├── local/                      # Local development only
│   │   ├── minikube/               # Creates Minikube cluster (scott-the-programmer/minikube)
│   │   │   ├── main.tf
│   │   │   ├── variables.tf
│   │   │   ├── outputs.tf
│   │   │   └── versions.tf
│   │   └── mysql/                  # MySQL via Bitnami Helm chart
│   │       ├── main.tf
│   │       ├── variables.tf
│   │       ├── outputs.tf
│   │       └── versions.tf
│   ├── k8s/                        # Kubernetes modules (any environment)
│   │   ├── argocd/                 # Deploys ArgoCD Helm chart
│   │   │   ├── main.tf
│   │   │   ├── variables.tf
│   │   │   ├── outputs.tf
│   │   │   └── versions.tf
│   │   └── logs/                   # Prometheus + Grafana (kube-prometheus-stack)
│   │       ├── main.tf
│   │       ├── variables.tf
│   │       ├── outputs.tf
│   │       └── versions.tf
│   └── aws/                        # AWS-specific modules
│       └── mysql/                  # RDS MySQL instance
│           ├── main.tf
│           ├── variables.tf
│           ├── outputs.tf
│           └── versions.tf
└── env/                            # Terragrunt HCL files per environment
    ├── terragrunt.hcl              # Root config (provider, backend, common vars)
    ├── local/                      # Local development on Minikube
    │   ├── terragrunt.hcl          # Local environment config
    │   ├── minikube/
    │   │   └── terragrunt.hcl      # Creates Minikube cluster
    │   ├── argocd/
    │   │   └── terragrunt.hcl      # Deploys ArgoCD (depends on minikube)
    │   ├── mysql/
    │   │   └── terragrunt.hcl      # Deploys MySQL (depends on minikube)
    │   └── api/
    │       └── terragrunt.hcl      # Placeholder (ArgoCD-managed)
    ├── staging/
    │   └── terragrunt.hcl
    └── production/
        └── terragrunt.hcl

helm/                               # Helm charts and values
├── charts/                         # Charts dedicated to Showbiz
│   └── showbiz-app/                # Generic chart for deploying Showbiz services
│       ├── Chart.yaml
│       ├── values.yaml
│       └── templates/
│           ├── _helpers.tpl
│           ├── deployment.yaml
│           ├── service.yaml
│           ├── ingress.yaml
│           ├── configmap.yaml
│           └── serviceaccount.yaml
└── local/                          # Local helm values deployed by ArgoCD
    ├── api/
    │   └── values.yaml             # Values override for the API service
    └── ui/
        └── values.yaml             # Values override for the UI service
```

---

## Terraform Module Conventions

Every Terraform module follows a standard file layout (see [ADR-020](./decisions.md#adr-020-terraform-module-file-conventions)):

| File | Contents |
|---|---|
| `main.tf` | Resources only |
| `variables.tf` | Input variables |
| `outputs.tf` | Output values |
| `versions.tf` | `required_version` and `required_providers` |

---

## Module Organization

Modules live under `infra/modules/<provider>/<module-name>`, organized by the provider or context they target:

| Provider | Purpose | Modules |
|---|---|---|
| `local/` | Local development only (Minikube) | `minikube`, `mysql` |
| `k8s/` | Kubernetes (any environment) | `argocd`, `logs` |
| `aws/` | AWS cloud | `mysql` (RDS) |

### local/minikube

Creates a Minikube cluster using the [`scott-the-programmer/minikube`](https://registry.terraform.io/providers/scott-the-programmer/minikube) provider (see [ADR-021](./decisions.md#adr-021-minikube-terraform-provider)). Used only in the `local` environment.

**Outputs:** `cluster_name`, `host`, `client_certificate`, `client_key`, `cluster_ca_certificate`

### local/mysql

Deploys MySQL via the Bitnami Helm chart. Used in the `local` environment where MySQL runs inside Minikube instead of a managed cloud service.

**Outputs:** `host`, `port`, `database`

### k8s/argocd

Deploys ArgoCD via the official Argo Helm chart. Provides GitOps-driven deployments for all Showbiz services.

**Outputs:** `namespace`, `release_name`

### k8s/logs

Deploys the `kube-prometheus-stack` Helm chart (Prometheus + Grafana). Provides metrics collection, alerting, and dashboards.

**Outputs:** `namespace`, `grafana_service`, `prometheus_service`

### aws/mysql

Creates an RDS MySQL instance with encryption, automated backups, and optional multi-AZ. Used in staging and production environments.

**Outputs:** `endpoint`, `host`, `port`, `database`, `arn`

---

## Terragrunt

[Terragrunt](https://terragrunt.gruntwork.io/) provides DRY configuration across environments.

### Root Config (`infra/env/terragrunt.hcl`)

Generates provider configurations (Helm, Kubernetes) and sets up a local backend. All environment-level configs inherit from this via `find_in_parent_folders()`.

### Environments

| Environment | Path | Description |
|---|---|---|
| `local` | `infra/env/local/` | Minikube + ArgoCD + Helm-based MySQL |
| `staging` | `infra/env/staging/` | Cloud-hosted (placeholder) |
| `production` | `infra/env/production/` | Cloud-hosted (placeholder) |

### Running Terragrunt

```bash
# Apply a single module
cd infra/env/local/mysql
terragrunt apply

# Apply all modules in an environment
cd infra/env/local
terragrunt run-all apply

# Plan before applying
cd infra/env/local
terragrunt run-all plan
```

---

## Helm

### showbiz-app Chart

A generic, reusable Helm chart for deploying any Showbiz service (API, UI, future services). Parameterized via values files.

Key template features:
- **Deployment** — configurable replicas, image, ports, env vars, health probes, resource limits
- **Service** — ClusterIP by default, configurable type
- **Ingress** — optional, with host and path rules
- **ConfigMap** — arbitrary key-value config injected as env vars
- **ServiceAccount** — optional, with annotations (e.g., for IAM roles)

### Local Values

Per-service values overrides in `helm/local/` are deployed by ArgoCD in the local Minikube environment:

- `helm/local/api/values.yaml` — API service (port 8080, MySQL connection env vars)
- `helm/local/ui/values.yaml` — UI service (port 5173, API URL env var)

---

## Deployment Flow

### Local (Minikube)

```
Terragrunt apply
  → Creates Minikube cluster (local/minikube)
  → Deploys MySQL via Helm (local/mysql)
  → Deploys ArgoCD via Helm (k8s/argocd)

ArgoCD watches helm/local/
  → Deploys API using showbiz-app chart + helm/local/api/values.yaml
  → Deploys UI using showbiz-app chart + helm/local/ui/values.yaml
```

### Cloud (staging/production)

```
Terragrunt apply
  → Creates RDS MySQL (aws/mysql)
  → Deploys ArgoCD (k8s/argocd)
  → Deploys monitoring (k8s/logs)

ArgoCD watches helm/<env>/
  → Deploys services using showbiz-app chart + per-env values
```
