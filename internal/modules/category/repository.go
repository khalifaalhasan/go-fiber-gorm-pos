package category

import (
	model "go-fiber-pos/internal/core"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetAll() ([]model.Category, error) {
	var categories []model.Category
	// Langsung sikat semua data, karena ini database milik 1 toko eksklusif
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindByName(name string) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}