package handlers

import (
	"encoding/json"
	"net/http"

	"chrono-player-profile/internal/models"
	"chrono-player-profile/internal/service"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// ProfileUpdateHandler обрабатывает запросы на обновление профиля
type ProfileUpdateHandler struct {
	service *service.PlayerService
	logger  *zap.Logger
}

// NewProfileUpdateHandler создает новый обработчик обновления профиля
func NewProfileUpdateHandler(service *service.PlayerService, logger *zap.Logger) *ProfileUpdateHandler {
	return &ProfileUpdateHandler{
		service: service,
		logger:  logger,
	}
}

// Update обрабатывает PUT /api/v1/profile/{id}
func (h *ProfileUpdateHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "player ID is required")
		return
	}

	var updates models.PlayerUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	player, err := h.service.UpdatePlayer(id, &updates)
	if err != nil {
		h.logger.Error("Failed to update player", zap.String("id", id), zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, player)
}

// Create обрабатывает POST /api/v1/profile
func (h *ProfileUpdateHandler) Create(w http.ResponseWriter, r *http.Request) {
	var player models.Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	createdPlayer, err := h.service.CreatePlayer(&player)
	if err != nil {
		h.logger.Error("Failed to create player", zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusCreated, createdPlayer)
}

// Delete обрабатывает DELETE /api/v1/profile/{id}
func (h *ProfileUpdateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		h.respondError(w, http.StatusBadRequest, "player ID is required")
		return
	}

	if err := h.service.DeletePlayer(id); err != nil {
		h.logger.Error("Failed to delete player", zap.String("id", id), zap.Error(err))
		h.respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "player deleted successfully"})
}

// respondJSON отправляет JSON ответ
func (h *ProfileUpdateHandler) respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Error("Failed to encode JSON response", zap.Error(err))
	}
}

// respondError отправляет ошибку в JSON формате
func (h *ProfileUpdateHandler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

