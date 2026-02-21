package payment

import (
	"go-fiber-pos/internal/core"

	"time"

	"github.com/google/uuid"
)

// PaymentGateway adalah PORT — abstraksi untuk semua payment gateway.
// Didefinisikan di sini (bukan di infrastructure) agar service layer
// hanya bergantung pada interface, bukan implementasi konkret.
type PaymentGateway interface {
	CreatePaymentLink(order *core.Order) (paymentURL string, transactionID string, err error)
	VerifySignature(payload WebhookPayload) bool
}

// PaymentRepository mendefinisikan kontrak akses data untuk Payment.
type PaymentRepository interface {
	Create(payment *core.Payment) error
	FindByIdempotencyKey(key string) (*core.Payment, error)
	FindByOrderID(orderID uuid.UUID) (*core.Payment, error)
	FindOrderByID(orderID uuid.UUID) (*core.Order, error)
	UpdateStatus(paymentID uuid.UUID, status string, paidAt *time.Time) error
	UpdateWebhookTimestamp(paymentID uuid.UUID, receivedAt time.Time) error
	UpdateOrderPaymentStatus(orderID uuid.UUID, status string) error
}

// PaymentService mendefinisikan kontrak business logic untuk Payment.
type PaymentService interface {
	InitiatePayment(req InitiatePaymentRequest) (*InitiatePaymentResponse, error)
	HandleWebhook(payload WebhookPayload) error
}
