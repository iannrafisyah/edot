package models

import (
	"time"
)

type Warehouse struct {
	ID        int
	Name      string
	Address   string
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
