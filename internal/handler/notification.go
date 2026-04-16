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
	// jobs — канал куда кладём задачи для воркеров
	jobs chan<- model.Notification
	// store — простое хранилище в памяти (пока без БД)
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
	// Декодируем тело запроса
	var req model.CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "невалидный запрос"}`, http.StatusBadRequest)
		return
	}

	// Валидация
	if req.UserID == "" || req.Title == "" || req.Body == "" {
		http.Error(w, `{"error": "user_id, title и body обязательны"}`, http.StatusBadRequest)
		return
	}

	// Создаём уведомление
	n := model.Notification{
		ID:        generateID(),
		UserID:    req.UserID,
		Type:      req.Type,
		Title:     req.Title,
		Body:      req.Body,
		Status:    model.StatusPending,
		CreatedAt: time.Now(),
	}

	// Сохраняем в store
	h.store[n.ID] = n

	// Отправляем в worker pool через канал
	h.jobs <- n

	// Отвечаем клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(n)
}

// GetByID — GET /api/v1/notifications/{id}
func (h *NotificationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Достаём id из URL
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/notifications/")
	if id == "" {
		http.Error(w, `{"error": "id обязателен"}`, http.StatusBadRequest)
		return
	}

	// Ищем в store
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
