package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = os.Getenv("JWT_SECRET")

var jwtSecret = []byte(secret)

type Claims struct {
	AdminID uint   `json:"admin_id"`
	Role    string `json:"role"`

	jwt.RegisteredClaims
}

func GenerateToken(adminID uint, role string) (string, error) {

	claims := Claims{
		AdminID: adminID,
		Role:    role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(time.Hour * 24),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {

			return jwtSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
