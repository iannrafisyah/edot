package seeds

import (
	"Edot/models"
	"Edot/packages/postgres"

	"github.com/brianvoe/gofakeit/v7"
)

func Product(db *postgres.DB, option string) error {
	var seed []*models.Product

	for i := 0; i < 5; i++ {
		seed = append(seed, &models.Product{
			Name:  gofakeit.ProductName(),
			Price: gofakeit.Price(1000, 10000),
		})
	}

	if option == "fresh" {
		if err := db.Gorm.Exec("ALTER TABLE products DISABLE TRIGGER ALL").Error; err != nil {
			return err
		}
		if err := db.Gorm.Exec("DELETE FROM products").Error; err != nil {
			return err
		}
		if err := db.Gorm.Exec("ALTER TABLE products ENABLE TRIGGER ALL").Error; err != nil {
			return err
		}
		if err := db.Gorm.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1").Error; err != nil {
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
