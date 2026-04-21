package handler

import (
	"sync"

	"github.com/icewem/notification-service/internal/model"
)

// NotificationStore — потокобезопасное хранилище уведомлений
type NotificationStore struct {
	mu    sync.RWMutex
	items map[string]model.Notification
}

// NewNotificationStore — конструктор
func NewNotificationStore() *NotificationStore {
	return &NotificationStore{
		items: make(map[string]model.Notification),
	}
}

// Set — сохранить уведомление
func (s *NotificationStore) Set(n model.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[n.ID] = n
}

// Get — получить уведомление по ID
func (s *NotificationStore) Get(id string) (model.Notification, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.items[id]
	return n, ok
}

// Count — количество уведомлений
func (s *NotificationStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}
