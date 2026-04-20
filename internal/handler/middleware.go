package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/icewem/notification-service/internal/apperror"
)

// ErrorResponse — структура ответа при ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}

// HandlerFunc — тип хендлера который возвращает ошибку
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// ErrorMiddleware — оборачивает хендлер и обрабатывает ошибки
func ErrorMiddleware(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}

		// Проверяем является ли ошибка AppError
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			// Известная ошибка — отвечаем клиенту
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(appErr.Code)
			json.NewEncoder(w).Encode(ErrorResponse{Error: appErr.Message})
			return
		}

		// Неизвестная ошибка — логируем и отвечаем 500
		log.Printf("неожиданная ошибка: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "внутренняя ошибка сервера"})
	}
}
