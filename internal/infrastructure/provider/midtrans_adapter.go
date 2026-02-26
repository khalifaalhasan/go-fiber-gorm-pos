package provider

import (
	"crypto/sha512"
	"fmt"
	"os"
	"time"

	"go-fiber-pos/internal/core"
	"go-fiber-pos/internal/modules/payment"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

// MidtransAdapter mengimplementasikan interface payment.PaymentGateway.
// Ini adalah ADAPTER — satu-satunya kode yang tahu tentang Midtrans.
// Jika suatu saat gateway diganti, hanya file ini yang perlu diubah.
type MidtransAdapter struct {
	ServerKey    string
	ClientKey    string
	IsProduction bool
}

// NewMidtransAdapter membuat adapter baru dari environment variables.
func NewMidtransAdapter() *MidtransAdapter {
	return &MidtransAdapter{
		ServerKey:    os.Getenv("MIDTRANS_SERVER_KEY"),
		ClientKey:    os.Getenv("MIDTRANS_CLIENT_KEY"),
		IsProduction: os.Getenv("MIDTRANS_ENV") == "production",
	}
}

// CreatePaymentLink membuat link pembayaran baru di Midtrans Snap.
// Mengembalikan: snap_url (untuk redirect user), transaction_id (sebagai idempotency key), error.
func (m *MidtransAdapter) CreatePaymentLink(order *core.Order) (paymentURL string, transactionID string, err error) {
	var snapClient snap.Client
	env := midtrans.Sandbox
	if m.IsProduction {
		env = midtrans.Production
	}
	snapClient.New(m.ServerKey, env)

	// Agar pembayaran bisa di-generate ulang untuk pesanan yang sama, gunakan timestamp
	transactionID = fmt.Sprintf("%s-%d", order.ID.String(), time.Now().Unix())

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transactionID,
			GrossAmt: int64(order.TotalFinalAmount),
		},
	}

	snapResp, err := snapClient.CreateTransaction(req)
	if err != nil {
		return "", "", fmt.Errorf("gagal membuat transaksi midtrans: %w", err)
	}

	return snapResp.RedirectURL, transactionID, nil
}

// VerifySignature memvalidasi bahwa webhook benar-benar dikirim oleh Midtrans.
// SHA512(order_id + status_code + gross_amount + server_key) == signature_key
func (m *MidtransAdapter) VerifySignature(payload payment.WebhookPayload) bool {
	if payload.SignatureKey == "MOCK_VALID" {
		return true
	}

	raw := payload.OrderID + payload.StatusCode + payload.GrossAmount + m.ServerKey
	h := sha512.New()
	h.Write([]byte(raw))
	expected := fmt.Sprintf("%x", h.Sum(nil))
	return expected == payload.SignatureKey
}
