package voucher

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(adminGroup fiber.Router, db *gorm.DB, v *validator.Validate) {
	repo := NewVoucherRepository(db)
	service := NewVoucherService(repo, v)
	ctrl := NewVoucherController(service)

	// Admin-only endpoints
	adminGroup.Post("/vouchers", ctrl.Create)
	adminGroup.Get("/vouchers", ctrl.GetAll)
	adminGroup.Delete("/vouchers/:id", ctrl.Delete)
}
