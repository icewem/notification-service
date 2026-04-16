package model

import "time"

// Тип уведомления
type NotificationType string

const (
    TypeEmail NotificationType = "email"
    TypePush  NotificationType = "push"
    TypeSMS   NotificationType = "sms"
)

// Статус уведомления
type NotificationStatus string

const (
    StatusPending    NotificationStatus = "pending"
    StatusProcessing NotificationStatus = "processing"
    StatusSent       NotificationStatus = "sent"
    StatusFailed     NotificationStatus = "failed"
)

// Notification — основная модель уведомления
type Notification struct {
    ID        string             `json:"id"`
    UserID    string             `json:"user_id"`
    Type      NotificationType   `json:"type"`
    Title     string             `json:"title"`
    Body      string             `json:"body"`
    Status    NotificationStatus `json:"status"`
    CreatedAt time.Time          `json:"created_at"`
}

// CreateNotificationRequest — запрос на создание уведомления
type CreateNotificationRequest struct {
    UserID string           `json:"user_id"`
    Type   NotificationType `json:"type"`
    Title  string           `json:"title"`
    Body   string           `json:"body"`
}
