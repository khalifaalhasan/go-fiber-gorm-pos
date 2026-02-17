package middleware

import (
	"go-fiber-pos/utils"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler{
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		secret := []byte(os.Getenv("JWT_SECRET"))

		token, err := jwt.ParseWithClaims(tokenString, &utils.JwtCustomClaims{}, func(token *jwt.Token)(interface{}, error){
			return secret, nil
		})

		if err != nil || !token.Valid{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error" : "Invalid or expired token"})


		}

		claims, ok := token.Claims.(*utils.JwtCustomClaims)
		if !ok{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error" : "Invalid token claims"})

		}

		c.Locals("user_id", claims.UserID)
		c.Locals("store_id", claims.StoreID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}