package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/icewem/notification-service/internal/model"
)

// WorkerPool — пул воркеров для обработки уведомлений
type WorkerPool struct {
	// jobs — канал из которого воркеры читают задачи
	jobs    <-chan model.Notification
	// workers — количество воркеров
	workers int
	// wg — ждём завершения всех воркеров
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

// worker — один воркер читает задачи из канала
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for n := range wp.jobs {
		// Читаем задачу из канала и обрабатываем
		wp.process(id, n)
	}
}

// process — обрабатываем одно уведомление
func (wp *WorkerPool) process(workerID int, n model.Notification) {
	fmt.Printf(
		"воркер %d обрабатывает уведомление id=%s type=%s user=%s\n",
		workerID, n.ID, n.Type, n.UserID,
	)

	// Эмулируем отправку уведомления
	time.Sleep(100 * time.Millisecond)

	fmt.Printf(
		"воркер %d отправил уведомление id=%s title=%s\n",
		workerID, n.ID, n.Title,
	)
}
