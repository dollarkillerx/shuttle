package models

import (
	"gorm.io/gorm"

	"time"
)

type BasicModel struct {
	ID        string         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"index" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
