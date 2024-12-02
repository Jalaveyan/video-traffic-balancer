package worker

import (
	"golang.org/x/net/context"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
	"videobalance/internal/logs"
)

var (
	maxGoroutines    int64         = 2000 // Максимальное количество активных горутин
	activeGoroutines int64                // Счётчик текущих активных горутин
	stopMonitoring   chan struct{}        // Канал для остановки мониторинга
	taskChannel      chan struct{}        // Канал для задач
	once             sync.Once
)

// Инициализация и запуск мониторинга пула горутин
func AdjustWorkerPoolSize() {
	taskChannel = make(chan struct{}, maxGoroutines)
	stopMonitoring = make(chan struct{})
	go monitorPoolSize(stopMonitoring)
}

// Функция мониторинга пула горутин с задержкой между проверками
func monitorPoolSize(stop chan struct{}) {
	ticker := time.NewTicker(2 * time.Second) // Логируем раз в 2 секунды для уменьшения нагрузки
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			logs.AsyncLog(slog.LevelInfo, "Остановка мониторинга пула горутин")
			return
		case <-ticker.C:
			currentGoroutines := atomic.LoadInt64(&activeGoroutines)
			// Логика увеличения или уменьшения размера пула горутин
			if currentGoroutines > maxGoroutines {
				logs.AsyncLog(slog.LevelInfo, "Уменьшаем размер пула горутин", "current", currentGoroutines)
				atomic.StoreInt64(&maxGoroutines, currentGoroutines)
			} else if currentGoroutines < maxGoroutines {
				logs.AsyncLog(slog.LevelInfo, "Размер пула горутин в норме", "current", currentGoroutines)
			} else {
				logs.AsyncLog(slog.LevelInfo, "Размер пула горутин в норме", "current", currentGoroutines)
			}
		}
	}
}

// Завершение мониторинга и завершение работы с горутинами
func Shutdown() {
	once.Do(func() {
		close(stopMonitoring) // Закрываем канал для остановки мониторинга горутин
		logs.AsyncLog(slog.LevelInfo, "Мониторинг горутин остановлен")
	})
}

// Структура Worker (для примера)
type Worker struct {
	Done chan bool
}

func (w *Worker) DoWork() {
	// Симуляция работы
	time.Sleep(100 * time.Millisecond) // Эмулируем работу
	w.Done <- true
}

// Выполнение задачи с таймаутом
func executeTask(worker *Worker) {
	select {
	case taskChannel <- struct{}{}:
		go func() {
			defer func() { <-taskChannel }()
			runWithTimeout(worker) // выполняем работу с таймаутом
		}()
	default:
		logs.AsyncLog(slog.LevelWarn, "Превышен лимит горутин, задача отклонена")
	}
}

// Запуск задачи с контролем таймаута
func runWithTimeout(worker *Worker) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		logs.AsyncLog(slog.LevelError, "Горутина не успела завершиться вовремя", "timeout", ctx.Err())
	case workerDone := <-worker.Done:
		if workerDone {
			logs.AsyncLog(slog.LevelInfo, "Горутина завершена вовремя")
		}
	}
}

// Основная функция, которая демонстрирует использование worker
func RunWorkers() {
	// Создаем работников
	workers := make([]*Worker, 5) // Например, 5 работников

	for i := 0; i < len(workers); i++ {
		workers[i] = &Worker{
			Done: make(chan bool),
		}
	}

	// Запускаем выполнение задач для работников
	for _, worker := range workers {
		executeTask(worker) // Отправляем задачу на выполнение
	}

	// Ждем, пока все горутины завершат работу
	time.Sleep(10 * time.Second)

	// Завершаем работу пула горутин и мониторинга
	Shutdown()
}
