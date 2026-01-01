package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	ID          uuid.UUID `gorm:"type:char(36);not null;primaryKey"`
	Title       string    `gorm:"not null;uniqueIndex"`
	Description string    `gorm:"not null"`
	AuthorID    *uuid.UUID
	Author      *Author `gorm:"foreignKey:AuthorID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (b *Book) BeforeCreate(_ *gorm.DB) error {
	b.ID = uuid.New()
	return nil
}
