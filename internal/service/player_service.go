package service

import (
	"fmt"

	"chrono-player-profile/internal/models"
	"chrono-player-profile/internal/storage"

	"go.uber.org/zap"
)

// PlayerService представляет сервисный слой для работы с игроками
type PlayerService struct {
	storage *storage.PostgresStorage
	logger  *zap.Logger
}

// NewPlayerService создает новый сервис игроков
func NewPlayerService(storage *storage.PostgresStorage, logger *zap.Logger) *PlayerService {
	return &PlayerService{
		storage: storage,
		logger:  logger,
	}
}

// GetPlayerByID получает игрока по ID
func (s *PlayerService) GetPlayerByID(id string) (*models.Player, error) {
	if id == "" {
		return nil, fmt.Errorf("player ID is required")
	}

	player, err := s.storage.GetPlayerByID(id)
	if err != nil {
		s.logger.Error("Failed to get player by ID", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("player not found: %w", err)
	}

	return player, nil
}

// GetPlayerByNickname получает игрока по никнейму
func (s *PlayerService) GetPlayerByNickname(nickname string) (*models.Player, error) {
	if nickname == "" {
		return nil, fmt.Errorf("nickname is required")
	}

	player, err := s.storage.GetPlayerByNickname(nickname)
	if err != nil {
		s.logger.Error("Failed to get player by nickname", zap.String("nickname", nickname), zap.Error(err))
		return nil, fmt.Errorf("player not found: %w", err)
	}

	return player, nil
}

// CreatePlayer создает нового игрока
func (s *PlayerService) CreatePlayer(player *models.Player) (*models.Player, error) {
	if player.Nickname == "" {
		return nil, fmt.Errorf("nickname is required")
	}

	// Проверяем, существует ли игрок с таким никнеймом
	existing, err := s.storage.GetPlayerByNickname(player.Nickname)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("player with nickname %s already exists", player.Nickname)
	}

	if err := s.storage.CreatePlayer(player); err != nil {
		s.logger.Error("Failed to create player", zap.Error(err))
		return nil, fmt.Errorf("failed to create player: %w", err)
	}

	s.logger.Info("Player created", zap.String("id", player.ID.String()), zap.String("nickname", player.Nickname))
	return player, nil
}

// UpdatePlayer обновляет данные игрока
func (s *PlayerService) UpdatePlayer(id string, updates *models.PlayerUpdateRequest) (*models.Player, error) {
	if id == "" {
		return nil, fmt.Errorf("player ID is required")
	}

	// Валидация обновлений
	if updates.Nickname != nil && *updates.Nickname == "" {
		return nil, fmt.Errorf("nickname cannot be empty")
	}

	player, err := s.storage.UpdatePlayer(id, updates)
	if err != nil {
		s.logger.Error("Failed to update player", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to update player: %w", err)
	}

	s.logger.Info("Player updated", zap.String("id", id))
	return player, nil
}

// DeletePlayer удаляет игрока
func (s *PlayerService) DeletePlayer(id string) error {
	if id == "" {
		return fmt.Errorf("player ID is required")
	}

	if err := s.storage.DeletePlayer(id); err != nil {
		s.logger.Error("Failed to delete player", zap.String("id", id), zap.Error(err))
		return fmt.Errorf("failed to delete player: %w", err)
	}

	s.logger.Info("Player deleted", zap.String("id", id))
	return nil
}

