package seeds

import (
	"Edot/models"
	"Edot/packages/postgres"

	"github.com/brianvoe/gofakeit/v7"
)

func Stock(db *postgres.DB, option string) error {
	var (
		seed       []*models.Stock
		products   []*models.Product
		warehouses []*models.Warehouse
	)

	if err := db.Gorm.Find(&products).Error; err != nil {
		return err
	}

	if err := db.Gorm.Find(&warehouses).Error; err != nil {
		return err
	}

	for _, warehouse := range warehouses {
		for _, product := range products {
			totalQty := gofakeit.IntRange(100, 200)
			seed = append(seed, &models.Stock{
				Operator:    models.StockOperatorIncrement,
				Qty:         totalQty,
				Latest:      totalQty,
				ProductID:   product.ID,
				WarehouseID: warehouse.ID,
			})
		}
	}

	if option == "fresh" {
		if err := db.Gorm.Exec("ALTER TABLE stocks DISABLE TRIGGER ALL").Error; err != nil {
			return err
		}
		if err := db.Gorm.Exec("DELETE FROM stocks").Error; err != nil {
			return err
		}
		if err := db.Gorm.Exec("ALTER TABLE stocks ENABLE TRIGGER ALL").Error; err != nil {
			return err
		}
		if err := db.Gorm.Exec("ALTER SEQUENCE stocks_id_seq RESTART WITH 1").Error; err != nil {
			return err
		}
		if err := db.Gorm.Create(seed).Error; err != nil {
			return err
		}
	} else {
		for _, v := range seed {
			if err := db.Gorm.FirstOrCreate(v).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
