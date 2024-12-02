# Используем официальное изображение Golang версии 1.23 для сборки
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект в контейнер
COPY . .

# Собираем бинарный файл приложения
RUN go build -o /video-balancer ./cmd/balancer/main.go

# Используем официальное изображение Golang 1.23 как базовый образ для запуска
FROM golang:1.23

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из стадии сборки
COPY --from=builder /video-balancer .

# Устанавливаем переменные окружения с возможностью их изменения при запуске
ENV CDN_HOST=cdn.example.com
ENV SERVER_PORT=:443

# Открываем порт 50051 для gRPC сервера
EXPOSE 443

# Запускаем приложение
CMD ["./video-balancer"]
