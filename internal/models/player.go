package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Player представляет игрока в системе
type Player struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Nickname  string    `gorm:"uniqueIndex;not null" json:"nickname"`
	Level     int       `gorm:"default:1" json:"level"`
	Rating    int       `gorm:"default:1000" json:"rating"`
	ELO       int       `gorm:"default:1000" json:"elo"`
	Role      string    `gorm:"type:varchar(50)" json:"role"`
	Region    string    `gorm:"type:varchar(50)" json:"region"`
	Language  string    `gorm:"type:varchar(10);default:'en'" json:"language"`
	
	// Статистика
	Wins      int `gorm:"default:0" json:"wins"`
	Losses    int `gorm:"default:0" json:"losses"`
	Rank      string `gorm:"type:varchar(50)" json:"rank"`
	
	// Косметика и настройки (храним как JSON)
	Cosmetics map[string]interface{} `gorm:"type:jsonb" json:"cosmetics"`
	Settings  map[string]interface{} `gorm:"type:jsonb" json:"settings"`
	
	// Предпочтения матчмейкинга
	PreferredMode string `gorm:"type:varchar(50)" json:"preferred_mode"`
	PreferredRole string `gorm:"type:varchar(50)" json:"preferred_role"`
	
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate hook для генерации UUID
func (p *Player) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// PlayerUpdateRequest представляет данные для обновления профиля
type PlayerUpdateRequest struct {
	Nickname      *string                `json:"nickname,omitempty"`
	Level         *int                   `json:"level,omitempty"`
	Rating        *int                   `json:"rating,omitempty"`
	ELO           *int                   `json:"elo,omitempty"`
	Role          *string                `json:"role,omitempty"`
	Region        *string                `json:"region,omitempty"`
	Language      *string                `json:"language,omitempty"`
	Wins          *int                   `json:"wins,omitempty"`
	Losses        *int                   `json:"losses,omitempty"`
	Rank          *string                `json:"rank,omitempty"`
	Cosmetics     map[string]interface{} `json:"cosmetics,omitempty"`
	Settings      map[string]interface{} `json:"settings,omitempty"`
	PreferredMode *string                `json:"preferred_mode,omitempty"`
	PreferredRole *string                `json:"preferred_role,omitempty"`
}

