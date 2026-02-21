package payment

import (
	"errors"

	"go-fiber-pos/internal/core"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PaymentController struct {
	service PaymentService
}

func NewPaymentController(service PaymentService) *PaymentController {
	return &PaymentController{service: service}
}

// InitiatePayment membuat link pembayaran untuk sebuah order.
// Endpoint: POST /admin/payments/initiate
func (ctrl *PaymentController) InitiatePayment(c *fiber.Ctx) error {
	var req InitiatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	resp, err := ctrl.service.InitiatePayment(req)
	if err != nil {
		var valErr validator.ValidationErrors
		if errors.As(err, &valErr) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Validasi gagal", "details": valErr.Error()})
		}
		if errors.Is(err, core.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order tidak ditemukan"})
		}
		if errors.Is(err, core.ErrOrderAlreadyPaid) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Link pembayaran berhasil dibuat",
		"data":    resp,
	})
}

// HandleWebhook menerima notifikasi pembayaran dari Midtrans.
// Endpoint: POST /webhook/payment (PUBLIC — tidak butuh JWT)
// Selalu return 200 agar Midtrans tidak melakukan retry yang tidak perlu.
func (ctrl *PaymentController) HandleWebhook(c *fiber.Ctx) error {
	var payload WebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		// Return 200 meskipun format invalid untuk menghentikan retry Midtrans
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "invalid payload format"})
	}

	err := ctrl.service.HandleWebhook(payload)
	if err != nil {
		if errors.Is(err, core.ErrInvalidSignature) {
			// Ini mungkin serangan — kembalikan 403 bukan 200
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Signature tidak valid"})
		}
		// Error internal jangan di-200 — biarkan Midtrans retry
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Selalu return 200 untuk notifikasi yang valid (termasuk yang sudah diproses/idempotent)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
}
