FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем workspace
COPY go.work go.work.sum ./

# Копируем все go.mod/go.sum для кэширования зависимостей
COPY apps/api/go.mod apps/api/go.sum ./apps/api/
COPY packages/config/go.mod packages/config/go.sum ./packages/config/
COPY packages/domain/go.mod packages/domain/go.sum ./packages/domain/
COPY packages/postgres/go.mod packages/postgres/go.sum ./packages/postgres/
COPY packages/repository/go.mod packages/repository/go.sum ./packages/repository/
COPY proto/go.mod proto/go.sum ./proto/

# Скачиваем зависимости
RUN cd apps/api && go mod download

# Копируем исходники
COPY apps/api/ ./apps/api/
COPY packages/ ./packages/
COPY proto/ ./proto/

# Собираем
RUN cd apps/api && go build -ldflags="-s -w" -o /app ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app .
CMD ["./app"]
