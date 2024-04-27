package models

import (
	"time"
)

type Cart struct {
	ID            int
	Name          string
	Price         float64
	Qty           int
	WarehouseID   int
	ToWarehouseID *int
	ProductID     int
	UserID        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
