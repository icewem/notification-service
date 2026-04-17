package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/icewem/notification-service/internal/model"
)

// WorkerPool — пул воркеров для обработки уведомлений
type WorkerPool struct {
	jobs    <-chan model.Notification
	workers int
	wg      sync.WaitGroup
}

// NewWorkerPool — конструктор пула воркеров
func NewWorkerPool(jobs <-chan model.Notification, workers int) *WorkerPool {
	return &WorkerPool{
		jobs:    jobs,
		workers: workers,
	}
}

// Start — запускаем воркеры
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	fmt.Printf("Worker pool запущен: %d воркеров\n", wp.workers)
}

// Stop — ждём завершения всех воркеров
func (wp *WorkerPool) Stop() {
	wp.wg.Wait()
	fmt.Println("Worker pool остановлен")
}

// worker — читает задачи из канала
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for n := range wp.jobs {
		// Создаём контекст с таймаутом на обработку
		// если за 5 секунд не успели — отменяем
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		if err := wp.process(ctx, id, n); err != nil {
			fmt.Printf("воркер %d ошибка обработки id=%s: %v\n", id, n.ID, err)
		}

		// Всегда освобождаем ресурсы контекста
		cancel()
	}
}

// process — обрабатываем одно уведомление
func (wp *WorkerPool) process(ctx context.Context, workerID int, n model.Notification) error {
	fmt.Printf(
		"воркер %d обрабатывает id=%s type=%s user=%s\n",
		workerID, n.ID, n.Type, n.UserID,
	)

	// Эмулируем отправку через select с ctx.Done()
	// чтобы можно было прервать если таймаут истёк
	select {
	case <-time.After(100 * time.Millisecond): // эмуляция работы
		fmt.Printf("воркер %d отправил id=%s title=%s\n", workerID, n.ID, n.Title)
		return nil
	case <-ctx.Done():
		// таймаут истёк или контекст отменён
		return fmt.Errorf("таймаут обработки: %w", ctx.Err())
	}
}
