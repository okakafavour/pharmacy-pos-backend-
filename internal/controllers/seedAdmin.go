package controllers

import (
	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {
	var admin models.User

	config.DB.Where("email = ?", "admin@pharmacy.com").First(&admin)

	if admin.ID == 0 {

		hashedPassword, _ := bcrypt.GenerateFromPassword(
			[]byte("Admin123@"),
			bcrypt.DefaultCost,
		)

		admin = models.User{
			Name:     "System Administrator",
			Email:    "admin@pharmacy.com",
			Password: string(hashedPassword),
			Role:     "admin",
		}

		config.DB.Create(&admin)
	}
}
