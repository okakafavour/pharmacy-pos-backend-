package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireRole(roles ...string) gin.HandlerFunc {

	return func(c *gin.Context) {

		user, exists := c.Get("user")

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		claims := user.(jwt.MapClaims)

		userRole := claims["role"].(string)

		for _, role := range roles {

			if userRole == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
		})

		c.Abort()
	}
}
