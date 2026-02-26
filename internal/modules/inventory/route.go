package inventory

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(router fiber.Router, db *gorm.DB, v *validator.Validate) {
	repo := NewInventoryRepository(db)
	service := NewInventoryService(repo, v)
	controller := NewInventoryController(service)

	// Protected Admin Routes (Sudah diamankan oleh JWT di api.go)
	admin := router.Group("/inventories")
	
	admin.Post("/adjust", controller.AdjustStock)
	admin.Get("/:productId", controller.GetStockByProductID)
	admin.Get("/:productId/movements", controller.GetMovements)
}
