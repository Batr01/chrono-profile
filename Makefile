.PHONY: build run test clean docker-up docker-down help

# Переменные
BINARY_NAME=pps
MAIN_PATH=cmd/main.go

help: ## Показать эту справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Собрать приложение
	@echo "Сборка приложения..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Готово: $(BINARY_NAME)"

run: ## Запустить приложение
	@echo "Запуск приложения..."
	@go run $(MAIN_PATH)

run-custom: ## Запустить с кастомными параметрами
	@go run $(MAIN_PATH) -port=8080 \
		-db-dsn="host=localhost user=postgres password=postgres dbname=chrono_profile port=5432 sslmode=disable" \
		-redis-addr="localhost:6379"

test: ## Запустить тесты
	@echo "Запуск тестов..."
	@go test -v ./...

test-coverage: ## Запустить тесты с покрытием
	@echo "Запуск тестов с покрытием..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

clean: ## Очистить собранные файлы
	@echo "Очистка..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out
	@go clean

deps: ## Установить зависимости
	@echo "Установка зависимостей..."
	@go mod download
	@go mod tidy

docker-up: ## Запустить PostgreSQL и Redis в Docker
	@echo "Запуск Docker контейнеров..."
	@docker-compose up -d
	@echo "Ожидание готовности PostgreSQL..."
	@sleep 5

docker-down: ## Остановить Docker контейнеры
	@echo "Остановка Docker контейнеров..."
	@docker-compose down

docker-logs: ## Показать логи Docker контейнеров
	@docker-compose logs -f

lint: ## Запустить линтер
	@echo "Проверка кода..."
	@golangci-lint run ./...

fmt: ## Форматировать код
	@echo "Форматирование кода..."
	@go fmt ./...

