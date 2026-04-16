package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"syscall"

	"github.com/icewem/notification-service/internal/handler"
	"github.com/icewem/notification-service/internal/model"
	"github.com/icewem/notification-service/internal/service"
)

func main() {
	// Канал для задач — буфер 100 уведомлений
	jobs := make(chan model.Notification, 100)

	// Запускаем worker pool — 3 воркера
	pool := service.NewWorkerPool(jobs, 3)
	pool.Start()

	// Создаём хендлер — передаём канал jobs
	notificationHandler := handler.NewNotificationHandler(jobs)

	// Роутер
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "ok"}`)
	})

	// Уведомления
	mux.HandleFunc("/api/v1/notifications", notificationHandler.Create)
	mux.HandleFunc("/api/v1/notifications/", notificationHandler.GetByID)

	// HTTP сервер
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Graceful shutdown — ловим сигналы остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в горутине
	go func() {
		log.Println("Сервис запущен на :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ждём сигнала остановки
	<-quit
	log.Println("Получен сигнал остановки...")

	// Останавливаем HTTP сервер
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка остановки сервера: %v", err)
	}

	// Закрываем канал jobs — воркеры завершат текущие задачи и остановятся
	close(jobs)
	pool.Stop()

	log.Println("Сервис остановлен")
}
