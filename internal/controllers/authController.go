package controllers

import (
	"context"
	"log"
	"os"

	"net/http"
	"pharmacy-pos-backend/config"
	"pharmacy-pos-backend/internal/models"
	"pharmacy-pos-backend/internal/utils"

	"google.golang.org/api/idtoken"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
}

type GoogleLoginInput struct {
	Token string `json:"token" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func CreateUser(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	config.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered sucessfully",
	})
}

// func Register(c *gin.Context) {

// 	var input RegisterInput

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	// Check if user already exists
// 	var existingUser models.User

// 	config.DB.Where("email = ?", input.Email).First(&existingUser)

// 	if existingUser.ID != 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "User already exists",
// 		})
// 		return
// 	}

// 	hashedPassword, _ := bcrypt.GenerateFromPassword(
// 		[]byte(input.Password),
// 		bcrypt.DefaultCost,
// 	)

// 	user := models.User{
// 		Name:     input.Name,
// 		Email:    input.Email,
// 		Password: string(hashedPassword),

// 		// Force admin role for public registration
// 		Role: "admin",
// 	}

// 	config.DB.Create(&user)

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Admin registered successfully",
// 	})
// }

func Login(c *gin.Context) {

	var input LoginInput
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Find user
	config.DB.Where("email = ?", input.Email).First(&user)

	// Debug information
	println("========== LOGIN DEBUG ==========")
	println("INPUT EMAIL:", input.Email)
	println("INPUT PASSWORD:", input.Password)
	println("USER ID:", user.ID)
	println("DB EMAIL:", user.Email)
	println("DB ROLE:", user.Role)

	var users []models.User

	config.DB.Find(&users)

	log.Println("====== USERS IN DATABASE ======")

	for _, u := range users {
		log.Printf("ID=%d EMAIL=%s ROLE=%s\n", u.ID, u.Email, u.Role)
	}

	// User not found
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	// Compare password
	err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	)

	if err != nil {

		println("PASSWORD CHECK FAILED")
		println("HASH IN DB:", user.Password)
		println("BCRYPT ERROR:", err.Error())

		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Password mismatch",
			"details": err.Error(),
		})
		return
	}

	println("LOGIN SUCCESSFUL")

	token, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func GoogleLogin(c *gin.Context) {

	var input GoogleLoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	payload, err := idtoken.Validate(
		context.Background(),
		input.Token,
		os.Getenv("GOOGLE_CLIENT_ID"),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid Google token",
		})
		return
	}

	email := payload.Claims["email"].(string)
	name := payload.Claims["name"].(string)

	var user models.User

	config.DB.Where("email = ?", email).First(&user)

	// create user automatically if not found
	if user.ID == 0 {

		user = models.User{
			Name:  name,
			Email: email,
			Role:  "staff",
		}

		config.DB.Create(&user)
	}

	token, err := utils.GenerateToken(
		user.ID,
		user.Email,
		user.Role,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Google login successful",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}
