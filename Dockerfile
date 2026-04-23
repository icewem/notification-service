# Этап 1 — сборка
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копируем go.mod
COPY go.mod ./
RUN go mod download

# Копируем код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Этап 2 — финальный образ
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
