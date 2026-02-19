package category

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes mendaftarkan endpoint khusus modul Category
func SetupRoutes(adminGroup fiber.Router, publicGroup fiber.Router, db *gorm.DB) {
	// 1. Dependency Injection Khusus Category
	repo := NewCategoryRepository(db)
	service := NewCategoryService(repo)
	adminCtrl := NewCategoryController(service)
	publicCtrl := NewPublicCategoryController(service)

	// 2. Daftarkan Endpoint
	
	// Rute Admin (Otomatis kena middleware JWT)
	adminGroup.Post("/categories", adminCtrl.Create)

	// Rute Public (Katalog Pelanggan / QR)
	publicGroup.Get("/menu/categories", publicCtrl.GetAllMenu)
}