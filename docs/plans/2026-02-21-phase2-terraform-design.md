# Phase 2: Terraform Infrastructure — Design Document

**Date:** 2026-02-21
**Status:** Approved

## Goal

Provision AWS infrastructure with one `terraform apply` command.

## Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Region | us-east-1 | Cheapest, user preference |
| State | S3 backend (optional, local by default) | Best practice, easy migration |
| App server | t3.small (2 vCPU, 2GB) | Enough for Go app + PostgreSQL + Nginx |
| Monitoring | t3.medium (2 vCPU, 4GB) | ELK needs 2-4GB RAM |
| AMI | Ubuntu 22.04 LTS | Standard for Docker, long-term support |
| Network | Single public subnet, 1 AZ | Sufficient for diploma scope |

## Architecture

```
VPC 10.0.0.0/16
  └── Public Subnet 10.0.1.0/24 (us-east-1a)
        ├── EC2 app-server (t3.small) ← EIP
        └── EC2 monitoring-server (t3.medium) ← EIP
```

## Security Groups

- Exporters (9100, 8081) open only from monitoring-sg
- Logstash (5044) open only from app-sg
- SSH restricted to user's IP
- Prometheus/Kibana restricted to user's IP
- HTTP/HTTPS/Grafana open to all
