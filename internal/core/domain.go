package core

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

// ==========================================
// CONSTANTS / ENUMS
// ==========================================

const (
	// Order Source
	OrderSourceCashier = "CASHIER"
	OrderSourceEMenu   = "E_MENU"

	// Order Status
	OrderStatusPending   = "PENDING"
	OrderStatusCompleted = "COMPLETED"
	OrderStatusCancelled = "CANCELLED"

	// Payment Status
	PaymentStatusUnpaid = "UNPAID"
	PaymentStatusPaid   = "PAID"
	PaymentStatusFailed = "FAILED"

	// Payment Method
	PaymentMethodCash     = "CASH"
	PaymentMethodQRIS     = "QRIS"
	PaymentMethodTransfer = "TRANSFER"

	// Voucher Discount Type
	DiscountTypePercentage = "PERCENTAGE"
	DiscountTypeFixed      = "FIXED"
)

// ==========================================
// APP CONFIGURATION (Single Row)
// ==========================================

// StoreProfile menyimpan konfigurasi toko. Single-tenant, hanya 1 baris.
type StoreProfile struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Address   string    `gorm:"type:text" json:"address"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	MarkupFee int       `gorm:"default:0" json:"markup_fee"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ==========================================
// USERS
// ==========================================

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Username     string         `gorm:"type:varchar(255);unique;not null" json:"username"`
	PasswordHash string         `gorm:"type:varchar(255);not null" json:"-"`
	Role         string         `gorm:"type:varchar(50);default:'CASHIER'" json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// ==========================================
// MASTER DATA (CATALOG)
// ==========================================

type Category struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);unique;not null" json:"name"`
	Slug      string         `gorm:"type:varchar(255);unique;index" json:"slug"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Product struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CategoryID     uuid.UUID      `gorm:"type:uuid;not null" json:"category_id"`
	Name           string         `gorm:"type:varchar(255);not null" json:"name"`
	Slug           string         `gorm:"type:varchar(255);index" json:"slug"`
	Description    string         `gorm:"type:text" json:"description"`
	ImageURL       string         `gorm:"type:varchar(255)" json:"image_url"`
	NormalPrice    int            `gorm:"not null" json:"normal_price"`
	IsAvailable    bool           `gorm:"default:true" json:"is_available"`
	IsPromoActive  bool           `gorm:"default:false" json:"is_promo_active"`
	PromoPrice     int            `json:"promo_price"`
	PromoStartTime string         `gorm:"type:varchar(5)" json:"promo_start_time"` // Format "HH:MM"
	PromoEndTime   string         `gorm:"type:varchar(5)" json:"promo_end_time"`   // Format "HH:MM"
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Category  *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Inventory *Inventory `gorm:"foreignKey:ProductID" json:"inventory,omitempty"`
}

// ==========================================
// INVENTORY
// ==========================================

type Inventory struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProductID    uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_product_inventory;not null" json:"product_id"`
	QtyAvailable int       `gorm:"not null;default:0" json:"qty_available"`
	QtyReserved  int       `gorm:"not null;default:0" json:"qty_reserved"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

type InventoryMovement struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	InventoryID   uuid.UUID `gorm:"type:uuid;not null;index" json:"inventory_id"`
	ReferenceType string    `gorm:"type:varchar(50);not null" json:"reference_type"` // ORDER, RESTOCK, ADJUSTMENT
	ReferenceID   string    `gorm:"type:varchar(255)" json:"reference_id"`
	QtyChange     int       `gorm:"not null" json:"qty_change"`
	Notes         string    `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time `json:"created_at"`

	Inventory *Inventory `gorm:"foreignKey:InventoryID" json:"inventory,omitempty"`
}

// ==========================================
// VOUCHERS
// ==========================================

type Voucher struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Code              string    `gorm:"type:varchar(50);unique;not null" json:"code"`
	DiscountType      string    `gorm:"type:varchar(50);not null" json:"discount_type"` // PERCENTAGE | FIXED
	DiscountValue     int       `gorm:"not null" json:"discount_value"`
	MinOrderAmount    int       `gorm:"default:0" json:"min_order_amount"`
	MaxDiscountAmount int       `gorm:"default:0" json:"max_discount_amount"` // 0 = no cap (untuk FIXED tidak relevan)
	ValidUntil        time.Time `json:"valid_until"`
	IsActive          bool      `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
}

// ==========================================
// DAILY COUNTER (Untuk atomic queue number)
// ==========================================

// DailyCounter adalah Single Source of Truth untuk nomor antrean per hari per source.
// Menggunakan FOR UPDATE pessimistic lock saat increment untuk mencegah race condition.
type DailyCounter struct {
	ID        string `gorm:"type:varchar(50);primaryKey" json:"id"`     // Format: "CASHIER-20260221"
	Date      string `gorm:"type:varchar(10);not null" json:"date"`      // Format: "20260221"
	Source    string `gorm:"type:varchar(50);not null" json:"source"`    // CASHIER | E_MENU
	LastCount int    `gorm:"not null;default:0" json:"last_count"`
}

// ==========================================
// TRANSACTIONAL
// ==========================================

type Order struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	VoucherID        *uuid.UUID `gorm:"type:uuid" json:"voucher_id"` // Pointer karena opsional
	OrderSource      string     `gorm:"type:varchar(50);not null;default:'CASHIER'" json:"order_source"` // CASHIER | E_MENU
	QueueNumber      string     `gorm:"type:varchar(20)" json:"queue_number"`                            // K-001 | E-001
	TableNumber      *string    `gorm:"type:varchar(50)" json:"table_number"`
	OrderStatus      string     `gorm:"type:varchar(50);default:'PENDING'" json:"order_status"`
	PaymentStatus    string     `gorm:"type:varchar(50);default:'UNPAID'" json:"payment_status"`
	TotalBasePrice   int        `gorm:"not null" json:"total_base_price"`
	TotalDiscount    int        `gorm:"default:0" json:"total_discount"`
	PlatformFee      int        `gorm:"default:0" json:"platform_fee"`
	TotalFinalAmount int        `gorm:"not null" json:"total_final_amount"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	Voucher  *Voucher    `gorm:"foreignKey:VoucherID" json:"voucher,omitempty"`
	Items    []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	Payments []Payment   `gorm:"foreignKey:OrderID" json:"payments"`
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Qty       int       `gorm:"not null" json:"qty"`
	UnitPrice int       `gorm:"not null" json:"unit_price"`
	Subtotal  int       `gorm:"not null" json:"subtotal"`
	Notes     string    `gorm:"type:varchar(255)" json:"notes"`
	CreatedAt time.Time `json:"created_at"`

	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

type Payment struct {
	ID                   uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrderID              uuid.UUID  `gorm:"type:uuid;not null" json:"order_id"`
	PaymentMethod        string     `gorm:"type:varchar(50);not null" json:"payment_method"` // CASH | QRIS | TRANSFER
	MidtransTransactionID *string   `gorm:"type:varchar(255)" json:"midtrans_transaction_id"`
	IdempotencyKey       string     `gorm:"type:varchar(255);uniqueIndex" json:"idempotency_key"` // Midtrans order_id, mencegah duplikasi webhook
	AmountPaid           int        `gorm:"not null" json:"amount_paid"`
	PaymentStatus        string     `gorm:"type:varchar(50);not null;default:'UNPAID'" json:"payment_status"` // UNPAID | PAID | FAILED
	PaidAt               *time.Time `gorm:"type:timestamptz" json:"paid_at"`
	WebhookReceivedAt    *time.Time `gorm:"type:timestamptz" json:"webhook_received_at"` // Timestamp saat webhook diterima pertama kali
	CreatedAt            time.Time  `json:"created_at"`
}

// ==========================================
// IDEMPOTENCY
// ==========================================

// IdempotencyRecord adalah source of truth untuk request yang sudah berhasil diproses.
type IdempotencyRecord struct {
	Key          string    `gorm:"type:varchar(255);primaryKey" json:"key"`
	ResponseBody string    `gorm:"type:text;not null" json:"response_body"`
	StatusCode   int       `gorm:"not null" json:"status_code"`
	CreatedAt    time.Time `json:"created_at"`
}

// ==========================================
// GORM HOOKS
// ==========================================

func (c *Category) BeforeSave(tx *gorm.DB) (err error) {
	c.Slug = slug.Make(c.Name)
	return
}

func (p *Product) BeforeSave(tx *gorm.DB) (err error) {
	p.Slug = slug.Make(p.Name)
	return
}

func (p *Product) AfterCreate(tx *gorm.DB) (err error) {
	inv := &Inventory{
		ID:           uuid.New(),
		ProductID:    p.ID,
		QtyAvailable: 0,
		QtyReserved:  0,
	}
	return tx.Create(inv).Error
}