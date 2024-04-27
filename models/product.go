package models

import (
	"time"

	"gorm.io/gorm"
)

type (
	Product struct {
		ID        int
		Name      string
		Price     float64
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt gorm.DeletedAt

		//Relation
		Stocks []Stock `gorm:"<-:false;foreignKey:ProductID;references:ID;"`

		// Addional field
		Filter Filter `json:"-" gorm:"-"`
	}
)
