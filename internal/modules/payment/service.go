package payment

import (
	"errors"
	"time"

	"go-fiber-pos/internal/core"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentService struct {
	repo    PaymentRepository
	gateway PaymentGateway
	v       *validator.Validate
}

func NewPaymentService(repo PaymentRepository, gateway PaymentGateway, v *validator.Validate) PaymentService {
	return &paymentService{repo: repo, gateway: gateway, v: v}
}

// InitiatePayment membuat payment record baru dan link pembayaran via gateway.
func (s *paymentService) InitiatePayment(req InitiatePaymentRequest) (*InitiatePaymentResponse, error) {
	if err := s.v.Struct(req); err != nil {
		return nil, err
	}

	// Ambil order
	order, err := s.repo.FindOrderByID(req.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrNotFound
		}
		return nil, core.ErrInternalServer
	}

	// Cek apakah order sudah lunas
	if order.PaymentStatus == core.PaymentStatusPaid {
		return nil, core.ErrOrderAlreadyPaid
	}

	// Buat link pembayaran via gateway (PORT & ADAPTER)
	paymentURL, transactionID, err := s.gateway.CreatePaymentLink(order)
	if err != nil {
		return nil, core.ErrInternalServer
	}

	// Simpan payment record dengan IdempotencyKey = transactionID dari gateway
	p := &core.Payment{
		ID:                    uuid.New(),
		OrderID:               order.ID,
		PaymentMethod:         req.PaymentMethod,
		MidtransTransactionID: &transactionID,
		IdempotencyKey:        transactionID, // Ini yang di-cek saat webhook masuk
		AmountPaid:            order.TotalFinalAmount,
		PaymentStatus:         core.PaymentStatusUnpaid,
	}

	if err := s.repo.Create(p); err != nil {
		return nil, core.ErrInternalServer
	}

	return &InitiatePaymentResponse{
		PaymentID:      p.ID.String(),
		PaymentURL:     paymentURL,
		TransactionID:  transactionID,
		IdempotencyKey: p.IdempotencyKey,
		AmountDue:      order.TotalFinalAmount,
		PaymentMethod:  req.PaymentMethod,
	}, nil
}

// HandleWebhook memproses notifikasi pembayaran dari Midtrans dengan pengecekan idempotency.
func (s *paymentService) HandleWebhook(payload WebhookPayload) error {
	// 1. Verifikasi bahwa webhook benar-benar dari Midtrans
	if !s.gateway.VerifySignature(payload) {
		return core.ErrInvalidSignature
	}

	// 2. Cari payment berdasarkan idempotency key (= Midtrans order_id)
	p, err := s.repo.FindByIdempotencyKey(payload.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.ErrNotFound
		}
		return core.ErrInternalServer
	}

	// 3. ⭐ IDEMPOTENCY CHECK — Inti keamanan webhook
	// Jika sudah PAID, abaikan webhook duplikat. Return nil agar Midtrans tidak retry.
	if p.PaymentStatus == core.PaymentStatusPaid {
		return nil
	}

	// 4. Proses berdasarkan status dari Midtrans
	switch payload.TransactionStatus {
	case "settlement", "capture":
		// Pembayaran berhasil
		now := time.Now()
		if err := s.repo.UpdateStatus(p.ID, core.PaymentStatusPaid, &now); err != nil {
			return core.ErrInternalServer
		}
		if err := s.repo.UpdateWebhookTimestamp(p.ID, now); err != nil {
			return core.ErrInternalServer
		}
		// Update juga payment_status di tabel orders
		if err := s.repo.UpdateOrderPaymentStatus(p.OrderID, core.PaymentStatusPaid); err != nil {
			return core.ErrInternalServer
		}

	case "cancel", "deny", "expire":
		// Pembayaran gagal
		if err := s.repo.UpdateStatus(p.ID, core.PaymentStatusFailed, nil); err != nil {
			return core.ErrInternalServer
		}

	// "pending" dan status lain tidak perlu tindakan
	}

	return nil
}
