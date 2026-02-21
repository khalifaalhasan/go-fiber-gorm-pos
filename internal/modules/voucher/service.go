package voucher

import (
	"errors"
	"time"

	"go-fiber-pos/internal/core"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type voucherService struct {
	repo VoucherRepository
	v    *validator.Validate
}

func NewVoucherService(repo VoucherRepository, v *validator.Validate) VoucherService {
	return &voucherService{repo: repo, v: v}
}

func (s *voucherService) CreateVoucher(req CreateVoucherRequest) (*core.Voucher, error) {
	if err := s.v.Struct(req); err != nil {
		return nil, err
	}

	// Cek duplikasi kode voucher
	existing, err := s.repo.FindByCode(req.Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, core.ErrInternalServer
	}
	if existing != nil {
		return nil, core.ErrAlreadyExists
	}

	// Validasi: ValidUntil harus di masa depan
	if req.ValidUntil.Before(time.Now()) {
		return nil, core.ErrVoucherInvalid
	}

	voucher := &core.Voucher{
		ID:                uuid.New(),
		Code:              req.Code,
		DiscountType:      req.DiscountType,
		DiscountValue:     req.DiscountValue,
		MinOrderAmount:    req.MinOrderAmount,
		MaxDiscountAmount: req.MaxDiscountAmount,
		ValidUntil:        req.ValidUntil,
		IsActive:          true,
	}

	if err := s.repo.Create(voucher); err != nil {
		return nil, core.ErrInternalServer
	}
	return voucher, nil
}

func (s *voucherService) GetAllVouchers() ([]core.Voucher, error) {
	vouchers, err := s.repo.GetAll()
	if err != nil {
		return nil, core.ErrInternalServer
	}
	return vouchers, nil
}

func (s *voucherService) DeleteVoucher(id uuid.UUID) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return core.ErrNotFound
	}
	if err := s.repo.Delete(id); err != nil {
		return core.ErrInternalServer
	}
	return nil
}
