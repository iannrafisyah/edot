package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type (
	TransactionProduct struct {
		ID            int
		Name          string
		Price         float64
		Qty           int
		ProductID     int
		WarehouseID   int
		ToWarehouseID *int
		TransactionID int
		Snapshot      TransactionProductSnapshot
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	TransactionProductSnapshot Product
)

func (r *TransactionProductSnapshot) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("could not parse transaction product snapshot to bytes")
	}

	result := TransactionProductSnapshot{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*r = result

	return nil
}

func (r TransactionProductSnapshot) Value() (driver.Value, error) {
	return json.Marshal(r)
}
