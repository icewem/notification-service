package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/icewem/notification-service/internal/apperror"
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
func (h *NotificationHandler) Create(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var req model.CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Оборачиваем ошибку с контекстом
		return apperror.BadRequest("невалидный JSON")
	}

	// Валидация
	if req.UserID == "" {
		return apperror.BadRequest("user_id обязателен")
	}
	if req.Title == "" {
		return apperror.BadRequest("title обязателен")
	}
	if req.Body == "" {
		return apperror.BadRequest("body обязателен")
	}
	if req.Type != model.TypeEmail && req.Type != model.TypePush && req.Type != model.TypeSMS {
		return apperror.BadRequest("type должен быть email, push или sms")
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

	select {
	case h.jobs <- n:
	case <-ctx.Done():
		return apperror.New(http.StatusRequestTimeout, "запрос отменён", ctx.Err())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(n)
	return nil
}

// GetByID — GET /api/v1/notifications/{id}
func (h *NotificationHandler) GetByID(w http.ResponseWriter, r *http.Request) error {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/notifications/")
	if id == "" {
		return apperror.BadRequest("id обязателен")
	}

	n, ok := h.store[id]
	if !ok {
		// Используем sentinel error + wrapping
		return apperror.NotFound(fmt.Sprintf("уведомление %s не найдено", id))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n)
	return nil
}

// generateID — генерирует простой уникальный ID
func generateID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}
