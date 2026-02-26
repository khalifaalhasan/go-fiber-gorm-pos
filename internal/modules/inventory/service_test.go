package inventory_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-fiber-pos/internal/core"
	"go-fiber-pos/internal/modules/inventory"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Dummy Repo Implementation since mockgen is unavailable
type MockInventoryRepository struct {
	FindByProductIDFunc                func(ctx context.Context, productID uuid.UUID) (*core.Inventory, error)
	CreateDefaultFunc                  func(ctx context.Context, productID uuid.UUID) (*core.Inventory, error)
	GetMovementsFunc                   func(ctx context.Context, inventoryID uuid.UUID) ([]core.InventoryMovement, error)
	FindByProductIDForUpdateWithTxFunc func(tx *gorm.DB, productID uuid.UUID) (*core.Inventory, error)
	UpdateWithTxFunc                   func(tx *gorm.DB, inv *core.Inventory) error
	CreateMovementWithTxFunc           func(tx *gorm.DB, movement *core.InventoryMovement) error
	DBFunc                             func() *gorm.DB
}

func (m *MockInventoryRepository) DB() *gorm.DB {
	if m.DBFunc != nil {
		return m.DBFunc()
	}
	return nil
}

func (m *MockInventoryRepository) FindByProductID(ctx context.Context, productID uuid.UUID) (*core.Inventory, error) {
	if m.FindByProductIDFunc != nil {
		return m.FindByProductIDFunc(ctx, productID)
	}
	return nil, nil
}

func (m *MockInventoryRepository) CreateDefault(ctx context.Context, productID uuid.UUID) (*core.Inventory, error) {
	if m.CreateDefaultFunc != nil {
		return m.CreateDefaultFunc(ctx, productID)
	}
	return nil, nil
}

func (m *MockInventoryRepository) GetMovements(ctx context.Context, inventoryID uuid.UUID) ([]core.InventoryMovement, error) {
	if m.GetMovementsFunc != nil {
		return m.GetMovementsFunc(ctx, inventoryID)
	}
	return nil, nil
}

func (m *MockInventoryRepository) FindByProductIDForUpdateWithTx(tx *gorm.DB, productID uuid.UUID) (*core.Inventory, error) {
	if m.FindByProductIDForUpdateWithTxFunc != nil {
		return m.FindByProductIDForUpdateWithTxFunc(tx, productID)
	}
	return nil, nil
}

func (m *MockInventoryRepository) UpdateWithTx(tx *gorm.DB, inv *core.Inventory) error {
	if m.UpdateWithTxFunc != nil {
		return m.UpdateWithTxFunc(tx, inv)
	}
	return nil
}

func (m *MockInventoryRepository) CreateMovementWithTx(tx *gorm.DB, movement *core.InventoryMovement) error {
	if m.CreateMovementWithTxFunc != nil {
		return m.CreateMovementWithTxFunc(tx, movement)
	}
	return nil
}

func TestGetStockByProductID(t *testing.T) {
	productID := uuid.New()
	inventoryID := uuid.New()

	testCases := []struct {
		name          string
		productID     uuid.UUID
		setupMock     func(*MockInventoryRepository)
		expectedError error
		expectedQty   int
	}{
		{
			name:      "Sukses Ambil Stok",
			productID: productID,
			setupMock: func(mock *MockInventoryRepository) {
				mock.FindByProductIDFunc = func(ctx context.Context, pid uuid.UUID) (*core.Inventory, error) {
					return &core.Inventory{
						ID:           inventoryID,
						ProductID:    pid,
						QtyAvailable: 15,
						UpdatedAt:    time.Now(),
					}, nil
				}
			},
			expectedError: nil,
			expectedQty:   15,
		},
		{
			name:      "Gagal - Stok Tidak Ditemukan",
			productID: productID,
			setupMock: func(mock *MockInventoryRepository) {
				mock.FindByProductIDFunc = func(ctx context.Context, pid uuid.UUID) (*core.Inventory, error) {
					return nil, gorm.ErrRecordNotFound
				}
			},
			expectedError: core.ErrNotFound,
			expectedQty:   0,
		},
		{
			name:      "Gagal - Internal Server Error",
			productID: productID,
			setupMock: func(mock *MockInventoryRepository) {
				mock.FindByProductIDFunc = func(ctx context.Context, pid uuid.UUID) (*core.Inventory, error) {
					return nil, errors.New("db disconnect")
				}
			},
			expectedError: core.ErrInternalServer,
			expectedQty:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockInventoryRepository{}
			tc.setupMock(mockRepo)

			v := validator.New()
			service := inventory.NewInventoryService(mockRepo, v)

			res, err := service.GetStockByProductID(context.Background(), tc.productID)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tc.expectedQty, res.QtyAvailable)
			}
		})
	}
}

func TestCreateDefaultStock(t *testing.T) {
	productID := uuid.New()

	mockRepo := &MockInventoryRepository{
		CreateDefaultFunc: func(ctx context.Context, pid uuid.UUID) (*core.Inventory, error) {
			return &core.Inventory{ProductID: pid, QtyAvailable: 0}, nil
		},
	}

	v := validator.New()
	service := inventory.NewInventoryService(mockRepo, v)

	err := service.CreateDefaultStock(context.Background(), productID)
	assert.NoError(t, err)
}
