package order

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-fiber-pos/internal/core"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderService struct {
	repo OrderRepository
	v    *validator.Validate
}

func NewOrderService(repo OrderRepository, v *validator.Validate) OrderService {
	return &orderService{repo: repo, v: v}
}

func (s *orderService) Checkout(req CheckoutRequest) (*core.Order, error) {
	// 1. Validasi DTO
	if err := s.v.Struct(req); err != nil {
		return nil, err
	}

	// 2. Buka transaksi database
	tx := s.repo.DB().Begin()
	if tx.Error != nil {
		return nil, core.ErrInternalServer
	}
	// Guard: jika terjadi panic, pastikan transaksi di-rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 3. Resolve Voucher (baca biasa, tanpa lock — tidak perlu)
	var voucher *core.Voucher
	var voucherID *uuid.UUID
	if req.VoucherCode != "" {
		v, err := s.repo.FindVoucherByCode(req.VoucherCode)
		if err != nil {
			tx.Rollback()
			return nil, core.ErrVoucherInvalid
		}
		if !v.IsActive || v.ValidUntil.Before(time.Now()) {
			tx.Rollback()
			return nil, core.ErrVoucherInvalid
		}
		voucher = v
		voucherID = &v.ID
	}

	// 4. ⭐ ANTI-DEADLOCK: Sort items berdasarkan ProductID ascending SEBELUM akuisisi lock.
	// Ini memastikan semua transaksi concurrent mengunci baris dalam urutan yang sama,
	// sehingga tidak ada circular wait → tidak ada deadlock.
	sort.Slice(req.Items, func(i, j int) bool {
		return strings.Compare(req.Items[i].ProductID.String(), req.Items[j].ProductID.String()) < 0
	})

	// 5. Generate nomor antrean menggunakan DailyCounter + FOR UPDATE (atomic)
	queueNumber, err := s.repo.GetNextQueueNumber(tx, req.OrderSource)
	if err != nil {
		tx.Rollback()
		return nil, core.ErrInternalServer
	}

	// 6. Loop setiap item — akuisisi lock dan potong stok
	var orderItems []core.OrderItem
	var totalBasePrice int

	for _, item := range req.Items {
		// a. Kunci baris produk dengan FOR UPDATE
		product, err := s.repo.LockAndGetProduct(tx, item.ProductID)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("produk dengan ID %s tidak ditemukan", item.ProductID)
			}
			return nil, core.ErrInternalServer
		}

		// b. Validasi stok SETELAH lock diperoleh (bukan sebelum!)
		if product.Stock < item.Qty {
			tx.Rollback()
			return nil, fmt.Errorf("%w: %s (tersisa %d)", core.ErrInsufficientStock, product.Name, product.Stock)
		}

		// c. Kurangi stok — masih dalam tx dan lock
		product.Stock -= item.Qty
		if err := s.repo.DeductStockWithTx(tx, product); err != nil {
			tx.Rollback()
			return nil, core.ErrInternalServer
		}

		// d. Tentukan harga satuan (promo jika aktif dan dalam rentang waktu)
		unitPrice := calculateUnitPrice(product)
		subtotal := unitPrice * item.Qty
		totalBasePrice += subtotal

		orderItems = append(orderItems, core.OrderItem{
			ID:        uuid.New(),
			ProductID: product.ID,
			Qty:       item.Qty,
			UnitPrice: unitPrice,
			Subtotal:  subtotal,
			Notes:     item.Notes,
		})
	}

	// 7. Hitung diskon voucher
	totalDiscount := 0
	if voucher != nil {
		if totalBasePrice < voucher.MinOrderAmount {
			tx.Rollback()
			return nil, core.ErrVoucherMinOrder
		}
		totalDiscount = calculateDiscount(voucher, totalBasePrice)
	}

	// 8. Ambil platform fee dari profil toko
	platformFee := s.repo.GetStoreMarkupFee()

	totalFinalAmount := totalBasePrice - totalDiscount + platformFee

	// 9. Buat entity Order dan simpan dalam transaksi
	order := &core.Order{
		ID:               uuid.New(),
		VoucherID:        voucherID,
		OrderSource:      req.OrderSource,
		QueueNumber:      queueNumber,
		TableNumber:      req.TableNumber,
		OrderStatus:      core.OrderStatusPending,
		PaymentStatus:    core.PaymentStatusUnpaid,
		TotalBasePrice:   totalBasePrice,
		TotalDiscount:    totalDiscount,
		PlatformFee:      platformFee,
		TotalFinalAmount: totalFinalAmount,
		Items:            orderItems,
	}

	if err := s.repo.CreateWithTx(tx, order); err != nil {
		tx.Rollback()
		return nil, core.ErrInternalServer
	}

	// 10. Commit — semua lock dilepas, semua perubahan permanen
	if err := tx.Commit().Error; err != nil {
		return nil, core.ErrInternalServer
	}

	return order, nil
}

func (s *orderService) GetAllOrders() ([]core.Order, error) {
	orders, err := s.repo.GetAll()
	if err != nil {
		return nil, core.ErrInternalServer
	}
	return orders, nil
}

func (s *orderService) GetOrderByID(id uuid.UUID) (*core.Order, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrNotFound
		}
		return nil, core.ErrInternalServer
	}
	return order, nil
}

// ===========================================
// HELPER FUNCTIONS (private)
// ===========================================

// calculateUnitPrice menentukan harga jual berdasarkan kondisi promo.
func calculateUnitPrice(product *core.Product) int {
	if !product.IsPromoActive || product.PromoPrice <= 0 {
		return product.NormalPrice
	}

	// Cek apakah saat ini berada dalam rentang waktu promo
	now := time.Now()
	currentTime := fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())

	if product.PromoStartTime != "" && product.PromoEndTime != "" {
		if currentTime >= product.PromoStartTime && currentTime <= product.PromoEndTime {
			return product.PromoPrice
		}
	}

	return product.NormalPrice
}

// calculateDiscount menghitung jumlah diskon berdasarkan tipe voucher.
func calculateDiscount(voucher *core.Voucher, totalBase int) int {
	switch voucher.DiscountType {
	case core.DiscountTypePercentage:
		discount := totalBase * voucher.DiscountValue / 100
		// Terapkan MaxDiscountAmount jika ada batasan
		if voucher.MaxDiscountAmount > 0 && discount > voucher.MaxDiscountAmount {
			return voucher.MaxDiscountAmount
		}
		return discount
	case core.DiscountTypeFixed:
		if voucher.DiscountValue > totalBase {
			return totalBase // Diskon tidak boleh melebihi total belanja
		}
		return voucher.DiscountValue
	default:
		return 0
	}
}
