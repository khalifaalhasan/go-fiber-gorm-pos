package voucher

import (
	"go-fiber-pos/internal/core"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type voucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) VoucherRepository {
	return &voucherRepository{db: db}
}

func (r *voucherRepository) Create(voucher *core.Voucher) error {
	return r.db.Create(voucher).Error
}

func (r *voucherRepository) GetAll() ([]core.Voucher, error) {
	var vouchers []core.Voucher
	err := r.db.Order("created_at DESC").Find(&vouchers).Error
	return vouchers, err
}

func (r *voucherRepository) FindByID(id uuid.UUID) (*core.Voucher, error) {
	var voucher core.Voucher
	err := r.db.First(&voucher, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &voucher, nil
}

func (r *voucherRepository) FindByCode(code string) (*core.Voucher, error) {
	var voucher core.Voucher
	err := r.db.Where("code = ?", code).First(&voucher).Error
	if err != nil {
		return nil, err
	}
	return &voucher, nil
}

func (r *voucherRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&core.Voucher{}, "id = ?", id).Error
}
