package util

import (
	"errors"
	"log/slog"
	"regexp"
)

// Регулярные выражения для извлечения частей URL
var (
	// Регулярное выражение для извлечения сервера из URL (например, s1, s2 и т.д.)
	serverRegex = regexp.MustCompile(`https?://(s\d+)\.origin-cluster`) // Поддержка http и https

	// Регулярное выражение для извлечения пути из URL (например, video/123/xcg2djHckad.m3u8)
	pathRegex = regexp.MustCompile(`https?://s\d+\.origin-cluster/(.*)`) // Поддержка http и https
)

// ParseVideoURL разбирает входной URL и возвращает сервер и путь
// url - входной URL в формате https://s1.origin-cluster/video/123/xcg2djHckad.m3u8
// Возвращает:
// - server (например, s1)
// - path (например, video/123/xcg2djHckad.m3u8)
// - err (ошибка, если URL не может быть разобран)
func ParseVideoURL(url string) (server string, path string, err error) {
	// Извлекаем сервер из URL с использованием регулярного выражения
	serverMatch := serverRegex.FindStringSubmatch(url)
	if len(serverMatch) < 2 {
		// Логируем ошибку, если сервер не может быть извлечен из URL
		slog.Error("Ошибка при разборе URL: не удалось извлечь сервер", "url", url)
		// Возвращаем ошибку, если сервер не найден
		return "", "", errors.New("не удалось извлечь сервер из URL")
	}
	// Присваиваем извлеченное значение серверу
	server = serverMatch[1]

	// Извлекаем путь из URL с использованием регулярного выражения
	pathMatch := pathRegex.FindStringSubmatch(url)
	if len(pathMatch) < 2 {
		// Логируем ошибку, если путь не может быть извлечен из URL
		slog.Error("Ошибка при разборе URL: не удалось извлечь путь", "url", url)
		// Возвращаем ошибку, если путь не найден
		return "", "", errors.New("не удалось извлечь путь из URL")
	}
	// Присваиваем извлеченное значение пути
	path = pathMatch[1]

	// Логируем успешное извлечение данных из URL
	slog.Info("URL успешно разобран", "сервер", server, "путь", path)

	// Возвращаем извлеченные данные и ошибку (если она есть)
	return server, path, nil
}
