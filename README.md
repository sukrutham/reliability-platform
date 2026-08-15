# Reliability Platform

A Go backend service built to demonstrate core Site Reliability and Platform
Engineering practices end-to-end: containerization, Kubernetes deployment,
CI/CD, observability, SLOs, and incident response.

## Why this project

Reliability engineering is best shown, not just described. This repo is a
working system built incrementally — each phase adds a real capability
(container image, CI pipeline, k8s deployment, metrics, alerting, runbooks)
rather than a one-off demo.

| Phase | Demonstrates |
|---|---|
| API + tests | Go, REST API design, automated testing |
| Containerization | Docker, multi-stage builds, distroless images |
| CI pipeline | GitHub Actions, automated build/test/publish |
| Kubernetes | Deployments, health probes, config management |
| Observability | Prometheus/Grafana, metrics-driven operations |
| SLOs & alerting | Reliability engineering, error budgets |
| Chaos testing | Failure injection, recovery measurement |
| IaC | Terraform/Ansible, reproducible infrastructure |

## Status

🚧 Work in progress. This project is being built incrementally.

- [x] Phase 1: Core REST API + unit tests
- [x] Phase 2: Containerization (Docker)
- [x] Phase 3: CI pipeline (GitHub Actions)
- [x] Phase 4: Kubernetes deployment
- [ ] Phase 5: Observability (Prometheus + Grafana)
- [ ] Phase 6: SLOs, alerting, runbooks
- [ ] Phase 7: Chaos / resilience testing
- [ ] Phase 8: Infrastructure as Code

## Phase 1: API

A simple Go REST API with:

- `GET /healthz` — liveness endpoint
- `GET /readyz` — readiness endpoint
- `GET /api/v1/todos` — list todos
- `POST /api/v1/todos` — create a todo
- `GET /api/v1/todos/{id}` — get a todo by id

### Run locally

```bash
go run ./cmd/api
```

### Run tests

```bash
go test ./... -v
```

## Phase 2: Containerization

Multi-stage `Dockerfile` builds a static binary and runs it in a
`distroless` base image for a small, low-attack-surface container.

### Build and run with Docker

```bash
docker build -t reliability-platform:local .
docker run -p 8080:8080 reliability-platform:local
```

### Run with Docker Compose

```bash
docker compose up --build
```

## Phase 3: CI Pipeline

GitHub Actions runs on every push/PR to `main`:

1. **Lint** — `gofmt` and `go vet`
2. **Test** — `go test ./... -race`
3. **Build & push** — multi-stage Docker image published to GitHub Container
   Registry (`ghcr.io`), tagged with `latest` and the commit SHA (main only)

See [`.github/workflows/ci.yml`](.github/workflows/ci.yml).

## Phase 4: Kubernetes Deployment

Manifests in [`deploy/k8s`](deploy/k8s) run the app as a 2-replica
Deployment with liveness/readiness probes wired to `/healthz` and
`/readyz`, config via ConfigMap, and a ClusterIP Service.

The API also handles `SIGTERM` for graceful shutdown — in-flight requests
are drained before the process exits, so Kubernetes rollouts and pod
evictions don't drop traffic. Pods run with a hardened `securityContext`:
non-root UID, read-only root filesystem, no privilege escalation, all
Linux capabilities dropped.

### Run locally with kind

```bash
kind create cluster --name reliability-platform
docker build -t ghcr.io/sukrutham/reliability-platform:latest .
kind load docker-image ghcr.io/sukrutham/reliability-platform:latest --name reliability-platform

kubectl apply -f deploy/k8s/namespace.yaml
kubectl apply -f deploy/k8s/configmap.yaml
kubectl apply -f deploy/k8s/deployment.yaml
kubectl apply -f deploy/k8s/service.yaml

kubectl -n reliability-platform rollout status deployment/reliability-platform
kubectl -n reliability-platform port-forward svc/reliability-platform 8080:80
```
