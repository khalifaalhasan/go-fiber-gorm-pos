package order

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(adminGroup fiber.Router, db *gorm.DB, v *validator.Validate) {
	repo := NewOrderRepository(db)
	service := NewOrderService(repo, v)
	ctrl := NewOrderController(service)

	// Semua order endpoint membutuhkan autentikasi
	adminGroup.Post("/orders/checkout", ctrl.Checkout)
	adminGroup.Get("/orders", ctrl.GetAll)
	adminGroup.Get("/orders/:id", ctrl.GetByID)
}
