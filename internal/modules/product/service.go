package product

import (
	"errors"
	"go-fiber-pos/internal/core"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type productService struct {
	repo ProductRepository
	v    *validator.Validate
}

func NewProductService(repo ProductRepository, v *validator.Validate) ProductService {
	return &productService{
		repo: repo,
		v:    v,
	}
}

func (s *productService) CreateProduct(req CreateProductRequest) (*core.Product, error) {
	// 1. Validasi Input dari DTO
	if err := s.v.Struct(req); err != nil {
		return nil, err
	}

	// 2. Cek Duplikasi Produk
	// FIX: Menggunakan s.repo langsung
	existing, err := s.repo.FindByName(req.Name)
	if err != nil {
		return nil, errors.New("terjadi kesalahan pada server")
	}

	if existing != nil {
		// FIX: String disamakan persis dengan expectedError di service_test.go
		return nil, errors.New("produk sudah ada") 
	}

	// 3. Mapping Request ke Entity Core
	product := &core.Product{
		ID:             uuid.New(),
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

	// 4. Simpan ke Database
	// FIX: Menggunakan s.repo langsung
	err = s.repo.Create(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// get product
func (s *productService) GetAllProducts() ([]core.Product, error) {
	// FIX: Menggunakan s.repo langsung
	return s.repo.GetAll()
}