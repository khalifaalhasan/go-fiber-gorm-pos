package inventory

import (
	"time"

	"github.com/google/uuid"
)

type AdjustStockRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	QtyChange int       `json:"qty_change" validate:"required,ne=0"`
	Notes     string    `json:"notes" validate:"required"`
}

type InventoryResponse struct {
	ID           uuid.UUID `json:"id"`
	ProductID    uuid.UUID `json:"product_id"`
	QtyAvailable int       `json:"qty_available"`
	QtyReserved  int       `json:"qty_reserved"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type InventoryMovementResponse struct {
	ID            uuid.UUID `json:"id"`
	InventoryID   uuid.UUID `json:"inventory_id"`
	ReferenceType string    `json:"reference_type"`
	ReferenceID   string    `json:"reference_id"`
	QtyChange     int       `json:"qty_change"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}
