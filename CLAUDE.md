# Diploma Project — DevOps Engineer

## Project Overview
Дипломный проект DevOps-инженера. Приложение — **Habits Tracker** (трекер привычек), написанное на Go.
Цель: полный DevOps-цикл — от докеризации до CI/CD, мониторинга и логирования.

## IMPORTANT CONSTRAINTS
- **Папка `habits-tracker/` — готовое приложение, НЕ ТРОГАТЬ** (кроме go.mod и go.sum)
- Вся работа — DevOps вокруг приложения: инфраструктура, CI/CD, мониторинг, логирование
- Код приложения считается готовым и финализированным

## Technology Stack (Decided)
| Компонент | Технология |
|-----------|-----------|
| Cloud | **AWS** |
| IaC | **Terraform** |
| Configuration | **Ansible** |
| CI/CD | **GitLab CI** |
| Target Deploy | **Docker Compose на EC2** |
| Monitoring | **Prometheus + Grafana** |
| Logging | **ELK Stack** (Elasticsearch + Logstash + Kibana) |
| Notifications | **Telegram Bot** |
| Reverse Proxy + TLS | **Nginx + Let's Encrypt** |
| Container Registry | **GitLab Container Registry** |

## Application (Read-Only)
- **Язык:** Go 1.25
- **БД:** PostgreSQL 15
- **Фреймворк:** стандартная библиотека Go (`net/http`, `html/template`)
- **Драйвер БД:** `pgx/v5`
- **Аутентификация:** bcrypt + cookie-сессии
- **Порт:** 8080

## Target Architecture
```
                    ┌─────────────────────────────────────────┐
                    │              AWS VPC                      │
                    │                                           │
   Internet ──────►│  ┌──────────────────────────────────┐    │
                    │  │  EC2: App Server                  │    │
                    │  │  ┌─────────────────────────────┐  │    │
                    │  │  │ Docker Compose               │  │    │
                    │  │  │  ├─ Nginx (TLS + proxy)     │  │    │
                    │  │  │  ├─ App (habits-tracker)    │  │    │
                    │  │  │  ├─ PostgreSQL 15           │  │    │
                    │  │  │  ├─ Node Exporter           │  │    │
                    │  │  │  └─ Filebeat/Promtail       │  │    │
                    │  │  └─────────────────────────────┘  │    │
                    │  └──────────────────────────────────┘    │
                    │                                           │
                    │  ┌──────────────────────────────────┐    │
                    │  │  EC2: Monitoring Server           │    │
                    │  │  ┌─────────────────────────────┐  │    │
                    │  │  │ Docker Compose               │  │    │
                    │  │  │  ├─ Prometheus              │  │    │
                    │  │  │  ├─ Grafana                 │  │    │
                    │  │  │  ├─ Elasticsearch           │  │    │
                    │  │  │  ├─ Logstash               │  │    │
                    │  │  │  ├─ Kibana                  │  │    │
                    │  │  │  └─ Nginx (TLS + proxy)    │  │    │
                    │  │  └─────────────────────────────┘  │    │
                    │  └──────────────────────────────────┘    │
                    └─────────────────────────────────────────┘

   GitLab CI ──► Build Image ──► Push to Registry ──► Deploy via SSH
                                                        │
                                                   Telegram Notify
```

## Target Project Structure
```
diplom/
├── habits-tracker/              # Application code (READ-ONLY)
│   ├── go.mod / go.sum          # Only modifiable files
│   ├── init.sql
│   ├── cmd/server/
│   ├── internal/
│   └── testutils/
│
├── Dockerfile                   # App Docker image (exists)
├── docker-compose.yml           # Local dev compose (exists)
├── .dockerignore                # (exists)
│
├── terraform/                   # Infrastructure as Code
│   ├── main.tf                  # Provider, backend config
│   ├── variables.tf             # Input variables
│   ├── outputs.tf               # Output values (IPs, DNS)
│   ├── vpc.tf                   # VPC, subnets, IGW, routes
│   ├── security_groups.tf       # Security groups
│   ├── ec2_app.tf               # App server instance
│   ├── ec2_monitoring.tf        # Monitoring server instance
│   └── terraform.tfvars.example # Example variables (no secrets)
│
├── ansible/                     # Configuration Management
│   ├── inventory/
│   │   └── hosts.yml            # Dynamic or static inventory
│   ├── playbooks/
│   │   ├── setup_common.yml     # Common setup (Docker, tools)
│   │   ├── deploy_app.yml       # Deploy application stack
│   │   ├── deploy_monitoring.yml # Deploy monitoring stack
│   │   └── setup_ssl.yml        # SSL/TLS with Let's Encrypt
│   ├── roles/
│   │   ├── docker/              # Install Docker + Compose
│   │   ├── app/                 # Deploy app compose stack
│   │   ├── monitoring/          # Deploy Prometheus + Grafana
│   │   ├── logging/             # Deploy ELK stack
│   │   └── nginx/               # Nginx reverse proxy + TLS
│   └── ansible.cfg
│
├── docker/                      # Docker Compose for deployment
│   ├── app/
│   │   ├── docker-compose.yml   # App + DB + Nginx + exporters
│   │   └── nginx/
│   │       └── nginx.conf       # Reverse proxy config
│   └── monitoring/
│       ├── docker-compose.yml   # Prometheus + Grafana + ELK
│       ├── prometheus/
│       │   └── prometheus.yml   # Scrape configs
│       ├── grafana/
│       │   └── dashboards/      # Pre-built dashboards
│       └── logstash/
│           └── pipeline.conf    # Log processing pipeline
│
├── .gitlab-ci.yml               # CI/CD Pipeline
│
├── scripts/                     # Helper scripts
│   ├── deploy.sh                # One-command full deploy
│   ├── notify_telegram.sh       # Telegram notification
│   └── init_infrastructure.sh   # Full infra bootstrap
│
├── README.md                    # Project documentation
└── CLAUDE.md                    # This file
```

## Database Schema (Reference)
- **users**: id (UUID), username, password_hash, created_at
- **habits**: id (UUID), user_id (FK), description, frequency, target_percent, created_at
- **habit_records**: id (UUID), habit_id (FK), date, done, UNIQUE(habit_id, date)
- **advice**: id (serial), message

## HTTP Routes (Reference)
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET/POST | /register | No | User registration |
| GET/POST | /login | No | User login |
| GET | /logout | No | User logout |
| GET | /habits | Yes | View habits |
| GET/POST | /habits/add | Yes | Add new habit |
| POST | /records/mark | Yes | Mark habit completion |
| GET | /report | Yes | View reports (period: week/month/year) |

## Current State
- [x] Go application with full CRUD for habits
- [x] PostgreSQL database with schema
- [x] Docker multi-stage build
- [x] docker-compose for local orchestration
- [x] Unit and integration tests
- [x] bcrypt password hashing
- [x] Non-root Docker user

## What Needs to Be Built
- [ ] Terraform: AWS VPC + EC2 instances
- [ ] Ansible: Docker installation + app deployment + monitoring setup
- [ ] Docker Compose: production stacks (app + monitoring)
- [ ] Nginx: reverse proxy + TLS/SSL (Let's Encrypt)
- [ ] GitLab CI: full pipeline (lint → test → build → deploy → notify)
- [ ] Prometheus + Grafana: metrics and dashboards
- [ ] ELK: centralized logging
- [ ] Telegram notifications
- [ ] Documentation (README)
- [ ] One-command infrastructure bootstrap script

## Key Commands
```bash
# Local development
docker-compose up --build

# Infrastructure
cd terraform && terraform init && terraform apply
cd ansible && ansible-playbook playbooks/setup_common.yml

# Full deploy from scratch
./scripts/deploy.sh

# Tests
cd habits-tracker && go test -v ./...
```

## GitLab CI Pipeline Stages
```
lint → test → build → deploy → notify
  │      │      │       │        │
  │      │      │       │        └─ Telegram: success/failure
  │      │      │       └─ SSH to EC2, docker-compose pull/up (only on main)
  │      │      └─ docker build + push to GitLab Registry
  │      └─ go test with PostgreSQL service
  └─ golangci-lint
```

## Conventions
- Terraform: HCL format, variables in terraform.tfvars
- Ansible: YAML, roles-based structure
- Docker: multi-stage builds, non-root users
- Secrets: through environment variables / Ansible Vault / GitLab CI variables
- All infrastructure reproducible from zero with one command
