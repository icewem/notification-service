package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/icewem/notification-service/internal/handler"
	"github.com/icewem/notification-service/internal/model"
	"github.com/icewem/notification-service/internal/service"
)

func main() {
	jobs := make(chan model.Notification, 100)

	pool := service.NewWorkerPool(jobs, 3)
	pool.Start()

	notificationHandler := handler.NewNotificationHandler(jobs)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"status": "ok"}`)
	})

	// Оборачиваем хендлеры в ErrorMiddleware
	mux.HandleFunc("/api/v1/notifications", handler.ErrorMiddleware(notificationHandler.Create))
	mux.HandleFunc("/api/v1/notifications/", handler.ErrorMiddleware(notificationHandler.GetByID))

	// Регистрируем pprof эндпоинты
	handler.RegisterDebugHandlers(mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Сервис запущен на :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска: %v", err)
		}
	}()

	<-quit
	log.Println("Получен сигнал остановки...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка остановки: %v", err)
	}
	log.Println("HTTP сервер остановлен")

	close(jobs)
	pool.Stop()

	log.Println("Сервис остановлен корректно")
}
