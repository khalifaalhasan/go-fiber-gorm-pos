package order

import (
	"fmt"
	"time"

	"go-fiber-pos/internal/core"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

// DB mengekspos koneksi database untuk pembuatan transaksi di service layer.
func (r *orderRepository) ExecuteTx(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

// CreateWithTx menyimpan order beserta semua OrderItem dalam satu transaksi.
func (r *orderRepository) CreateWithTx(tx *gorm.DB, order *core.Order) error {
	return tx.Create(order).Error
}

// LockAndGetProduct mengambil produk dengan FOR UPDATE — row-level lock.
// Transaksi concurrent yang ingin mengakses baris yang sama AKAN MENUNGGU hingga lock dilepas.
func (r *orderRepository) LockAndGetProduct(tx *gorm.DB, productID uuid.UUID) (*core.Product, error) {
	var product core.Product
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&product, "id = ?", productID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetNextQueueNumber menghasilkan nomor antrean yang dijamin unik dan atomic.
// Menggunakan DailyCounter dengan FOR UPDATE untuk mencegah race condition.
func (r *orderRepository) GetNextQueueNumber(tx *gorm.DB, source string) (string, error) {
	prefix := "K"
	if source == core.OrderSourceEMenu {
		prefix = "E"
	}

	today := time.Now().Format("20060102") // "20260221"
	counterID := source + "-" + today      // "CASHIER-20260221"

	var counter core.DailyCounter

	// Kunci baris DailyCounter dengan FOR UPDATE dalam tx yang sama.
	// Jika baris tidak ada, gorm.ErrRecordNotFound akan dikembalikan.
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&counter, "id = ?", counterID).Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return "", err
		}
		// Baris belum ada untuk hari ini — buat baru dengan LastCount = 1
		counter = core.DailyCounter{
			ID:        counterID,
			Date:      today,
			Source:    source,
			LastCount: 1,
		}
		if createErr := tx.Create(&counter).Error; createErr != nil {
			return "", createErr
		}
	} else {
		// Baris sudah ada — increment dan simpan (masih dalam tx, masih dalam lock)
		counter.LastCount++
		if saveErr := tx.Save(&counter).Error; saveErr != nil {
			return "", saveErr
		}
	}

	return fmt.Sprintf("%s-%03d", prefix, counter.LastCount), nil
}

// FindVoucherByCode mencari voucher berdasarkan kode (baca biasa tanpa lock).
func (r *orderRepository) FindVoucherByCode(code string) (*core.Voucher, error) {
	var voucher core.Voucher
	err := r.db.Where("code = ?", code).First(&voucher).Error
	if err != nil {
		return nil, err
	}
	return &voucher, nil
}

// GetStoreMarkupFee mengambil markup fee dari profil toko, returns 0 jika belum dikonfigurasi.
func (r *orderRepository) GetStoreMarkupFee() int {
	var profile core.StoreProfile
	if err := r.db.First(&profile).Error; err != nil {
		return 0
	}
	return profile.MarkupFee
}

func (r *orderRepository) FindByID(id uuid.UUID) (*core.Order, error) {
	var order core.Order
	err := r.db.
		Preload("Items.Product").
		Preload("Voucher").
		Preload("Payments").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetAll() ([]core.Order, error) {
	var orders []core.Order
	err := r.db.
		Preload("Items.Product").
		Preload("Voucher").
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}
