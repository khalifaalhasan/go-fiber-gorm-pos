package order

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"go-fiber-pos/internal/core"
	"go-fiber-pos/internal/modules/inventory"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderService struct {
	repo       OrderRepository
	invService inventory.InventoryService
	v          *validator.Validate
}

func NewOrderService(repo OrderRepository, invService inventory.InventoryService, v *validator.Validate) OrderService {
	return &orderService{repo: repo, invService: invService, v: v}
}

func (s *orderService) Checkout(req CheckoutRequest) (*core.Order, error) {
	// 1. Validasi DTO
	if err := s.v.Struct(req); err != nil {
		return nil, err
	}

	var finalOrder *core.Order

	// ⭐ ANTI-DEADLOCK: Sort items berdasarkan ProductID ascending SEBELUM akuisisi lock.
	// Ini memastikan semua transaksi concurrent mengunci baris dalam urutan yang sama,
	// sehingga tidak ada circular wait → tidak ada deadlock.
	sort.Slice(req.Items, func(i, j int) bool {
		return strings.Compare(req.Items[i].ProductID.String(), req.Items[j].ProductID.String()) < 0
	})

	err := s.repo.ExecuteTx(func(tx *gorm.DB) error {
		var voucher *core.Voucher
		var voucherID *uuid.UUID

		if req.VoucherCode != "" {
			var err error
			voucher, err = s.repo.FindVoucherByCode(req.VoucherCode)
			if err != nil {
				return err
			}
			if !voucher.IsActive || voucher.ValidUntil.Before(time.Now()) {
				return core.ErrVoucherInvalid
			}
			voucherID = &voucher.ID
		}

		// 4. Hitung Platform Fee (Contoh: diambil dari settings toko)
		platformFee := s.repo.GetStoreMarkupFee()

		// 5. Generate Queue Number (Menggunakan FOR UPDATE di tabel daily_counters)
		queueNumber, err := s.repo.GetNextQueueNumber(tx, req.OrderSource)
		if err != nil {
			return err
		}

		// 6. Loop setiap item — akuisisi lock dan potong stok
		var orderItems []core.OrderItem
		var totalBasePrice int
		
		orderID := uuid.New().String()

		for _, item := range req.Items {
			// a. Kunci baris produk dengan FOR UPDATE
			product, err := s.repo.LockAndGetProduct(tx, item.ProductID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("%w: product not found", core.ErrNotFound)
				}
				return core.ErrInternalServer
			}

			// b. Kunci stok produk dan cek ketersediaan stok di modul inventory (tidak dipotong dulu)
			if err := s.invService.CheckStockWithTx(tx.Statement.Context, tx, product.ID, item.Qty); err != nil {
				if errors.Is(err, core.ErrInsufficientStock) {
					return fmt.Errorf("%w: %s", core.ErrInsufficientStock, product.Name)
				}
				return err
			}

			// c. Tentukan harga satuan (promo jika aktif dan dalam rentang waktu)
			unitPrice := calculateUnitPrice(product)
			subtotal := unitPrice * item.Qty
			totalBasePrice += subtotal

			// d. Siapkan entity OrderItem
			orderItems = append(orderItems, core.OrderItem{
				ID:        uuid.New(),
				ProductID: product.ID,
				Qty:       item.Qty,
				UnitPrice: unitPrice,
				Subtotal:  subtotal,
				Notes:     item.Notes, // Keep notes from original
			})
		}

		// 7. Hitung diskon voucher
		totalDiscount := 0
		if voucher != nil {
			if totalBasePrice < voucher.MinOrderAmount {
				return core.ErrVoucherMinOrder
			}
			totalDiscount = calculateDiscount(voucher, totalBasePrice)
		}

		// 8. Hitung Grand Total
		totalFinalAmount := totalBasePrice - totalDiscount + platformFee

		// 9. Buat entity Order dan simpan dalam transaksi
		orderUUID, _ := uuid.Parse(orderID)
		order := &core.Order{
			ID:               orderUUID,
			VoucherID:        voucherID,
			OrderSource:      req.OrderSource,
			QueueNumber:      queueNumber,
			TableNumber:      req.TableNumber, // Keep TableNumber from original
			OrderStatus:      core.OrderStatusPending, // Keep OrderStatus from original
			PaymentStatus:    core.PaymentStatusUnpaid, // Keep PaymentStatus from original
			TotalBasePrice:   totalBasePrice,
			TotalDiscount:    totalDiscount,
			PlatformFee:      platformFee,
			TotalFinalAmount: totalFinalAmount,
			Items:            orderItems,
		}

		if err := s.repo.CreateWithTx(tx, order); err != nil {
			return core.ErrInternalServer
		}

		finalOrder = order
		return nil
	})

	if err != nil {
		return nil, err
	}

	return finalOrder, nil
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
