package handler

import (
	"sync"

	"github.com/icewem/notification-service/internal/model"
)

// Store — интерфейс хранилища уведомлений
// позволяет подменять реализацию в тестах
type Store interface {
	Set(n model.Notification)
	Get(id string) (model.Notification, bool)
	Count() int
}

// NotificationStore — потокобезопасное хранилище
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

func (s *NotificationStore) Set(n model.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[n.ID] = n
}

func (s *NotificationStore) Get(id string) (model.Notification, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.items[id]
	return n, ok
}

func (s *NotificationStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}
