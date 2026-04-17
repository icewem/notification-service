package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/icewem/notification-service/internal/model"
)

// NotificationHandler — хендлер для работы с уведомлениями
type NotificationHandler struct {
	jobs  chan<- model.Notification
	store map[string]model.Notification
}

// NewNotificationHandler — конструктор хендлера
func NewNotificationHandler(jobs chan<- model.Notification) *NotificationHandler {
	return &NotificationHandler{
		jobs:  jobs,
		store: make(map[string]model.Notification),
	}
}

// Create — POST /api/v1/notifications
func (h *NotificationHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Берём контекст запроса — если клиент отключится,
	// ctx автоматически отменится
	ctx := r.Context()

	var req model.CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "невалидный запрос"}`, http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.Title == "" || req.Body == "" {
		http.Error(w, `{"error": "user_id, title и body обязательны"}`, http.StatusBadRequest)
		return
	}

	n := model.Notification{
		ID:        generateID(),
		UserID:    req.UserID,
		Type:      req.Type,
		Title:     req.Title,
		Body:      req.Body,
		Status:    model.StatusPending,
		CreatedAt: time.Now(),
	}

	h.store[n.ID] = n

	// Отправляем в воркер pool через select с ctx
	// если клиент отключился — не кладём задачу в канал
	select {
	case h.jobs <- n:
		// задача успешно добавлена в очередь
	case <-ctx.Done():
		// клиент отключился — не тратим ресурсы
		http.Error(w, `{"error": "запрос отменён"}`, http.StatusRequestTimeout)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(n)
}

// GetByID — GET /api/v1/notifications/{id}
func (h *NotificationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Берём контекст — пригодится когда добавим БД
	ctx := r.Context()
	_ = ctx // пока не используем явно

	id := strings.TrimPrefix(r.URL.Path, "/api/v1/notifications/")
	if id == "" {
		http.Error(w, `{"error": "id обязателен"}`, http.StatusBadRequest)
		return
	}

	n, ok := h.store[id]
	if !ok {
		http.Error(w, `{"error": "уведомление не найдено"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n)
}

// generateID — генерирует простой уникальный ID
func generateID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}
