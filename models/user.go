package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int
	FullName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
