package models

import (
	"APIs/internal/common/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Weather struct {
	UpdatedAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	DeletedAt gorm.DeletedAt

	Town string               `gorm:"size:256;not null"`
	Temp float32              `gorm:"not null"`
	Type entities.WeatherType `gorm:"not null"`

	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
}

func (m *Weather) TableName() string {
	return "weather"
}
