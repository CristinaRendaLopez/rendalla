package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// Initializes environment variables from a .env file.
// Logs a warning if the file is not found.
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file, using default values")
	}
}

// getJWTSecret retrieves the JWT secret from the environment.
// Falls back to a default value if none is defined.
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("Warning: JWT_SECRET is not defined in .env, using default secret")
		return []byte("default_secret")
	}
	return []byte(secret)
}

// JWTAuthMiddleware is a Gin middleware that enforces JWT authentication.
// It checks the Authorization header for a valid Bearer token, parses it,
// and extracts the username claim to make it available in the context.
//
// If the token is invalid, expired, or missing, it responds with 401 Unauthorized.
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
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

			c.JSON(http.StatusUnauthorized, gin.H{"error": errorMsg})
			c.Abort()
			return
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				c.Abort()
				return
			}
		}

		username, _ := claims["username"].(string)
		c.Set("username", username)
		c.Next()
	}
}
