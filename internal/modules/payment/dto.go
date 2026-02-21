package payment

import "github.com/google/uuid"

// InitiatePaymentRequest adalah DTO untuk request membuat link pembayaran.
type InitiatePaymentRequest struct {
	OrderID       uuid.UUID `json:"order_id" validate:"required"`
	PaymentMethod string    `json:"payment_method" validate:"required,oneof=CASH QRIS TRANSFER"`
}

// InitiatePaymentResponse berisi data payment yang berhasil dibuat.
type InitiatePaymentResponse struct {
	PaymentID      string `json:"payment_id"`
	PaymentURL     string `json:"payment_url"`     // URL redirect untuk QRIS/Transfer
	TransactionID  string `json:"transaction_id"`  // Midtrans transaction ID
	IdempotencyKey string `json:"idempotency_key"` // Untuk tracking webhook
	AmountDue      int    `json:"amount_due"`
	PaymentMethod  string `json:"payment_method"`
}

// WebhookPayload adalah DTO untuk menerima notifikasi dari Midtrans.
type WebhookPayload struct {
	OrderID           string `json:"order_id"`            // Ini adalah IdempotencyKey kita
	TransactionID     string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status"`  // settlement | capture | cancel | deny | expire | pending
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	FraudStatus       string `json:"fraud_status"`
}
