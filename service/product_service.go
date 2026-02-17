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

func (s *productService) CreateProduct(storeID uuid.UUID, req model.CreateProductRequest) (*model.Product, error) {
	if err := utils.Validate.Struct(req); err != nil {
		return nil, errors.New("validasi gagal: " + err.Error())
	}

	product := &model.Product{
		StoreID:        storeID, // Didapat dari parameter (hasil bongkar Token JWT)
		CategoryID:     req.CategoryID,
		Name:           req.Name,
		Description:    req.Description,
		ImageURL:       req.ImageURL,
		NormalPrice:    req.NormalPrice,
		IsAvailable:    true, // Defaultnya selalu true saat pertama bikin
		IsPromoActive:  req.IsPromoActive,
		PromoPrice:     req.PromoPrice,
		PromoStartTime: req.PromoStartTime,
		PromoEndTime:   req.PromoEndTime,
	}

	err := s.repo.Create(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// get product
func (s *productService) GetProductsByStore(storeID uuid.UUID) ([]model.Product, error){
	return s.repo.FindByStoreID(storeID)
}