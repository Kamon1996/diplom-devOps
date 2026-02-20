# ---- Стадия 1: Сборка ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем зависимости отдельно для кэширования слоя
COPY habits-tracker/go.mod habits-tracker/go.sum ./
RUN go mod download

# Копируем весь код
COPY habits-tracker/ .

# Собираем статический бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o habit-tracker ./cmd/server

# ---- Стадия 2: Финальный образ ----
FROM alpine:3.21

# Устанавливаем сертификаты
RUN apk add --no-cache ca-certificates

# Создаем непривилегированного пользователя
RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

WORKDIR /app

# Копируем бинарник и шаблоны из builder
COPY --from=builder /app/habit-tracker /app/habit-tracker
COPY --from=builder /app/internal/templates /app/internal/templates

# Права и переключение на non-root
RUN chown -R appuser:appgroup /app
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO /dev/null http://localhost:8080/login || exit 1

CMD ["/app/habit-tracker"]
