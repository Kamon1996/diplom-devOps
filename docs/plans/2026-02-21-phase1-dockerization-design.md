# Phase 1: Dockerization — Design Document

**Date:** 2026-02-21
**Status:** Implemented

## Goal

Prepare production-ready Docker configuration for deploying the Habits Tracker application on AWS EC2 with Docker Compose.

## Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Compose split | App and Monitoring separate | ELK needs 2-4GB RAM, resource isolation |
| Reverse proxy | Nginx | Industry standard, TLS-ready, security headers |
| Network isolation | frontend + backend (internal) | DB not exposed externally |
| Secrets management | .env files + .gitignore | Simple, no hardcoded values in compose |
| Metrics exporters | Node Exporter + cAdvisor | Ready for Prometheus when monitoring phase starts |
| Base image | alpine:3.21 (pinned) | Reproducible builds, no surprise breakage |

## Architecture

```
Internet → :80 → [Nginx] → frontend network → [App :8080]
                                                    ↓
                                              backend (internal)
                                                    ↓
                                              [PostgreSQL :5432]

[Node Exporter :9100] ← Prometheus (future)
[cAdvisor :8081]      ← Prometheus (future)
```

- `backend` network is `internal: true` — no external access to PostgreSQL
- App connects to both networks (receives traffic from Nginx, queries DB)
- Exporters on frontend network for external Prometheus scraping

## Files Created

| File | Purpose |
|------|---------|
| `Dockerfile` | Optimized: pinned alpine:3.21, added HEALTHCHECK |
| `docker/app/docker-compose.yml` | Production stack: Nginx + App + DB + Exporters |
| `docker/app/.env.example` | Environment template with placeholder secrets |
| `docker/app/nginx/nginx.conf` | Reverse proxy with security headers |
| `.gitignore` | Excludes .env, terraform state, SSH keys |

## Security Measures

- No hardcoded credentials anywhere
- PostgreSQL isolated on internal network, no port exposure
- Nginx adds X-Frame-Options, X-Content-Type-Options, X-XSS-Protection headers
- Hidden files blocked (location ~ /\.)
- Non-root user in app container (UID 10001)
- JSON log driver with size rotation (max 10m x 5 files)
- client_max_body_size limited to 1m
