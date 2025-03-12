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

// Cargar variables de entorno desde .env
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env, usando valores predeterminados")
	}
}

// getJWTSecret obtiene la clave secreta desde .env
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("Advertencia: JWT_SECRET no está definido en .env, usando clave predeterminada")
		return []byte("clave_predeterminada")
	}
	return []byte(secret)
}

// JWTAuthMiddleware - Middleware para verificar tokens JWT
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Verificar formato del header
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "Token JWT requerido",
					"code":    "JWT_REQUIRED",
				},
			})
			c.Abort()
			return
		}

		// Extraer el token eliminando el prefijo "Bearer "
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parsear y validar el token JWT
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		})

		// Manejo de errores en el parsing del token
		if err != nil {
			var errorMsg string

			if errors.Is(err, jwt.ErrTokenExpired) {
				errorMsg = "El token ha expirado"
			} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				errorMsg = "Firma de token inválida"
			} else {
				errorMsg = "Token JWT inválido"
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

		// Extraer claims del token y almacenarlas en el contexto
		if token.Valid {
			username, _ := claims["username"].(string)
			c.Set("username", username)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"message": "Token JWT inválido",
					"code":    "JWT_INVALID",
				},
			})
			c.Abort()
		}
	}
}
