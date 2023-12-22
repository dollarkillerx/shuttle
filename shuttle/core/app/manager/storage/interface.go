package storage

import (
	"gorm.io/gorm"
)

type Interface interface {
	DB() *gorm.DB
}
