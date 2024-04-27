package seeds

import (
	"Edot/models"
	"Edot/packages/postgres"

	"github.com/brianvoe/gofakeit/v7"
)

func Warehouse(db *postgres.DB, option string) error {
	var seed []*models.Warehouse

	for i := 0; i < 3; i++ {
		seed = append(seed, &models.Warehouse{
			Name:    gofakeit.Name(),
			Address: gofakeit.Address().Address,
			Status:  true,
		})
	}

	if option == "fresh" {
		if err := db.Gorm.Exec("ALTER TABLE warehouses DISABLE TRIGGER ALL").Error; err != nil {
			return err
		}
		if err := db.Gorm.Exec("DELETE FROM warehouses").Error; err != nil {
			return err
		}
		if err := db.Gorm.Exec("ALTER TABLE warehouses ENABLE TRIGGER ALL").Error; err != nil {
			return err
		}
		if err := db.Gorm.Exec("ALTER SEQUENCE warehouses_id_seq RESTART WITH 1").Error; err != nil {
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
