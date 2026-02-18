package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JwtCustomClaims: store_id RESMI DIHAPUS. 
// (Typo kurang tanda kutip tutup di json tag juga sudah diperbaiki)
type JwtCustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken: Parameter storeID dihapus
func GenerateToken(userID uuid.UUID, role string) (string, error) {
	secretStr := os.Getenv("JWT_SECRET")
	if secretStr == "" {
		// Fallback biar nggak panic kalau lupa set .env
		secretStr = "supersecretkey_ganti_nanti_di_env" 
	}
	secret := []byte(secretStr)

	claims := &JwtCustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}