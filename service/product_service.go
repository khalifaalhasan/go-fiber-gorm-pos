package service

import (
	"errors"
	"go-fiber-pos/model"
	"go-fiber-pos/utils"

	"github.com/google/uuid"
)



type  productService struct {
	repo model.ProductRepository
}

func NewProductService(repo model.ProductRepository)model.ProductService{
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(req model.CreateProductRequest) (*model.Product, error) {
    if err := utils.Validate.Struct(req); err != nil {
        return nil, errors.New("validasi gagal: " + err.Error())
    }

	// cek duplikat product
	existing, err := s.repo.FindByName(req.Name)
    if err != nil {
        return nil, errors.New("terjadi kesalahan pada server")
    }
    
    if existing != nil {
        return nil, errors.New("product dengan nama tersebut sudah ada")
    }

	product := &model.Product{
		ID: uuid.New(),
		CategoryID:     req.CategoryID,
		Name:           req.Name,
		Description:    req.Description,
		ImageURL:       req.ImageURL,
		NormalPrice:    req.NormalPrice,
		IsAvailable:    true, 
		IsPromoActive:  req.IsPromoActive,
		PromoPrice:     req.PromoPrice,
		PromoStartTime: req.PromoStartTime,
		PromoEndTime:   req.PromoEndTime,
	}

	err = s.repo.Create(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// get product
func (s *productService) GetAllProducts() ([]model.Product, error){
	return s.repo.GetAll()
}

