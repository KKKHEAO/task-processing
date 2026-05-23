FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем все go.mod/go.sum для кэширования зависимостей
COPY apps/worker/go.mod apps/worker/go.sum ./apps/worker/
COPY packages/config/go.mod packages/config/go.sum ./packages/config/
COPY packages/domain/go.mod packages/domain/go.sum ./packages/domain/
COPY packages/kafka/go.mod packages/kafka/go.sum ./packages/kafka/
COPY packages/logger/go.mod packages/logger/go.sum ./packages/logger/
COPY packages/postgres/go.mod packages/postgres/go.sum ./packages/postgres/
COPY packages/repository/go.mod packages/repository/go.sum ./packages/repository/

# Скачиваем зависимости
RUN cd apps/worker && go mod download

# Копируем исходники
COPY apps/worker/ ./apps/worker/
COPY packages/ ./packages/

# Собираем
RUN cd apps/worker && GOFLAGS=-mod=mod go build -ldflags="-s -w" -o /app/main ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app .
CMD ["/app/main"]
