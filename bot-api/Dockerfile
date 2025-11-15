FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение (правильный путь)
RUN go build -o bot-api ./cmd/bot-api/main.go

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder stage
COPY --from=builder /app/bot-api .

# Экспортируем порт (если нужно для health checks)
EXPOSE 8080

# Запускаем приложение
CMD ["./bot-api"]