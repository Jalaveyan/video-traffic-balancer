package config

import (
	"log/slog"
	"os"
)

// Config представляет конфигурацию приложения
type Config struct {
	CDNHost    string // Хост CDN, используемый для передачи данных
	ServerPort string // Порт для gRPC сервера, по которому сервер будет принимать соединения
}

// LoadConfig считывает переменные окружения или использует значения по умолчанию
// Эта функция пытается загрузить значения конфигурации из переменных окружения.
// Если переменная окружения не задана, то используется значение по умолчанию.
func LoadConfig() (*Config, error) {
	// Получаем значение переменной окружения CDN_HOST.
	cdnHost := os.Getenv("CDN_HOST")
	if cdnHost == "" {
		// Если переменная не задана, используем значение по умолчанию и предупреждаем об этом.
		cdnHost = "cdn.example.com" // Значение по умолчанию для CDN хоста
		slog.Warn("Переменная CDN_HOST не задана. Используется значение по умолчанию", "CDN_HOST", cdnHost)
	} else {
		// Если переменная задана, логируем её значение.
		slog.Info("Переменная CDN_HOST загружена", "CDN_HOST", cdnHost)
	}

	// Получаем значение переменной окружения SERVER_PORT.
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		// Если переменная не задана, используем значение по умолчанию и предупреждаем об этом.
		serverPort = ":443" // Значение по умолчанию для порта
		slog.Warn("Переменная SERVER_PORT не задана. Используется значение по умолчанию", "SERVER_PORT", serverPort)
	} else {
		// Если переменная задана, логируем её значение.
		slog.Info("Переменная SERVER_PORT загружена", "SERVER_PORT", serverPort)
	}

	// Возвращаем структуру конфигурации с загруженными значениями.
	return &Config{
		CDNHost:    cdnHost,
		ServerPort: serverPort,
	}, nil
}
