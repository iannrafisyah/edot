package models

import (
	"time"
)

type (
	StockOperator int

	Stock struct {
		ID            int
		Latest        int
		Previous      int
		Qty           int
		Operator      StockOperator
		WarehouseID   int
		ProductID     int
		TransactionID *int
		CreatedAt     time.Time
		UpdatedAt     time.Time

		// Relation
		Warehouse Warehouse `gorm:"<-:false;foreignKey:WarehouseID;references:ID;"`
	}
)

const (
	StockOperatorIncrement StockOperator = 0
	StockOperatorDecrement StockOperator = 1
)
