package models

import (
	"fmt"
	"time"

	"Edot/utilities"
)

type (
	TransactionStatus int
	TransactionType   int

	Transaction struct {
		ID         int
		Invoice    string
		GrandTotal float64
		Tax        float64
		Amount     float64
		UserID     int
		Status     TransactionStatus
		Type       TransactionType
		CreatedAt  time.Time
		UpdatedAt  time.Time

		//Relation
		TransactionProducts []TransactionProduct `gorm:"<-:false;foreignKey:TransactionID;references:ID;"`
	}
)

const (
	TransactionStatusUnpaid TransactionStatus = 0
	TransactionStatusPaid   TransactionStatus = 1
	TransactionStatusCancel TransactionStatus = 2

	TransactionTypeOrder         TransactionType = 0
	TransactionTypeTransferStock TransactionType = 1
)

func (r TransactionStatus) IsValid() error {
	switch r {
	case TransactionStatusUnpaid,
		TransactionStatusPaid,
		TransactionStatusCancel:
		return nil
	}

	return fmt.Errorf(utilities.DataNotFound, "transaction status")
}

func (r TransactionType) IsValid() error {
	switch r {
	case TransactionTypeOrder,
		TransactionTypeTransferStock:
		return nil
	}

	return fmt.Errorf(utilities.DataNotFound, "transaction type")
}

func (r TransactionStatus) String() string {
	switch r {
	case TransactionStatusUnpaid:
		return "Unpaid"
	case TransactionStatusPaid:
		return "Paid"
	case TransactionStatusCancel:
		return "Cancel"
	default:
		return "Unknown"
	}
}

func (r TransactionType) String() string {
	switch r {
	case TransactionTypeOrder:
		return "Order"
	case TransactionTypeTransferStock:
		return "Transfer Stock"
	default:
		return "Unknown"
	}
}
