# Diploma Project — Phases & Sequence

## Overview

| # | Фаза | Статус | Зависимости |
|---|------|--------|-------------|
| 1 | Докеризация | DONE | — |
| 2 | Инфраструктура (Terraform + AWS) | DONE | — |
| 3 | Конфигурация (Ansible) | DONE | Phase 2 |
| 4 | CI/CD (GitLab CI) | DONE | Phase 1, 3 |
| 5 | Мониторинг (Prometheus + Grafana) | TODO | Phase 3 |
| 6 | Логирование (ELK) | TODO | Phase 3 |
| 7 | Скрипт "одной команды" | TODO | Phase 2-6 |
| 8 | Документация и защита | TODO | Phase 1-7 |

---

## Phase 1: Докеризация — DONE

**Коммит:** `25058b6`

- [x] Dockerfile: pin alpine:3.21, HEALTHCHECK
- [x] Production docker-compose: Nginx + App + PostgreSQL + Node Exporter + cAdvisor
- [x] Nginx reverse proxy с security headers
- [x] .env.example (без хардкоженных секретов)
- [x] .gitignore

---

## Phase 2: Инфраструктура (Terraform + AWS)

**Цель:** Поднять AWS-инфраструктуру одной командой `terraform apply`.

### Ресурсы:
- VPC + public subnet + Internet Gateway + Route Table
- Security Groups (SSH, HTTP/HTTPS, exporters, ELK)
- EC2 `app-server` (t3.small) — приложение
- EC2 `monitoring-server` (t3.medium) — мониторинг + логи
- Elastic IP для стабильных адресов
- SSH Key Pair
- S3 backend для terraform state

### Файлы:
```
terraform/
├── main.tf               # Provider, backend
├── variables.tf          # Input variables
├── outputs.tf            # IPs, DNS, SSH commands
├── vpc.tf                # VPC, subnet, IGW, routes
├── security_groups.tf    # SG rules
├── ec2.tf                # EC2 instances + EIP
└── terraform.tfvars.example
```

---

## Phase 3: Конфигурация (Ansible)

**Цель:** Настроить серверы после создания Terraform: Docker, деплой стеков.

### Роли:
- `docker` — установка Docker + Compose на все серверы
- `app` — деплой app docker-compose стека
- `monitoring` — деплой Prometheus + Grafana
- `logging` — деплой ELK стека
- `nginx-ssl` — TLS/SSL через Let's Encrypt (когда будет домен)

### Файлы:
```
ansible/
├── ansible.cfg
├── inventory/
│   └── hosts.yml         # Из terraform output
├── playbooks/
│   ├── setup_common.yml
│   ├── deploy_app.yml
│   ├── deploy_monitoring.yml
│   └── deploy_logging.yml
└── roles/
    ├── docker/
    ├── app/
    ├── monitoring/
    └── logging/
```

---

## Phase 4: CI/CD (GitLab CI)

**Цель:** Автоматическая сборка, тестирование, деплой при каждом коммите.

### Pipeline:
```
lint → test → build → deploy → notify
```

- **lint** — golangci-lint (все ветки)
- **test** — go test с PostgreSQL service (все ветки)
- **build** — docker build + push в GitLab Registry (все ветки)
- **deploy** — SSH + docker-compose pull/up (только main/master)
- **notify** — Telegram бот (всегда, success/failure)

### Файлы:
```
.gitlab-ci.yml
scripts/
├── notify_telegram.sh
└── deploy.sh
```

---

## Phase 5: Мониторинг (Prometheus + Grafana)

**Цель:** Визуализация метрик инфраструктуры и приложения.

### Компоненты:
- Prometheus — сбор метрик с Node Exporter, cAdvisor, PostgreSQL Exporter
- Grafana — дашборды (Infrastructure, Docker, PostgreSQL, App Health)
- Alerting rules (опционально)

### Файлы:
```
docker/monitoring/
├── docker-compose.yml
├── prometheus/
│   └── prometheus.yml
└── grafana/
    └── provisioning/
        ├── datasources/datasource.yml
        └── dashboards/dashboard.yml
```

---

## Phase 6: Логирование (ELK)

**Цель:** Централизованный сбор и поиск по логам.

### Компоненты:
- Filebeat (на app-сервере) → собирает Docker-логи
- Logstash (на monitoring-сервере) → обрабатывает
- Elasticsearch → хранит
- Kibana → визуализирует

### Файлы:
```
docker/monitoring/          # Добавляем к Phase 5
├── logstash/
│   └── pipeline/logstash.conf
docker/app/                 # Добавляем к Phase 1
├── filebeat/
│   └── filebeat.yml
```

---

## Phase 7: Скрипт "одной команды"

**Цель:** `./scripts/deploy.sh` — вся инфраструктура с нуля.

### Алгоритм:
1. `terraform init && terraform apply -auto-approve`
2. Генерация Ansible inventory из `terraform output`
3. `ansible-playbook setup_common.yml`
4. `ansible-playbook deploy_app.yml`
5. `ansible-playbook deploy_monitoring.yml`
6. `ansible-playbook deploy_logging.yml`
7. Health checks
8. Telegram-уведомление

---

## Phase 8: Документация и защита

- README.md с архитектурой, инструкциями, скриншотами
- Презентация (3-5 мин)
- Демонстрация CI/CD (10-12 мин)
