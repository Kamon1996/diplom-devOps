# ---- Стадия 1: Сборка ----
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем зависимости
COPY habits-tracker/go.mod habits-tracker/go.sum ./
RUN go mod download

# Копируем весь код
COPY habits-tracker/ .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o habit-tracker ./cmd/server

# ---- Стадия 2: Финальный образ ----
FROM alpine:latest

# Устанавливаем сертификаты
RUN apk add --no-cache ca-certificates

# Создаем пользователя
RUN addgroup -g 10001 -S appgroup && \
    adduser -u 10001 -S appuser -G appgroup

# Создаем структуру папок
WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/habit-tracker /app/habit-tracker

# Копируем шаблоны из builder (ВАЖНО!)
COPY --from=builder /app/internal/templates /app/internal/templates

# Проверяем, что шаблоны скопировались (для отладки)
RUN ls -la /app/internal/templates/

# Даем права пользователю
RUN chown -R appuser:appgroup /app

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт
EXPOSE 8080

# Запускаем
CMD ["/app/habit-tracker"]
