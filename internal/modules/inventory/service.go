package inventory

import (
	"context"
	"errors"

	"go-fiber-pos/internal/core"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type inventoryService struct {
	repo InventoryRepository
	v    *validator.Validate
}

func NewInventoryService(repo InventoryRepository, v *validator.Validate) InventoryService {
	return &inventoryService{repo: repo, v: v}
}

func (s *inventoryService) GetStockByProductID(ctx context.Context, productID uuid.UUID) (*InventoryResponse, error) {
	inv, err := s.repo.FindByProductID(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrNotFound
		}
		return nil, core.ErrInternalServer
	}

	return &InventoryResponse{
		ID:           inv.ID,
		ProductID:    inv.ProductID,
		QtyAvailable: inv.QtyAvailable,
		QtyReserved:  inv.QtyReserved,
		UpdatedAt:    inv.UpdatedAt,
	}, nil
}

func (s *inventoryService) AdjustStock(ctx context.Context, req AdjustStockRequest) error {
	if err := s.v.Struct(req); err != nil {
		return err
	}

	tx := s.repo.DB().WithContext(ctx).Begin()
	if tx.Error != nil {
		return core.ErrInternalServer
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	inv, err := s.repo.FindByProductIDForUpdateWithTx(tx, req.ProductID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return core.ErrNotFound
		}
		return core.ErrInternalServer
	}

	inv.QtyAvailable += req.QtyChange

	if err := s.repo.UpdateWithTx(tx, inv); err != nil {
		tx.Rollback()
		return core.ErrInternalServer
	}

	movement := &core.InventoryMovement{
		ID:            uuid.New(),
		InventoryID:   inv.ID,
		ReferenceType: "ADJUSTMENT",
		ReferenceID:   "MANUAL_ADJUSTMENT",
		QtyChange:     req.QtyChange,
		Notes:         req.Notes,
	}

	if err := s.repo.CreateMovementWithTx(tx, movement); err != nil {
		tx.Rollback()
		return core.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		return core.ErrInternalServer
	}

	return nil
}

func (s *inventoryService) GetMovements(ctx context.Context, productID uuid.UUID) ([]InventoryMovementResponse, error) {
	inv, err := s.repo.FindByProductID(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrNotFound
		}
		return nil, core.ErrInternalServer
	}

	movements, err := s.repo.GetMovements(ctx, inv.ID)
	if err != nil {
		return nil, core.ErrInternalServer
	}

	var res []InventoryMovementResponse
	for _, m := range movements {
		res = append(res, InventoryMovementResponse{
			ID:            m.ID,
			InventoryID:   m.InventoryID,
			ReferenceType: m.ReferenceType,
			ReferenceID:   m.ReferenceID,
			QtyChange:     m.QtyChange,
			Notes:         m.Notes,
			CreatedAt:     m.CreatedAt,
		})
	}
	return res, nil
}

func (s *inventoryService) CheckStockWithTx(ctx context.Context, tx *gorm.DB, productID uuid.UUID, qty int) error {
	inv, err := s.repo.FindByProductIDForUpdateWithTx(tx, productID)
	if err != nil {
		return err
	}

	if inv.QtyAvailable < qty {
		return core.ErrInsufficientStock
	}

	return nil
}

func (s *inventoryService) DeductStockWithTx(ctx context.Context, tx *gorm.DB, productID uuid.UUID, qty int, referenceType, referenceID string) error {
	inv, err := s.repo.FindByProductIDForUpdateWithTx(tx, productID)
	if err != nil {
		return err
	}

	if inv.QtyAvailable < qty {
		return core.ErrInsufficientStock
	}

	inv.QtyAvailable -= qty

	if err := s.repo.UpdateWithTx(tx, inv); err != nil {
		return err
	}

	movement := &core.InventoryMovement{
		ID:            uuid.New(),
		InventoryID:   inv.ID,
		ReferenceType: referenceType,
		ReferenceID:   referenceID,
		QtyChange:     -qty,
		Notes:         "Order Checkout",
	}

	if err := s.repo.CreateMovementWithTx(tx, movement); err != nil {
		return err
	}

	return nil
}

func (s *inventoryService) CreateDefaultStock(ctx context.Context, productID uuid.UUID) error {
	_, err := s.repo.CreateDefault(ctx, productID)
	return err
}
