package utils

import (
	"github.com/golang-jwt/jwt/v5"
)

// TokenGenerator defines an interface for generating signed JWT tokens
// from a given set of claims.
type TokenGenerator interface {
	GenerateToken(claims jwt.MapClaims) (string, error)
}

// JWTTokenGenerator implements TokenGenerator using HMAC with a provided secret.
type JWTTokenGenerator struct {
	Secret []byte
}

// GenerateToken creates a signed JWT token with the given claims using HS256.
// Returns:
//   - the signed token as a string on success
//   - error if signing fails
func (g *JWTTokenGenerator) GenerateToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(g.Secret)
}

// Ensure JWTTokenGenerator satisfies the TokenGenerator interface.
var _ TokenGenerator = (*JWTTokenGenerator)(nil)
