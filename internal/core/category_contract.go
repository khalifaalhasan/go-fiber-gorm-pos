package core

import "github.com/google/uuid"

// DTO untuk Validasi Input
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=3"`
}

// Kontrak untuk Repository (Hanya ngobrol ke Database)
type CategoryRepository interface {
	Create(category *Category) error
	GetAll() ([]Category, error)             // Berubah: Mengambil semua kategori di database
	FindByName(name string) (*Category, error) // Berubah: Cukup cari berdasarkan nama saja
}

// Kontrak untuk Service (Otak Bisnis)
type CategoryService interface {
	CreateCategory(req CreateCategoryRequest) (*Category, error) // Berubah: storeID dihapus
	GetAllCategories() ([]Category, error)                       // Berubah: storeID dihapus
}

// Response DTO (Kardus rapi)
type CategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

// Mapper Function: Domain -> Response
func ToCategoryResponse(domain *Category) CategoryResponse {
	return CategoryResponse{
		ID:   domain.ID,
		Name: domain.Name,
		Slug: domain.Slug,
	}
}

// Mapper Function untuk List
func ToCategoryResponseList(domains []Category) []CategoryResponse {
	responses := []CategoryResponse{}
	for _, domain := range domains {
		responses = append(responses, ToCategoryResponse(&domain))
	}
	return responses
}