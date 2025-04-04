package models

import (
	"APIs/internal/common/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Promocode struct {
	UpdatedAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	DeletedAt gorm.DeletedAt

	Name             string     `gorm:"size:256;not null"`
	AdvantagePercent int64      `gorm:"not null"`
	Restrictions     utils.JSON `gorm:"not null"`

	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
}

func (m *Promocode) TableName() string {
	return "promocode"
}
