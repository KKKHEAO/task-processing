FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем workspace
COPY go.work go.work.sum ./

# Копируем все go.mod/go.sum для кэширования зависимостей
COPY apps/outboxer/go.mod apps/outboxer/go.sum ./apps/outboxer/
COPY packages/config/go.mod packages/config/go.sum ./packages/config/
COPY packages/domain/go.mod packages/domain/go.sum ./packages/domain/
COPY packages/postgres/go.mod packages/postgres/go.sum ./packages/postgres/
COPY packages/repository/go.mod packages/repository/go.sum ./packages/repository/
COPY packages/kafka/go.mod packages/kafka/go.sum ./packages/kafka/

# Скачиваем зависимости
RUN cd apps/outboxer && go mod download

# Копируем исходники
COPY apps/outboxer/ ./apps/outboxer/
COPY packages/ ./packages/

# Собираем
RUN cd apps/outboxer && go build -ldflags="-s -w" -o /app ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app .
CMD ["./app"]
