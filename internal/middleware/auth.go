package middleware

import (
	pkg "go-fiber-pos/pkg/jwt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		
		// Fallback secret (menyamakan dengan yang di Auth Service)
		secretStr := os.Getenv("JWT_SECRET")
		if secretStr == "" {
			secretStr = "supersecretkey_ganti_nanti_di_env" 
		}
		secret := []byte(secretStr)

		token, err := jwt.ParseWithClaims(tokenString, &pkg.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
		}

		claims, ok := token.Claims.(*pkg.JwtCustomClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		// HANYA simpan user_id dan role. store_id resmi DIBUANG! 🚀
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}