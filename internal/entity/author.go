package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Author struct {
	ID        uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Name      string    `gorm:"not null;unique"`
	Book      []Book    `gorm:"constraint:OnDelete:SET NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *Author) BeforeCreate(_ *gorm.DB) error {
	a.ID = uuid.New()
	return nil
}
