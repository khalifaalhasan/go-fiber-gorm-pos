package core

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

// ==========================================
// APP CONFIGURATION (Hanya 1 Baris Data Nanti)
// ==========================================



// StoreProfile menggantikan Store. Dipakai untuk header struk kasir/UI.
type StoreProfile struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Address   string         `gorm:"type:text" json:"address"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	MarkupFee int            `gorm:"default:0" json:"markup_fee"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

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
	Name      string         `gorm:"type:varchar(255);unique;not null" json:"name"` // Langsung unique global
	Slug 	  string 		 `gorm:"type:varchar(255);unique;index" json:"slug"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Product struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CategoryID     uuid.UUID      `gorm:"type:uuid;not null" json:"category_id"`
	Name           string         `gorm:"type:varchar(255);not null" json:"name"` // Bisa ditambah unique kalau mau
	Slug           string         `gorm:"type:varchar(255);index" json:"slug"` // Bisa ditambah unique kalau mau
	Description    string         `gorm:"type:text" json:"description"`
	ImageURL       string         `gorm:"type:varchar(255)" json:"image_url"`
	NormalPrice    int            `gorm:"not null" json:"normal_price"`
	IsAvailable    bool           `gorm:"default:true" json:"is_available"`
	IsPromoActive  bool           `gorm:"default:false" json:"is_promo_active"`
	PromoPrice     int            `json:"promo_price"`
	PromoStartTime string         `gorm:"type:varchar(5)" json:"promo_start_time"`
	PromoEndTime   string         `gorm:"type:varchar(5)" json:"promo_end_time"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// ==========================================
// VOUCHERS
// ==========================================

type Voucher struct {
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Code              string    `gorm:"type:varchar(50);unique;not null" json:"code"`
	DiscountType      string    `gorm:"type:varchar(50);not null" json:"discount_type"` // PERCENTAGE / FIXED
	DiscountValue     int       `gorm:"not null" json:"discount_value"`
	MinOrderAmount    int       `gorm:"default:0" json:"min_order_amount"`
	MaxDiscountAmount int       `gorm:"default:0" json:"max_discount_amount"`
	ValidUntil        time.Time `json:"valid_until"`
	IsActive          bool      `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
}

// ==========================================
// TRANSACTIONAL
// ==========================================

type Order struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	VoucherID        *uuid.UUID `gorm:"type:uuid" json:"voucher_id"` // Pointer karena opsional
	TableNumber      *string    `gorm:"type:varchar(50)" json:"table_number"`
	OrderStatus      string     `gorm:"type:varchar(50);default:'PENDING'" json:"order_status"`
	PaymentStatus    string     `gorm:"type:varchar(50);default:'UNPAID'" json:"payment_status"`
	TotalBasePrice   int        `gorm:"not null" json:"total_base_price"`
	TotalDiscount    int        `gorm:"default:0" json:"total_discount"`
	PlatformFee      int        `gorm:"default:0" json:"platform_fee"`
	TotalFinalAmount int        `gorm:"not null" json:"total_final_amount"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	Voucher  *Voucher    `gorm:"foreignKey:VoucherID" json:"voucher,omitempty"` // Perbaikan: pointer biar aman kalau nil
	Items    []OrderItem `gorm:"foreignKey:OrderID" json:"items"`    
	Payments []Payment   `gorm:"foreignKey:OrderID" json:"payments"` 
}

// (OrderItem dan Payment tetap sama seperti kodemu sebelumnya)
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
	ID                    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrderID               uuid.UUID `gorm:"type:uuid;not null" json:"order_id"`
	PaymentMethod         string    `gorm:"type:varchar(50);not null" json:"payment_method"` 
	MidtransTransactionID *string   `gorm:"type:varchar(255)" json:"midtrans_transaction_id"`
	AmountPaid            int       `gorm:"not null" json:"amount_paid"`
	PaymentStatus         string    `gorm:"type:varchar(50);not null" json:"payment_status"` 
	PaidAt                time.Time `json:"paid_at"`
	CreatedAt             time.Time `json:"created_at"`
}

func (c *Category) BeforeSave(tx *gorm.DB) (err error) {
    // 3. Pastikan memanggil slug.Make (package.Fungsi)
    c.Slug = slug.Make(c.Name)
    return
}

func (p *Product) BeforeSave(tx *gorm.DB) (err error) {
    p.Slug = slug.Make(p.Name)
    return
}