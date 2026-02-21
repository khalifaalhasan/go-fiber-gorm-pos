package voucher

import (
	"go-fiber-pos/internal/core"

	"github.com/google/uuid"
)

// VoucherRepository mendefinisikan kontrak akses data untuk Voucher.
type VoucherRepository interface {
	Create(voucher *core.Voucher) error
	GetAll() ([]core.Voucher, error)
	FindByID(id uuid.UUID) (*core.Voucher, error)
	FindByCode(code string) (*core.Voucher, error)
	Delete(id uuid.UUID) error
}

// VoucherService mendefinisikan kontrak business logic untuk Voucher.
type VoucherService interface {
	CreateVoucher(req CreateVoucherRequest) (*core.Voucher, error)
	GetAllVouchers() ([]core.Voucher, error)
	DeleteVoucher(id uuid.UUID) error
}
