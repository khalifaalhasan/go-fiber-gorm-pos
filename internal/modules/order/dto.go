package order

import "github.com/google/uuid"

// CheckoutItemInput adalah DTO untuk satu item dalam request checkout.
type CheckoutItemInput struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Qty       int       `json:"qty" validate:"required,min=1"`
	Notes     string    `json:"notes"`
}

// CheckoutRequest adalah DTO untuk request checkout order baru.
type CheckoutRequest struct {
	OrderSource string              `json:"order_source" validate:"required,oneof=CASHIER E_MENU"`
	TableNumber *string             `json:"table_number"`
	VoucherCode string              `json:"voucher_code"`
	Items       []CheckoutItemInput `json:"items" validate:"required,min=1,dive"`
}

type CheckoutResponse struct {
	OrderID     uuid.UUID `json:"order_id"`
	OrderNumber string    `json:"order_number"`
	OrderSource string    `json:"order_source"`
	TableNumber *string   `json:"table_number"`
	VoucherCode string    `json:"voucher_code"`
	Items       []CheckoutItemInput `json:"items" validate:"required,min=1,dive"`
}
