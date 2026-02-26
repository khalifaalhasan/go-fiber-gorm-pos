package payment

import (
	"errors"
	"fmt"
	"time"

	"go-fiber-pos/internal/core"
	"go-fiber-pos/internal/modules/inventory"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentService struct {
	repo       PaymentRepository
	gateway    PaymentGateway
	invService inventory.InventoryService
	v          *validator.Validate
}

func NewPaymentService(repo PaymentRepository, gateway PaymentGateway, invService inventory.InventoryService, v *validator.Validate) PaymentService {
	return &paymentService{repo: repo, gateway: gateway, invService: invService, v: v}
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
		fmt.Printf("ERROR calling gateway: %v\n", err)
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
		fmt.Printf("ERROR saving payment: %v\n", err)
		return nil, err
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
		// Pembayaran berhasil - Gunakan transaksi untuk update status & potong stok
		errTx := s.repo.ExecuteTx(func(tx *gorm.DB) error {
			now := time.Now()
			
			// Ambil order beserta itemnya untuk potong stok
			order, err := s.repo.FindOrderByID(p.OrderID)
			if err != nil {
				return err
			}

			// Lakukan pemotongan stok untuk setiap item
			for _, item := range order.Items {
				// Gunakan context background atau dari controller jika memungkinkan (di sini kita pakai context.Background)
				err := s.invService.DeductStockWithTx(tx.Statement.Context, tx, item.ProductID, item.Qty, "PAYMENT_SUCCESS", p.ID.String())
				if err != nil {
					// Jika gagal (misal out of stock), transaksi akan rollback.
					// Kita bisa kembalikan error spesifik agar dicatat.
					return fmt.Errorf("gagal memotong stok untuk produk %s: %w", item.ProductID, err)
				}
			}

			// Update status payment & order
			if err := tx.Model(&core.Payment{}).Where("id = ?", p.ID).
				Updates(map[string]interface{}{"payment_status": core.PaymentStatusPaid, "paid_at": &now, "webhook_received_at": now}).Error; err != nil {
				return err
			}

			if err := tx.Model(&core.Order{}).Where("id = ?", p.OrderID).
				Update("payment_status", core.PaymentStatusPaid).Error; err != nil {
				return err
			}

			return nil
		})

		if errTx != nil {
			// Jika gagal karena stok habis atau masalah DB
			fmt.Printf("Webhook Fulfillment Error: %v\n", errTx)
			// Idealnya di sini kita mencatat kegagalan pemenuhan pesanan (misalnya status PAID_BUT_UNFULFILLED)
			// Namun untuk sementara, kembalikan 500 agar webhook diretry atau error terekam.
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
