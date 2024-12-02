package server

import (
	"context"
	"fmt"
	_ "github.com/hashicorp/golang-lru"
	"golang.org/x/sync/semaphore"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
	"videobalance/internal/cache"
	"videobalance/internal/util"
	"videobalance/internal/worker"
	pb "videobalance/proto"
)

const (
	maxConcurrentRequests = 5000 // Максимальное количество параллельных запросов
	defaultWorkerPoolSize = 500  // Начальный размер пула горутин
)

var (
	// Инициализация семафора для ограничения количества параллельных запросов
	sem = semaphore.NewWeighted(maxConcurrentRequests)

	// Пул горутин для обработки запросов
	workerPool = make(chan struct{}, defaultWorkerPoolSize)

	// Счетчики запросов по каждому видео
	videoRequestCounts sync.Map
)

// Балансировщик запросов
type BalancerServer struct {
	pb.UnimplementedBalancerServer
	balancerDomain string
	cdnHost        string
	logger         *slog.Logger
	mu             sync.Mutex // для защиты локального счетчика от гонок
}

// Конструктор балансировщика
func NewBalancerServer(balancerDomain, cdnHost string) *BalancerServer {
	worker.AdjustWorkerPoolSize()

	return &BalancerServer{
		balancerDomain: balancerDomain,
		cdnHost:        cdnHost,
		logger:         slog.Default(),
	}
}

func (s *BalancerServer) Redirect(ctx context.Context, req *pb.RedirectRequest) (*pb.RedirectResponse, error) {
	// Устанавливаем тайм-аут для обработки запроса
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Попытка захватить семафор для ограничения параллельных запросов
	if err := sem.Acquire(ctx, 1); err != nil {
		s.logger.Warn("Не удалось захватить семафор", "error", err, "video", req.Video)
		return nil, fmt.Errorf("не удалось захватить семафор: %v", err)
	}
	defer sem.Release(1) // Освобождаем семафор

	// Попытка захватить место в пуле горутин
	select {
	case workerPool <- struct{}{}: // Блокируем горутину, если пул заполнился
		defer func() { <-workerPool }() // Освобождаем место в пуле
	case <-ctx.Done():
		s.logger.Warn("Пул горутин заполнен, запрос ожидает")
		return nil, fmt.Errorf("запрос был отменен или превышен тайм-аут")
	}

	// Проверка наличия URL в кэше
	if cachedURL, found := cache.GetFromCache(req.Video); found {
		s.logger.Info("URL найден в кэше", "url", cachedURL)
		return &pb.RedirectResponse{TargetUrl: cachedURL}, nil

	}

	// Используем функцию из util для разбора видео URL
	server, path, err := util.ParseVideoURL(req.Video)
	if err != nil {
		s.logger.Error("Не удалось разобрать URL", "url", req.Video, "error", err)
		return nil, err
	}

	// Получаем текущий счетчик запросов
	count := s.incrementRequestCount(req.Video)

	// Логика перенаправления для каждого 10-го запроса
	if count%10 == 0 {
		// Печатаем логи для диагностики
		s.logger.Info("Перенаправление на оригинальный URL", "url", req.Video, "номер_запроса", count)
		return &pb.RedirectResponse{TargetUrl: req.Video}, nil
	}

	// Если cdnHost пуст, используем оригинальный URL
	if s.cdnHost == "" {

		s.logger.Info("Перенаправление на оригинальный URL, CDN не указан", "url", req.Video)

		return &pb.RedirectResponse{TargetUrl: req.Video}, nil
	}

	// Формируем URL для перенаправления на CDN
	cdnURL := fmt.Sprintf("http://%s/%s", s.cdnHost, server, path)

	s.logger.Info("Перенаправление на CDN", "url", cdnURL)

	// Кэшируем URL
	return &pb.RedirectResponse{TargetUrl: cdnURL}, nil
}

// Получение и обновление локального счетчика запросов
func (s *BalancerServer) incrementRequestCount(video string) uint64 {
	countInterface, _ := videoRequestCounts.LoadOrStore(video, uint64(0))
	count := countInterface.(uint64)
	newCount := atomic.AddUint64(&count, 1)
	videoRequestCounts.Store(video, newCount)
	// Возвращаем новое значение счетчика
	return newCount
}
