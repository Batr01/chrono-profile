package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"chrono-player-profile/internal/handlers"
	"chrono-player-profile/internal/service"
	"chrono-player-profile/internal/storage"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	// Параметры командной строки
	port := flag.String("port", "8080", "HTTP server port")
	dbDSN := flag.String("db-dsn", "host=localhost user=postgres password=postgres dbname=chrono_profile port=5432 sslmode=disable", "PostgreSQL DSN")
	redisAddr := flag.String("redis-addr", "localhost:6379", "Redis address")
	flag.Parse()

	// Инициализация логгера
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting Player Profile Service",
		zap.String("port", *port),
		zap.String("db-dsn", *dbDSN),
		zap.String("redis-addr", *redisAddr),
	)

	// Инициализация хранилища
	storage, err := storage.NewPostgresStorage(*dbDSN, *redisAddr, logger)
	if err != nil {
		logger.Fatal("Failed to initialize storage", zap.Error(err))
	}
	defer storage.Close()

	// Инициализация сервиса
	playerService := service.NewPlayerService(storage, logger)

	// Инициализация обработчиков
	profileGetHandler := handlers.NewProfileGetHandler(playerService, logger)
	profileUpdateHandler := handlers.NewProfileUpdateHandler(playerService, logger)

	// Настройка роутера
	router := mux.NewRouter()

	// API v1 маршруты
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Профиль игрока
	apiV1.HandleFunc("/profile", profileUpdateHandler.Create).Methods("POST")
	apiV1.HandleFunc("/profile/{id}", profileGetHandler.GetByID).Methods("GET")
	apiV1.HandleFunc("/profile/{id}", profileUpdateHandler.Update).Methods("PUT")
	apiV1.HandleFunc("/profile/{id}", profileUpdateHandler.Delete).Methods("DELETE")
	apiV1.HandleFunc("/profile/nickname/{nickname}", profileGetHandler.GetByNickname).Methods("GET")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	// Запуск HTTP сервера
	serverAddr := fmt.Sprintf(":%s", *port)
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		logger.Info("HTTP server started", zap.String("address", serverAddr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	if err := srv.Close(); err != nil {
		logger.Error("Error closing server", zap.Error(err))
	}

	logger.Info("Server stopped")
}

