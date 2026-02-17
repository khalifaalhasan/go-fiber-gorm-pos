package repository

import (
	"go-fiber-pos/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) model.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByStoreID(storeID uuid.UUID) ([]model.Product, error) {
	var product []model.Product
	err := r.db.Where("store_id = ?", storeID).Find(&product).Error 
	return product, err
}

