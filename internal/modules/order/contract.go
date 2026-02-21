package order

import (
	"go-fiber-pos/internal/core"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrderRepository mendefinisikan kontrak akses data untuk Order.
type OrderRepository interface {
	// CreateWithTx menyimpan order dan semua item-nya dalam satu transaksi database.
	CreateWithTx(tx *gorm.DB, order *core.Order) error
	// LockAndGetProduct mengambil product dengan FOR UPDATE pessimistic lock untuk mencegah race condition stok.
	LockAndGetProduct(tx *gorm.DB, productID uuid.UUID) (*core.Product, error)
	// DeductStockWithTx memperbarui stok produk dalam transaksi yang sudah ada.
	DeductStockWithTx(tx *gorm.DB, product *core.Product) error
	// GetNextQueueNumber menggunakan DailyCounter + FOR UPDATE untuk generate nomor antrean atomic.
	GetNextQueueNumber(tx *gorm.DB, source string) (string, error)
	// FindByCode mencari voucher berdasarkan kode (tanpa lock, baca biasa).
	FindVoucherByCode(code string) (*core.Voucher, error)
	// GetStoreMarkupFee mengambil markup fee dari profil toko.
	GetStoreMarkupFee() int
	FindByID(id uuid.UUID) (*core.Order, error)
	GetAll() ([]core.Order, error)
	DB() *gorm.DB
}

// OrderService mendefinisikan kontrak business logic untuk Order.
type OrderService interface {
	Checkout(req CheckoutRequest) (*core.Order, error)
	GetAllOrders() ([]core.Order, error)
	GetOrderByID(id uuid.UUID) (*core.Order, error)
}
