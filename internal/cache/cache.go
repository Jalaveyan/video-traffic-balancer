package cache

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"hash/fnv"
	"log/slog"
	"sync"
	"time"
)

const (
	cacheTTL                = 10 * time.Minute // Время жизни кэша (TTL)
	maxCacheSize            = 5000             // Максимальный размер кэша
	cacheGCInterval         = 5 * time.Minute  // Интервал для очистки кэша
	frequentAccessThreshold = 100              // Порог запросов, после которого URL остаётся дольше в кэше
)

var (
	UrlCache         *lru.Cache           // LRU кэш для хранения URL
	cacheTimestamp   map[string]time.Time // Карта для хранения времени для каждого ключа
	cacheAccessCount map[string]int       // Карта для подсчёта количества запросов на видео
	cacheMutex       sync.RWMutex         // Мьютекс для синхронизации доступа
	CleanerStop      chan struct{}        // Канал для остановки очистки
)

type CacheEntry struct {
	URL       string
	Timestamp time.Time
}

// Инициализация и очистка кэша по TTL
func init() {
	var err error
	UrlCache, err = lru.New(maxCacheSize)
	if err != nil {
		slog.Info("Ошибка при создании кэша: %v", err)
	}

	cacheTimestamp = make(map[string]time.Time)
	cacheAccessCount = make(map[string]int)

	// Инициализация канала для остановки очистки
	CleanerStop = make(chan struct{})

	// Запуск горутины для очистки устаревших данных
	go startCacheCleaner()
}

// Функция для хеширования видео URL
func hashVideo(video string) string {
	h := fnv.New64a()
	h.Write([]byte(video))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Очистка кэша от устаревших записей
func startCacheCleaner() {
	ticker := time.NewTicker(cacheGCInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cleanExpiredCache()
		case <-CleanerStop: // Ожидаем закрытия канала
			slog.Info("Остановка очистки кэша.")
			return
		}
	}
}

// Очистка устаревших записей из кэша
func cleanExpiredCache() {
	now := time.Now()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Проходим по всем ключам и удаляем устаревшие записи
	for key, timestamp := range cacheTimestamp {
		if now.Sub(timestamp) > cacheTTL {
			// Удаление устаревших элементов
			UrlCache.Remove(key)
			delete(cacheTimestamp, key)
			delete(cacheAccessCount, key)
			slog.Info("Удалена устаревшая запись из кэша: %v", key)
		}
	}
}

// Получение URL из кэша
func GetFromCache(video string) (string, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	hashedVideo := hashVideo(video)
	entryInterface, found := UrlCache.Get(hashedVideo)
	if !found || time.Since(entryInterface.(CacheEntry).Timestamp) > cacheTTL {
		return "", false
	}

	cacheAccessCount[video]++
	if cacheAccessCount[video] > frequentAccessThreshold {
		cacheTimestamp[video] = time.Now().Add(cacheTTL)
	}
	return entryInterface.(CacheEntry).URL, true
}
