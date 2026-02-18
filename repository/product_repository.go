package repository

import (
	"go-fiber-pos/model"

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



func (r* productRepository) GetAll() ([]model.Product, error){
	var products []model.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *productRepository) FindByName(name string) (*model.Product, error) {
    var product model.Product
    err := r.db.Where("name = ?", name).First(&product).Error
    
    if err != nil {
        // Jika errornya adalah 'Record Not Found', kembalikan nil tanpa error
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    
    return &product, nil
}



