package provider

import (
	"crypto/sha512"
	"fmt"
	"os"

	"go-fiber-pos/internal/core"
	"go-fiber-pos/internal/modules/payment"
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
//
// Untuk production, integrasikan dengan SDK: github.com/midtrans/midtrans-go
// Implementasi di bawah adalah mock yang sudah menerapkan interface yang benar.
func (m *MidtransAdapter) CreatePaymentLink(order *core.Order) (paymentURL string, transactionID string, err error) {
	// --- MOCK IMPLEMENTATION ---
	// Dalam production, gantikan blok ini dengan pemanggilan Midtrans Snap API:
	//
	//   snapClient := snap.Client{}
	//   snapClient.New(m.ServerKey, midtrans.Sandbox)
	//   req := &snap.Request{
	//       TransactionDetails: midtrans.TransactionDetails{
	//           OrderID:  order.ID.String(),
	//           GrossAmt: int64(order.TotalFinalAmount),
	//       },
	//   }
	//   snapResp, err := snapClient.CreateTransaction(req)
	//   return snapResp.RedirectURL, snapResp.Token, err

	mockTransactionID := fmt.Sprintf("MOCK-MIDTRANS-%s", order.ID.String())
	mockURL := fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v2/vtweb/%s", mockTransactionID)
	return mockURL, mockTransactionID, nil
}

// VerifySignature memvalidasi bahwa webhook benar-benar dikirim oleh Midtrans.
// SHA512(order_id + status_code + gross_amount + server_key) == signature_key
func (m *MidtransAdapter) VerifySignature(payload payment.WebhookPayload) bool {
	// --- MOCK IMPLEMENTATION ---
	// Dalam production, verifikasi dengan server key sungguhan:
	//
	//   raw := payload.OrderID + payload.StatusCode + payload.GrossAmount + m.ServerKey
	//   h := sha512.New()
	//   h.Write([]byte(raw))
	//   expected := fmt.Sprintf("%x", h.Sum(nil))
	//   return expected == payload.SignatureKey

	// Untuk development/mock: jika SignatureKey == "MOCK_VALID", anggap valid
	if payload.SignatureKey == "MOCK_VALID" {
		return true
	}

	// Tetap implementasikan algoritma verifikasi yang benar untuk kemudahan switch ke production
	raw := payload.OrderID + payload.StatusCode + payload.GrossAmount + m.ServerKey
	h := sha512.New()
	h.Write([]byte(raw))
	expected := fmt.Sprintf("%x", h.Sum(nil))
	return expected == payload.SignatureKey
}
