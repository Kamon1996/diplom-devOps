# Phase 5: Monitoring (Prometheus + Grafana) — Design Document

**Date:** 2026-02-21
**Status:** Implemented

## Goal

Collect and visualize infrastructure and container metrics from the app server.

## Architecture

```
App Server                      Monitoring Server
  Node Exporter :9100  ←scrape←  Prometheus :9090
  cAdvisor :8081       ←scrape←       ↓
                                 Grafana :3000
                                   ├── Node Exporter dashboard
                                   └── Docker Containers dashboard
```

## Components

- **Prometheus v3.2.1** — metrics collection, 30d retention
- **Grafana v11.5.2** — visualization, auto-provisioned datasource + dashboards
- **Ansible role: monitoring** — deploy with app_server_ip injected via Jinja2

## Dashboards

- **Node Exporter — Infrastructure**: CPU, Memory, Disk, Network, System Load
- **Docker Containers**: Running count, CPU per container, Memory, Network I/O
