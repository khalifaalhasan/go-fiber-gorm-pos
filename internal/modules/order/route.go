package order

import (
	"go-fiber-pos/internal/config"
	"go-fiber-pos/internal/middleware"
	"go-fiber-pos/internal/modules/inventory"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(adminGroup fiber.Router, db *gorm.DB, v *validator.Validate) {
	invRepo := inventory.NewInventoryRepository(db)
	invService := inventory.NewInventoryService(invRepo, v)

	repo := NewOrderRepository(db)
	service := NewOrderService(repo, invService, v)
	ctrl := NewOrderController(service)

	// Middleware khusus idempotency dari Redis dan DB untuk route Checkout
	idempotencyMiddleware := middleware.IdempotencyMiddleware(config.RedisClient, db)

	// Semua order endpoint membutuhkan autentikasi
	adminGroup.Post("/orders/checkout", idempotencyMiddleware, ctrl.Checkout)
	adminGroup.Get("/orders", ctrl.GetAll)
	adminGroup.Get("/orders/:id", ctrl.GetByID)
}
