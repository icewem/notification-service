package apperror

import "fmt"

// Sentinel errors — заранее определённые ошибки
// используем для проверки через errors.Is()
var (
	ErrNotFound       = fmt.Errorf("не найдено")
	ErrInvalidRequest = fmt.Errorf("невалидный запрос")
	ErrInternal       = fmt.Errorf("внутренняя ошибка")
	ErrTimeout        = fmt.Errorf("таймаут")
)

// AppError — кастомный тип ошибки сервиса
type AppError struct {
	// Code — HTTP код ответа
	Code int
	// Message — сообщение для клиента
	Message string
	// Err — оригинальная ошибка
	Err error
}

// Error — реализуем интерфейс error
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap — позволяет errors.Is() и errors.As() работать
// с вложенными ошибками
func (e *AppError) Unwrap() error {
	return e.Err
}

// New — создать новую AppError
func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NotFound — 404
func NotFound(message string) *AppError {
	return New(404, message, ErrNotFound)
}

// BadRequest — 400
func BadRequest(message string) *AppError {
	return New(400, message, ErrInvalidRequest)
}

// Internal — 500
func Internal(err error) *AppError {
	return New(500, "внутренняя ошибка сервера", fmt.Errorf("%w: %w", ErrInternal, err))
}
