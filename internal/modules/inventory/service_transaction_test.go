package inventory_test

import (
	"context"
	"testing"
	"time"

	"go-fiber-pos/internal/core"
	"go-fiber-pos/internal/modules/inventory"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeductStockWithTx(t *testing.T) {
	productID := uuid.New()
	inventoryID := uuid.New()
	tx := &gorm.DB{} // Mock tx

	testCases := []struct {
		name          string
		qty           int
		setupMock     func(*MockInventoryRepository)
		expectedError error
	}{
		{
			name: "Sukses Deduct Stok",
			qty:  5,
			setupMock: func(mock *MockInventoryRepository) {
				mock.FindByProductIDForUpdateWithTxFunc = func(tx *gorm.DB, pid uuid.UUID) (*core.Inventory, error) {
					return &core.Inventory{
						ID:           inventoryID,
						ProductID:    pid,
						QtyAvailable: 10,
					}, nil
				}
				mock.UpdateWithTxFunc = func(tx *gorm.DB, inv *core.Inventory) error {
					assert.Equal(t, 5, inv.QtyAvailable)
					return nil
				}
				mock.CreateMovementWithTxFunc = func(tx *gorm.DB, mov *core.InventoryMovement) error {
					assert.Equal(t, -5, mov.QtyChange)
					assert.Equal(t, "ORDER", mov.ReferenceType)
					return nil
				}
			},
			expectedError: nil,
		},
		{
			name: "Gagal - Stok Tidak Cukup",
			qty:  15,
			setupMock: func(mock *MockInventoryRepository) {
				mock.FindByProductIDForUpdateWithTxFunc = func(tx *gorm.DB, pid uuid.UUID) (*core.Inventory, error) {
					return &core.Inventory{
						ID:           inventoryID,
						ProductID:    pid,
						QtyAvailable: 10,
					}, nil
				}
			},
			expectedError: core.ErrInsufficientStock,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockInventoryRepository{}
			tc.setupMock(mockRepo)

			service := inventory.NewInventoryService(mockRepo, nil)

			err := service.DeductStockWithTx(context.Background(), tx, productID, tc.qty, "ORDER", "Ref-123")

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetMovements(t *testing.T) {
	productID := uuid.New()
	inventoryID := uuid.New()

	mockRepo := &MockInventoryRepository{
		FindByProductIDFunc: func(ctx context.Context, pid uuid.UUID) (*core.Inventory, error) {
			return &core.Inventory{ID: inventoryID, ProductID: pid}, nil
		},
		GetMovementsFunc: func(ctx context.Context, id uuid.UUID) ([]core.InventoryMovement, error) {
			return []core.InventoryMovement{
				{
					ID:            uuid.New(),
					InventoryID:   id,
					ReferenceType: "RESTOCK",
					QtyChange:     10,
					CreatedAt:     time.Now(),
				},
				{
					ID:            uuid.New(),
					InventoryID:   id,
					ReferenceType: "ORDER",
					QtyChange:     -2,
					CreatedAt:     time.Now().Add(-1 * time.Hour),
				},
			}, nil
		},
	}

	service := inventory.NewInventoryService(mockRepo, nil)

	res, err := service.GetMovements(context.Background(), productID)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "RESTOCK", res[0].ReferenceType)
	assert.Equal(t, 10, res[0].QtyChange)
	assert.Equal(t, "ORDER", res[1].ReferenceType)
	assert.Equal(t, -2, res[1].QtyChange)
}
