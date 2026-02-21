package payment

import (
	"go-fiber-pos/internal/core"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *core.Payment) error {
	return r.db.Create(payment).Error
}

// FindByIdempotencyKey adalah kunci dari webhook idempotency.
// Mencari payment berdasarkan IdempotencyKey (= Midtrans order_id).
func (r *paymentRepository) FindByIdempotencyKey(key string) (*core.Payment, error) {
	var p core.Payment
	err := r.db.Where("idempotency_key = ?", key).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) FindByOrderID(orderID uuid.UUID) (*core.Payment, error) {
	var p core.Payment
	err := r.db.Where("order_id = ?", orderID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepository) FindOrderByID(orderID uuid.UUID) (*core.Order, error) {
	var order core.Order
	err := r.db.First(&order, "id = ?", orderID).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *paymentRepository) UpdateStatus(paymentID uuid.UUID, status string, paidAt *time.Time) error {
	updates := map[string]interface{}{"payment_status": status}
	if paidAt != nil {
		updates["paid_at"] = paidAt
	}
	return r.db.Model(&core.Payment{}).Where("id = ?", paymentID).Updates(updates).Error
}

func (r *paymentRepository) UpdateWebhookTimestamp(paymentID uuid.UUID, receivedAt time.Time) error {
	return r.db.Model(&core.Payment{}).
		Where("id = ?", paymentID).
		Update("webhook_received_at", receivedAt).Error
}

func (r *paymentRepository) UpdateOrderPaymentStatus(orderID uuid.UUID, status string) error {
	return r.db.Model(&core.Order{}).
		Where("id = ?", orderID).
		Update("payment_status", status).Error
}
