package handlers

import (
	"encoding/json"
	"net/http"

	"chrono-player-profile/internal/service"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// ProfileGetHandler обрабатывает запросы на получение профиля
type ProfileGetHandler struct {
	service *service.PlayerService
	logger  *zap.Logger
}

// NewProfileGetHandler создает новый обработчик получения профиля
func NewProfileGetHandler(service *service.PlayerService, logger *zap.Logger) *ProfileGetHandler {
	return &ProfileGetHandler{
		service: service,
		logger:  logger,
	}
}

// GetByID обрабатывает GET /api/v1/profile/{id}
func (h *ProfileGetHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "player ID is required")
		return
	}

	player, err := h.service.GetPlayerByID(id)
	if err != nil {
		h.logger.Error("Failed to get player", zap.String("id", id), zap.Error(err))
		h.respondError(w, http.StatusNotFound, "player not found")
		return
	}

	h.respondJSON(w, http.StatusOK, player)
}

// GetByNickname обрабатывает GET /api/v1/profile/nickname/{nickname}
func (h *ProfileGetHandler) GetByNickname(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]

	if nickname == "" {
		h.respondError(w, http.StatusBadRequest, "nickname is required")
		return
	}

	player, err := h.service.GetPlayerByNickname(nickname)
	if err != nil {
		h.logger.Error("Failed to get player by nickname", zap.String("nickname", nickname), zap.Error(err))
		h.respondError(w, http.StatusNotFound, "player not found")
		return
	}

	h.respondJSON(w, http.StatusOK, player)
}

// respondJSON отправляет JSON ответ
func (h *ProfileGetHandler) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Error("Failed to encode JSON response", zap.Error(err))
	}
}

// respondError отправляет ошибку в JSON формате
func (h *ProfileGetHandler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

