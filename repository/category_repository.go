package repository

import (
	"go-fiber-pos/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) model.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindByStoreID(storeID uuid.UUID) ([]model.Category, error) {
	var categories []model.Category
	// Filter ketat berdasarkan StoreID agar data tidak bocor antar kafe
	err := r.db.Where("store_id = ?", storeID).Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindByNameAndStoreID(name string, storeID uuid.UUID) (*model.Category, error) {
    var category model.Category
    err := r.db.Where("name = ? AND store_id = ?", name, storeID).First(&category).Error
    
    // Kalau tidak ketemu, jangan anggap sebagai error yang merusak sistem
    if err != nil {
        return nil, err
    }
    
    return &category, nil
}