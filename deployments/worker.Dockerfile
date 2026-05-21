FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем workspace
COPY go.work go.work.sum ./

# Копируем все go.mod/go.sum для кэширования зависимостей
COPY apps/worker/go.mod apps/worker/go.sum ./apps/worker/
COPY packages/domain/go.mod packages/domain/go.sum ./packages/domain/
COPY packages/kafka/go.mod packages/kafka/go.sum ./packages/kafka/

# Скачиваем зависимости
RUN cd apps/worker && go mod download

# Копируем исходники
COPY apps/worker/ ./apps/worker/
COPY packages/ ./packages/

# Собираем
RUN cd apps/worker && go build -ldflags="-s -w" -o /app ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app .
CMD ["./app"]
