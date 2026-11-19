
# Этап 1: сборка приложения
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Сначала копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Этап 2: запуск приложения
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8000

CMD ["./server"]
