package product

// Sesuaikan dengan nama module di go.mod kamu
import (
	"go-fiber-pos/internal/core"
)


//go:generate mockgen -destination=mocks/mock_repository.go -package=mocks -source=contract.go ProductRepository

type ProductRepository interface {
	Create(product *core.Product) error
	GetAll() ([]core.Product, error)
	FindByName(name string) (*core.Product, error)
}

type ProductService interface {
	// Lihat! Sekarang dia menerima tipe dari package dto
	CreateProduct(req CreateProductRequest) (*core.Product, error) 
	GetAllProducts() ([]core.Product, error)
}