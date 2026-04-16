# 📬 Notification Service

Сервис уведомлений на Go (учебный проект).

## Стек
- Go 1.22
- PostgreSQL 16
- Docker

## Запуск

```bash
docker-compose up -d
go run cmd/server/main.go
```

## API

```bash
# Health check
curl http://localhost:8080/health

# Создать уведомление
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user-123", "type": "email", "title": "Привет!", "body": "Текст"}'

# Получить по ID
curl http://localhost:8080/api/v1/notifications/{id}
```
