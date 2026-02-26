package inventory

import (
	"context"

	"go-fiber-pos/internal/core"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	DB() *gorm.DB
	FindByProductID(ctx context.Context, productID uuid.UUID) (*core.Inventory, error)
	CreateDefault(ctx context.Context, productID uuid.UUID) (*core.Inventory, error)
	GetMovements(ctx context.Context, inventoryID uuid.UUID) ([]core.InventoryMovement, error)
	
	// Transactional methods
	FindByProductIDForUpdateWithTx(tx *gorm.DB, productID uuid.UUID) (*core.Inventory, error)
	UpdateWithTx(tx *gorm.DB, inv *core.Inventory) error
	CreateMovementWithTx(tx *gorm.DB, movement *core.InventoryMovement) error
}

type InventoryService interface {
	GetStockByProductID(ctx context.Context, productID uuid.UUID) (*InventoryResponse, error)
	AdjustStock(ctx context.Context, req AdjustStockRequest) error
	GetMovements(ctx context.Context, productID uuid.UUID) ([]InventoryMovementResponse, error)
	
	// Called internally by Order Module during checkout
	CheckStockWithTx(ctx context.Context, tx *gorm.DB, productID uuid.UUID, qty int) error

	// Called internally by Order Module after successful payment
	DeductStockWithTx(ctx context.Context, tx *gorm.DB, productID uuid.UUID, qty int, referenceType, referenceID string) error
	
	// Called internally by Product Module during creation
	CreateDefaultStock(ctx context.Context, productID uuid.UUID) error
}
