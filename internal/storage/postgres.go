package storage

import (
	"context"
	"fmt"
	"time"

	"chrono-player-profile/internal/models"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresStorage представляет хранилище данных PostgreSQL
type PostgresStorage struct {
	db     *gorm.DB
	cache  *redis.Client
	logger *zap.Logger
}

// NewPostgresStorage создает новое хранилище PostgreSQL
func NewPostgresStorage(dsn string, redisAddr string, logger *zap.Logger) (*PostgresStorage, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Автомиграция моделей
	if err := db.AutoMigrate(&models.Player{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Подключение к Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // без пароля
		DB:       0,  // используем DB по умолчанию
	})

	// Проверка подключения к Redis
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Warn("Failed to connect to Redis, continuing without cache", zap.Error(err))
		rdb = nil
	}

	return &PostgresStorage{
		db:     db,
		cache:  rdb,
		logger: logger,
	}, nil
}

// GetPlayerByID получает игрока по ID с кешированием
func (s *PostgresStorage) GetPlayerByID(id string) (*models.Player, error) {
	// Пытаемся получить из кеша
	if s.cache != nil {
		cacheKey := fmt.Sprintf("player:%s", id)
		ctx := context.Background()
		
		// Проверяем кеш для ELO (быстрый доступ)
		eloStr, err := s.cache.Get(ctx, fmt.Sprintf("%s:elo", cacheKey)).Result()
		if err == nil {
			// Если ELO в кеше, можно сделать быстрый ответ (опционально)
			s.logger.Debug("ELO found in cache", zap.String("player_id", id))
		}
		_ = eloStr // используем если нужно
	}

	var player models.Player
	if err := s.db.Where("id = ?", id).First(&player).Error; err != nil {
		return nil, err
	}

	// Сохраняем ELO в кеш для быстрого доступа
	if s.cache != nil {
		cacheKey := fmt.Sprintf("player:%s", id)
		ctx := context.Background()
		s.cache.Set(ctx, fmt.Sprintf("%s:elo", cacheKey), player.ELO, 5*time.Minute)
	}

	return &player, nil
}

// GetPlayerByNickname получает игрока по никнейму
func (s *PostgresStorage) GetPlayerByNickname(nickname string) (*models.Player, error) {
	var player models.Player
	if err := s.db.Where("nickname = ?", nickname).First(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

// CreatePlayer создает нового игрока
func (s *PostgresStorage) CreatePlayer(player *models.Player) error {
	return s.db.Create(player).Error
}

// UpdatePlayer обновляет данные игрока
func (s *PostgresStorage) UpdatePlayer(id string, updates *models.PlayerUpdateRequest) (*models.Player, error) {
	var player models.Player
	if err := s.db.Where("id = ?", id).First(&player).Error; err != nil {
		return nil, err
	}

	// Обновляем поля, если они указаны
	if updates.Nickname != nil {
		player.Nickname = *updates.Nickname
	}
	if updates.Level != nil {
		player.Level = *updates.Level
	}
	if updates.Rating != nil {
		player.Rating = *updates.Rating
	}
	if updates.ELO != nil {
		player.ELO = *updates.ELO
	}
	if updates.Role != nil {
		player.Role = *updates.Role
	}
	if updates.Region != nil {
		player.Region = *updates.Region
	}
	if updates.Language != nil {
		player.Language = *updates.Language
	}
	if updates.Wins != nil {
		player.Wins = *updates.Wins
	}
	if updates.Losses != nil {
		player.Losses = *updates.Losses
	}
	if updates.Rank != nil {
		player.Rank = *updates.Rank
	}
	if updates.Cosmetics != nil {
		player.Cosmetics = updates.Cosmetics
	}
	if updates.Settings != nil {
		player.Settings = updates.Settings
	}
	if updates.PreferredMode != nil {
		player.PreferredMode = *updates.PreferredMode
	}
	if updates.PreferredRole != nil {
		player.PreferredRole = *updates.PreferredRole
	}

	if err := s.db.Save(&player).Error; err != nil {
		return nil, err
	}

	// Обновляем кеш
	if s.cache != nil {
		cacheKey := fmt.Sprintf("player:%s", id)
		ctx := context.Background()
		s.cache.Set(ctx, fmt.Sprintf("%s:elo", cacheKey), player.ELO, 5*time.Minute)
	}

	return &player, nil
}

// DeletePlayer удаляет игрока (мягкое удаление)
func (s *PostgresStorage) DeletePlayer(id string) error {
	return s.db.Where("id = ?", id).Delete(&models.Player{}).Error
}

// Close закрывает соединения
func (s *PostgresStorage) Close() error {
	if s.cache != nil {
		return s.cache.Close()
	}
	return nil
}

