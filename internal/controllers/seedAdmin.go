package controllers

import (
	"log"
	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {
	log.Println("========== SEED ADMIN STARTED ==========")

	var admin models.User

	result := config.DB.Where("email = ?", "admin@pharmacy.com").First(&admin)

	if result.Error != nil {
		log.Println("Find admin error:", result.Error)
	}

	if admin.ID == 0 {

		log.Println("Creating admin user...")

		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte("Admin123@"),
			bcrypt.DefaultCost,
		)

		if err != nil {
			log.Println("Hash error:", err)
			return
		}

		admin = models.User{
			Name:     "System Administrator",
			Email:    "admin@pharmacy.com",
			Password: string(hashedPassword),
			Role:     "admin",
		}

		result := config.DB.Create(&admin)

		if result.Error != nil {
			log.Println("ADMIN CREATE ERROR:", result.Error)
		} else {
			log.Println("Admin created successfully. ID:", admin.ID)
		}

	} else {
		log.Println("Admin already exists. ID:", admin.ID)
	}

	log.Println("========== SEED ADMIN FINISHED ==========")
}
