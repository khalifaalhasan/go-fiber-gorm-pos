package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


type JwtCustomClaims struct {
	UserID uuid.UUID `json:"user_id`
	StoreID uuid.UUID `json:"store_id`
	Role string `json:"role"`
	jwt.RegisteredClaims

}


func GenerateToken(userID, storeID uuid.UUID, role string)(string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	claims := &JwtCustomClaims{
		UserID : userID,
		StoreID: storeID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}