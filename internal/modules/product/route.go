package product

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func SetupRoutes(adminGroup fiber.Router, publicGroup fiber.Router, db *gorm.DB, v *validator.Validate) {
	
	repo := NewProductRepository(db)
	service := NewProductService(repo, v)
	adminCtrl := NewProductController(service)
	publicCtrl := NewPublicProductController(service)

	
	
	adminGroup.Post("/products", adminCtrl.Create)

	
	publicGroup.Get("/menu/products", publicCtrl.GetAllMenu)
}