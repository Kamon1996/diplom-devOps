# Phase 4: GitLab CI/CD — Design Document

**Date:** 2026-02-21
**Status:** Implemented

## Goal

Automated CI/CD pipeline: lint, test, build, deploy, notify on every commit.

## Pipeline

```
lint → test → build → deploy (master only) → notify (always)
```

| Stage | Image | Purpose | Branches |
|-------|-------|---------|----------|
| lint | golangci-lint:v1.64 | Static analysis | all |
| test | golang:1.26 + postgres | Unit/integration tests | all |
| build | docker:27 + dind | Build + push to GitLab Registry | all |
| deploy | alpine + ssh | SSH to EC2, pull + up | master |
| notify | alpine + curl | Telegram success/failure | all |

## Deploy Flow

1. CI builds image → pushes to `registry.gitlab.com/user/project:sha`
2. SSH to app-server
3. Update APP_IMAGE in .env
4. docker login to GitLab Registry
5. docker compose pull && up
6. Health check (12 retries x 5s = 60s timeout)

## GitLab CI Variables Required

| Variable | Type | Description |
|----------|------|-------------|
| SSH_PRIVATE_KEY | File | SSH key for EC2 access |
| APP_SERVER_IP | Variable | App server Elastic IP |
| TELEGRAM_BOT_TOKEN | Variable, masked | Telegram bot token |
| TELEGRAM_CHAT_ID | Variable | Telegram chat ID |

## Notifications

- Success: green message with pipeline link
- Failure: red message with pipeline link
- Both include: project, branch, commit hash and title
