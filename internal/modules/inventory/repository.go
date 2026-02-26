package inventory

import (
	"context"

	"go-fiber-pos/internal/core"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) DB() *gorm.DB {
	return r.db
}

func (r *inventoryRepository) FindByProductID(ctx context.Context, productID uuid.UUID) (*core.Inventory, error) {
	var inv core.Inventory
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).First(&inv).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *inventoryRepository) CreateDefault(ctx context.Context, productID uuid.UUID) (*core.Inventory, error) {
	inv := &core.Inventory{
		ID:           uuid.New(),
		ProductID:    productID,
		QtyAvailable: 0,
		QtyReserved:  0,
	}
	if err := r.db.WithContext(ctx).Create(inv).Error; err != nil {
		return nil, err
	}
	return inv, nil
}

func (r *inventoryRepository) GetMovements(ctx context.Context, inventoryID uuid.UUID) ([]core.InventoryMovement, error) {
	var movements []core.InventoryMovement
	if err := r.db.WithContext(ctx).Where("inventory_id = ?", inventoryID).Order("created_at desc").Find(&movements).Error; err != nil {
		return nil, err
	}
	return movements, nil
}

// Transactional methods

func (r *inventoryRepository) FindByProductIDForUpdateWithTx(tx *gorm.DB, productID uuid.UUID) (*core.Inventory, error) {
	var inv core.Inventory
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("product_id = ?", productID).First(&inv).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *inventoryRepository) UpdateWithTx(tx *gorm.DB, inv *core.Inventory) error {
	return tx.Save(inv).Error
}

func (r *inventoryRepository) CreateMovementWithTx(tx *gorm.DB, movement *core.InventoryMovement) error {
	return tx.Create(movement).Error
}
