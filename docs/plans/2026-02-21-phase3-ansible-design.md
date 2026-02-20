# Phase 3: Ansible Configuration — Design Document

**Date:** 2026-02-21
**Status:** Implemented

## Goal

Configure EC2 servers after Terraform provisioning: install Docker, deploy the application stack.

## Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Secrets | Ansible Vault | Encrypted in repo, best practice |
| Docker install | Official Docker apt repo | Stable, up-to-date versions |
| Deploy method | Copy files + docker compose up | Matches Phase 1 production compose |
| Inventory | Static hosts.yml | Simple, IPs from terraform output |

## Roles

**docker** — Install Docker CE + Compose plugin on all servers:
- apt dependencies → Docker GPG key → Docker repo → install → enable service

**app** — Deploy application stack on app server:
- Create /opt/habits-tracker/ → copy compose + nginx + init.sql → generate .env from vault → docker compose up → health check

## Playbooks

- `setup_common.yml` — runs docker role on all hosts
- `deploy_app.yml` — runs app role on app_servers group

## Usage

```bash
cd ansible
# 1. Fill inventory IPs from terraform output
# 2. Create and encrypt vault
cp inventory/group_vars/vault.yml.example inventory/group_vars/vault.yml
ansible-vault encrypt inventory/group_vars/vault.yml
# 3. Run
ansible-playbook playbooks/setup_common.yml --ask-vault-pass
ansible-playbook playbooks/deploy_app.yml --ask-vault-pass
```
