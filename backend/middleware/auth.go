package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file, using default values")
	}
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("Warning: JWT_SECRET is not defined in .env, using default secret")
		return []byte("default_secret")
	}
	return []byte(secret)
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "JWT token required",
					"code":    "JWT_REQUIRED",
				},
			})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		})

		if err != nil {
			var errorMsg string

			if errors.Is(err, jwt.ErrTokenExpired) {
				errorMsg = "Token has expired"
			} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				errorMsg = "Invalid token signature"
			} else {
				errorMsg = "Invalid JWT token"
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": errorMsg,
					"code":    "JWT_INVALID",
				},
			})
			c.Abort()
			return
		}

		if token.Valid {
			username, _ := claims["username"].(string)
			c.Set("username", username)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "Invalid JWT token",
					"code":    "JWT_INVALID",
				},
			})
			c.Abort()
		}
	}
}
