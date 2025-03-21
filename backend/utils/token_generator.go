package utils

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenGenerator struct {
	Secret []byte
}

func (g *JWTTokenGenerator) GenerateToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(g.Secret)
}
