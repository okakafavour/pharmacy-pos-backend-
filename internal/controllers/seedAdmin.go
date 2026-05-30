package controllers

import (
	"fmt"
	"log"
	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {
	log.Println("========== SEED ADMIN STARTED ==========")

	var admin models.User

	config.DB.Where("email = ?", "admin@pharmacy.com").First(&admin)

	if admin.ID == 0 {

		fmt.Println("Creating admin user...")

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

		fmt.Println("Admin created successfully")
	} else {
		fmt.Println("Admin already exists")
	}
	log.Println("========== SEED ADMIN FINISHED ==========")

}
