#!/usr/bin/env bash
# Telegram notification script
# Usage: ./notify_telegram.sh <success|failure> [message]
#
# Required env vars:
#   TELEGRAM_BOT_TOKEN  — Bot token from @BotFather
#   TELEGRAM_CHAT_ID    — Chat/group ID (get from @userinfobot)

set -euo pipefail

STATUS="${1:-unknown}"
MESSAGE="${2:-}"

if [[ -z "${TELEGRAM_BOT_TOKEN:-}" || -z "${TELEGRAM_CHAT_ID:-}" ]]; then
  echo "ERROR: TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID must be set"
  exit 1
fi

if [[ "$STATUS" == "success" ]]; then
  ICON="✅"
  LABEL="SUCCESS"
elif [[ "$STATUS" == "failure" ]]; then
  ICON="❌"
  LABEL="FAILED"
else
  ICON="ℹ️"
  LABEL="INFO"
fi

TEXT="${ICON} <b>${LABEL}</b>"
if [[ -n "$MESSAGE" ]]; then
  TEXT="${TEXT}\n${MESSAGE}"
fi

curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
  -d chat_id="${TELEGRAM_CHAT_ID}" \
  -d parse_mode="HTML" \
  -d text="${TEXT}" > /dev/null

echo "Telegram notification sent: ${LABEL}"
