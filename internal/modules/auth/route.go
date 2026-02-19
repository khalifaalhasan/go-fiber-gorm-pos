package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes mendaftarkan endpoint khusus modul Authentication
func SetupRoutes(route fiber.Router, db *gorm.DB) {
	// 1. Dependency Injection Khusus Auth
	repo := NewAuthRepository(db)
	service := NewAuthService(repo)
	ctrl := NewAuthController(service)

	// 2. Daftarkan Endpoint
	route.Post("/register", ctrl.Register)
	route.Post("/login", ctrl.Login)
}