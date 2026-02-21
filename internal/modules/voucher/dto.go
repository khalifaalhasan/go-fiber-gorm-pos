package voucher

import "time"

// CreateVoucherRequest adalah DTO untuk request pembuatan voucher baru.
type CreateVoucherRequest struct {
	Code              string    `json:"code" validate:"required,min=3,max=50"`
	DiscountType      string    `json:"discount_type" validate:"required,oneof=PERCENTAGE FIXED"`
	DiscountValue     int       `json:"discount_value" validate:"required,min=1"`
	MinOrderAmount    int       `json:"min_order_amount" validate:"min=0"`
	MaxDiscountAmount int       `json:"max_discount_amount" validate:"min=0"`
	ValidUntil        time.Time `json:"valid_until" validate:"required"`
}

// VoucherResponse adalah DTO untuk response data voucher.
type VoucherResponse struct {
	ID                string    `json:"id"`
	Code              string    `json:"code"`
	DiscountType      string    `json:"discount_type"`
	DiscountValue     int       `json:"discount_value"`
	MinOrderAmount    int       `json:"min_order_amount"`
	MaxDiscountAmount int       `json:"max_discount_amount"`
	ValidUntil        time.Time `json:"valid_until"`
	IsActive          bool      `json:"is_active"`
}
