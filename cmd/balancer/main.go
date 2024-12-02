package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"videobalance/internal/worker"
	pb "videobalance/proto"

	"log"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof" // Инициализация пакета pprof
	"os"
	"os/signal"
	"syscall"
	"time"
	"videobalance/internal/config"
	"videobalance/internal/server"
	_ "videobalance/proto"
)

// healthCheckHandler для проверки состояния через HTTP
func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	// Это простая проверка состояния
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	// Логирование для структурированных логов
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	worker.AdjustWorkerPoolSize()

	// Настроим pprof
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil)) // Порт для pprof
	}()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Ошибка загрузки конфигурации", "ошибка", err)
		return
	}
	slog.Info("Конфигурация загружена", "CDN_HOST", cfg.CDNHost, "SERVER_PORT", cfg.ServerPort)

	// Настройка gRPC сервера
	lis, err := net.Listen("tcp", cfg.ServerPort)
	if err != nil {
		slog.Error("Ошибка при создании соединения", "ошибка", err, "порт", cfg.ServerPort)
		return
	}
	slog.Info("gRPC сервер слушает порт", "порт", cfg.ServerPort)

	// Создание нового экземпляра сервера балансировщика
	balancerServer := server.NewBalancerServer("balancer-domain.com", cfg.CDNHost)

	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(200000),
		grpc.MaxRecvMsgSize(100*1024*1024),
		grpc.MaxSendMsgSize(100*1024*1024),
		grpc.WriteBufferSize(256*1024*1024),
		grpc.ReadBufferSize(256*1024*1024),
	)

	// Регистрация сервиса
	pb.RegisterBalancerServer(grpcServer, balancerServer)

	// Настройка Health Check для gRPC
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	// Устанавливаем статус "SERVING" для нашего сервиса
	healthServer.SetServingStatus("videobalance", grpc_health_v1.HealthCheckResponse_SERVING)

	// Инициализация пула горутин и запуск мониторинга
	worker.RunWorkers() // Запускаем рабочие горутины

	// Запуск gRPC сервера в горутине
	go func() {
		startTime := time.Now()
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("Ошибка запуска gRPC сервера", "ошибка", err)
		}
		duration := time.Since(startTime)
		slog.Info("gRPC сервер завершил работу", "время_работы", duration)
	}()

	// Запуск HTTP сервера для health check
	go func() {
		// Добавление обработчика health check
		http.HandleFunc("/health", healthCheckHandler) // Путь для health check
		log.Fatal(http.ListenAndServe(":8080", nil))   // Порт для health check
	}()

	// Канал для graceful shutdown
	stopChan := make(chan struct{}, 1)

	// Обработка сигналов ОС для graceful shutdown
	go func() {
		// Ожидаем сигнал завершения работы
		stopSignal := make(chan os.Signal, 1)
		signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)

		<-stopSignal // Блокируемся, пока не получим сигнал
		slog.Info("Получен сигнал завершения работы, начинаем graceful shutdown")

		// Завершаем сервер
		grpcServer.GracefulStop()

		// Завершаем работу пула горутин и gRPC сервера
		worker.Shutdown() // Завершаем мониторинг горутин

		// Закрытие канала graceful shutdown
		close(stopChan)
	}()

	// Задержка ожидания сигнала завершения
	<-stopChan

	// Закрытие других ресурсов, если необходимо

	slog.Info("Остановка сервиса завершена")
}
