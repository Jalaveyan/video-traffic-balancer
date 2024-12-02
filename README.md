
# Video Traffic Load Balancer

## Описание
**Video Traffic Load Balancer** — это высокопроизводительный gRPC-сервер, который перенаправляет запросы на видео на оптимальные серверы (например, CDN). Сервер поддерживает балансировку нагрузки, кэширование и асинхронное логирование, а также имеет встроенные механизмы для мониторинга производительности.

## Основные возможности
- **gRPC API** для балансировки трафика.
- Кэширование запросов с использованием LRU-алгоритма.
- **Пул горутин** для ограничения ресурсов и повышения производительности.
- Поддержка **health checks** (gRPC и HTTP).
- Интеграция с профилировщиком **pprof**.
- Асинхронное логирование через канал для уменьшения задержек.

## Требования
- **Go** 1.20+
- **Protocol Buffers** 3.15+
- gRPC-библиотеки: `google.golang.org/grpc`, `google.golang.org/grpc/health`.
- LRU-кэш: `github.com/hashicorp/golang-lru`.

## Установка и запуск

### 1. Клонирование репозитория
```bash
git clone https://github.com/yourusername/videobalance.git
cd videobalance
```

### 2. Установка зависимостей
```bash
go mod tidy
```

### 3. Переменные окружения
Убедитесь, что переменные окружения настроены:
- `CDN_HOST` — адрес CDN (по умолчанию `cdn.example.com`).
- `SERVER_PORT` — порт gRPC сервера (по умолчанию `:443`).

Пример:
```bash
export CDN_HOST=cdn.example.com
export SERVER_PORT=:443
```

### 4. Запуск сервера
```bash
go run cmd/server/main.go
```

### 5. Мониторинг
- **HTTP health check:** доступен по адресу `http://localhost:8080/health`.
- **pprof:** доступен по адресу `http://localhost:6060/debug/pprof/`.

## gRPC API

### Метод `Redirect`
#### Запрос
```protobuf
message RedirectRequest {
  string video = 1; // URL видео для перенаправления.
}
```

#### Ответ
```protobuf
message RedirectResponse {
  string target_url = 1; // Перенаправленный URL.
}
```

### Пример gRPC-запроса с использованием grpcurl:
```bash
grpcurl -d '{"video": "https://s1.origin-cluster/video/123/xcg2djHckad.m3u8"}'   -plaintext localhost:443 videobalance.Balancer/Redirect
```

## Структура проекта
```
videobalance/
│
├── cmd/
│   └── server/         # Точка входа для запуска gRPC сервера
├── internal/
│   ├── cache/          # Модуль для управления LRU-кэшем
│   ├── config/         # Загрузка и обработка конфигурации
│   ├── logs/           # Асинхронное логирование
│   ├── server/         # Логика gRPC сервера
│   ├── util/           # Вспомогательные функции
│   └── worker/         # Управление пулом горутин
├── proto/              # gRPC-протоколы и сообщения
│   └── balancer.proto  # Файлы описания API
├── go.mod              # Управление зависимостями Go
└── README.md           # Документация
```
