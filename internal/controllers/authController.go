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

	// Query user
	result := config.DB.Where("email = ?", input.Email).First(&user)

	log.Println("========== LOGIN DEBUG ==========")
	log.Println("INPUT EMAIL:", input.Email)

	if result.Error != nil {
		log.Println("QUERY ERROR:", result.Error)
	}

	log.Println("ROWS FOUND:", result.RowsAffected)
	log.Println("USER ID:", user.ID)
	log.Println("DB EMAIL:", user.Email)
	log.Println("DB ROLE:", user.Role)

	// Show all users visible to GORM
	var users []models.User

	allUsersResult := config.DB.Find(&users)

	if allUsersResult.Error != nil {
		log.Println("FIND USERS ERROR:", allUsersResult.Error)
	}

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

		log.Println("PASSWORD CHECK FAILED")
		log.Println("BCRYPT ERROR:", err)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Password mismatch",
		})
		return
	}

	log.Println("LOGIN SUCCESSFUL")

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
func GetUsers(c *gin.Context) {
	var users []models.User

	config.DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}
