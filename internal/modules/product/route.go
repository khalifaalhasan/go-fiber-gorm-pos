package product

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes menerima Group dari terminal pusat (admin & public)
func SetupRoutes(adminGroup fiber.Router, publicGroup fiber.Router, db *gorm.DB) {
	// 1. Dependency Injection Khusus Product
	repo := NewProductRepository(db)
	service := NewProductService(repo)
	adminCtrl := NewProductController(service)
	publicCtrl := NewPublicProductController(service)

	// 2. Daftarkan Endpoint
	// Rute Admin (Otomatis kena middleware JWT dari pusat)
	adminGroup.Post("/products", adminCtrl.Create)

	// Rute Public (Katalog QR)
	publicGroup.Get("/menu/products", publicCtrl.GetAllMenu)
}