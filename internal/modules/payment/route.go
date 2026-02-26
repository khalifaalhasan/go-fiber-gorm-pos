package payment

import (
	"go-fiber-pos/internal/modules/inventory"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(adminGroup fiber.Router, webhookGroup fiber.Router, db *gorm.DB, v *validator.Validate, gateway PaymentGateway, invService inventory.InventoryService) {
	repo := NewPaymentRepository(db)
	service := NewPaymentService(repo, gateway, invService, v)
	ctrl := NewPaymentController(service)

	// Admin: membuat link pembayaran
	adminGroup.Post("/payments/initiate", ctrl.InitiatePayment)

	// Webhook: public endpoint (tanpa JWT) — Midtrans mengirim notifikasi ke sini
	webhookGroup.Post("/payment", ctrl.HandleWebhook)
}
