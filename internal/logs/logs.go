package logs

import (
	"log/slog"
	"time"
)

type LogMessage struct {
	Level slog.Level
	Msg   string
	Args  []interface{}
}

var (
	maxChannelSize = 100000                                // Размер канала логирования
	logChannel     = make(chan LogMessage, maxChannelSize) // Канал для асинхронного логирования
	logBuffer      = make([]LogMessage, 0, maxChannelSize) // Буфер для накопления логов
)

// Функция для ленивой инициализации канала логирования
func getLogChannel() chan LogMessage {
	if logChannel == nil {
		// Инициализация канала только при необходимости
		logChannel = make(chan LogMessage, maxChannelSize)
	}
	return logChannel
}

// Асинхронное логирование
func AsyncLog(level slog.Level, msg string, args ...interface{}) {
	// Ограничение размера данных в логах
	for i, arg := range args {
		switch v := arg.(type) {
		case string:
			if len(v) > 100 { // Ограничиваем строковые данные до 100 символов
				args[i] = v[:100]
			}
		}
	}

	logMsg := LogMessage{Level: level, Msg: msg, Args: args}
	select {
	case getLogChannel() <- logMsg:
	default:
		// Если канал переполнен, пробуем повторить через несколько миллисекунд
		time.Sleep(100 * time.Millisecond)
		select {
		case getLogChannel() <- logMsg:
		default:
			slog.Warn("Не удалось записать лог, канал переполнен", "msg", msg)
		}
	}
}

func processLogs(logs []LogMessage) {
	for _, logMsg := range logs {
		switch logMsg.Level {
		case slog.LevelInfo:
			slog.Info(logMsg.Msg, logMsg.Args...)
		case slog.LevelError:
			slog.Error(logMsg.Msg, logMsg.Args...)
		default:
			slog.Debug(logMsg.Msg, logMsg.Args...)
		}
	}
}

func init() {
	// Создаем несколько горутин для обработки логов
	// Горутин для асинхронного логирования
	go func() {
		for logMsg := range logChannel {
			logBuffer = append(logBuffer, logMsg)

			// Когда количество логов в буфере достигает максимума, обрабатываем их
			if len(logBuffer) >= maxChannelSize {
				processLogs(logBuffer)
				logBuffer = nil // Очищаем буфер
			}
		}
	}()
}
