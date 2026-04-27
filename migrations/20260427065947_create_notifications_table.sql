-- +goose Up
-- +goose StatementBegin
CREATE TABLE notifications (
    id          TEXT PRIMARY KEY,
    user_id     TEXT NOT NULL,
    type        TEXT NOT NULL CHECK (type IN ('email', 'push', 'sms')),
    title       TEXT NOT NULL,
    body        TEXT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'pending'
                CHECK (status IN ('pending', 'processing', 'sent', 'failed')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индекс для поиска уведомлений пользователя
CREATE INDEX idx_notifications_user_id ON notifications(user_id);

-- Составной индекс для фильтрации по пользователю и статусу
CREATE INDEX idx_notifications_user_status ON notifications(user_id, status);

-- Partial индекс — только pending уведомления
-- используется когда ищем необработанные задачи
CREATE INDEX idx_notifications_pending ON notifications(created_at)
WHERE status = 'pending';

-- BRIN индекс для временных меток — большие таблицы
CREATE INDEX idx_notifications_created_brin ON notifications
USING brin(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications;
-- +goose StatementEnd
