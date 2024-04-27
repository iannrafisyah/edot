package seeds

import (
	"Edot/models"
	"Edot/packages/postgres"

	"github.com/brianvoe/gofakeit/v7"
	"golang.org/x/crypto/bcrypt"
)

func User(db *postgres.DB, option string) error {
	var seed []*models.User

	password, err := bcrypt.GenerateFromPassword([]byte("testing123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	seed = append(seed, &models.User{
		FullName: gofakeit.Name(),
		Password: string(password),
		Email:    "user@mail.com",
	})

	if option == "fresh" {
		if err := db.Gorm.Exec("ALTER TABLE users DISABLE TRIGGER ALL").Error; err != nil {
			return err
		}
		if err := db.Gorm.Exec("DELETE FROM users").Error; err != nil {
			return err
		}
		if err := db.Gorm.Exec("ALTER TABLE users ENABLE TRIGGER ALL").Error; err != nil {
			return err
		}
		if err := db.Gorm.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1").Error; err != nil {
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
