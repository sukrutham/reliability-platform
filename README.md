# Reliability Platform

A small backend service used as the foundation for exploring Site Reliability
Engineering and Platform Engineering practices: containerization, Kubernetes
deployment, CI/CD, observability, SLOs, and incident response.

## Status

🚧 Work in progress. This project is being built incrementally.

- [x] Phase 1: Core REST API + unit tests
- [ ] Phase 2: Containerization (Docker)
- [ ] Phase 3: CI pipeline (GitHub Actions)
- [ ] Phase 4: Kubernetes deployment
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
