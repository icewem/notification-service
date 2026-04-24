package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/icewem/notification-service/internal/model"
)

// MockStore — мок хранилища для тестов
type MockStore struct {
	items map[string]model.Notification
}

func NewMockStore() *MockStore {
	return &MockStore{items: make(map[string]model.Notification)}
}

func (m *MockStore) Set(n model.Notification) {
	m.items[n.ID] = n
}

func (m *MockStore) Get(id string) (model.Notification, bool) {
	n, ok := m.items[id]
	return n, ok
}

func (m *MockStore) Count() int {
	return len(m.items)
}

// TestCreate — тесты на создание уведомления
func TestCreate(t *testing.T) {
	tests := []struct {
		name         string
		body         map[string]string
		expectedCode int
		expectedErr  string
	}{
		{
			name: "успешное создание email",
			body: map[string]string{
				"user_id": "user-123",
				"type":    "email",
				"title":   "Тест",
				"body":    "Текст",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name: "успешное создание push",
			body: map[string]string{
				"user_id": "user-456",
				"type":    "push",
				"title":   "Пуш",
				"body":    "Текст",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "пустой user_id",
			body:         map[string]string{"type": "email", "title": "Тест", "body": "Текст"},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "user_id обязателен",
		},
		{
			name:         "пустой title",
			body:         map[string]string{"user_id": "user-123", "type": "email", "body": "Текст"},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "title обязателен",
		},
		{
			name:         "неверный тип",
			body:         map[string]string{"user_id": "user-123", "type": "telegram", "title": "Тест", "body": "Текст"},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "type должен быть email, push или sms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаём буферизированный канал и мок store
			jobs := make(chan model.Notification, 10)
			store := NewMockStore()
			h := NewNotificationHandlerWithStore(jobs, store)

			// Сериализуем тело запроса
			bodyBytes, _ := json.Marshal(tt.body)

			// Создаём тестовый запрос
			req := httptest.NewRequest(http.MethodPost, "/api/v1/notifications", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Создаём тестовый ResponseWriter
			w := httptest.NewRecorder()

			// Вызываем хендлер через middleware
			handler := ErrorMiddleware(h.Create)
			handler(w, req)

			// Проверяем код ответа
			assert.Equal(t, tt.expectedCode, w.Code)

			// Если ожидаем ошибку — проверяем сообщение
			if tt.expectedErr != "" {
				var resp map[string]string
				json.NewDecoder(w.Body).Decode(&resp)
				assert.Equal(t, tt.expectedErr, resp["error"])
			}

			// Если успех — проверяем что уведомление сохранено
			if tt.expectedCode == http.StatusCreated {
				assert.Equal(t, 1, store.Count())
			}
		})
	}
}

// TestGetByID — тесты на получение уведомления
func TestGetByID(t *testing.T) {
	tests := []struct {
		name         string
		id           string
		setupStore   func(store *MockStore)
		expectedCode int
		expectedErr  string
	}{
		{
			name: "успешное получение",
			id:   "test-id-123",
			setupStore: func(store *MockStore) {
				store.Set(model.Notification{
					ID:     "test-id-123",
					UserID: "user-123",
					Type:   model.TypeEmail,
					Title:  "Тест",
					Body:   "Текст",
					Status: model.StatusPending,
				})
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "не найдено",
			id:           "несуществующий-id",
			setupStore:   func(store *MockStore) {},
			expectedCode: http.StatusNotFound,
			expectedErr:  "уведомление несуществующий-id не найдено",
		},
		{
			name:         "пустой id",
			id:           "",
			setupStore:   func(store *MockStore) {},
			expectedCode: http.StatusBadRequest,
			expectedErr:  "id обязателен",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobs := make(chan model.Notification, 10)
			store := NewMockStore()
			tt.setupStore(store)
			h := NewNotificationHandlerWithStore(jobs, store)

			url := "/api/v1/notifications/" + tt.id
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			handler := ErrorMiddleware(h.GetByID)
			handler(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedErr != "" {
				var resp map[string]string
				json.NewDecoder(w.Body).Decode(&resp)
				assert.Equal(t, tt.expectedErr, resp["error"])
			}
		})
	}
}
