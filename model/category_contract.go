package model

import "github.com/google/uuid"

// DTO (Data Transfer Object) untuk Validasi Input
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=3"`
}

// Kontrak untuk Repository (Hanya ngobrol ke Database)
type CategoryRepository interface {
	Create(category *Category) error
	FindByStoreID(storeID uuid.UUID) ([]Category, error)
	FindByNameAndStoreID(name string, storeID uuid.UUID) (*Category, error)
}

// Kontrak untuk Service (Otak Bisnis)
type CategoryService interface {
	CreateCategory(storeID uuid.UUID, req CreateCategoryRequest) (*Category, error)
	GetCategoriesByStore(storeID uuid.UUID) ([]Category, error)
}

// Response DTO untuk Kategori (Biar rapi di JSON)
type CategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// Mapper Function: Domain -> Response
func ToCategoryResponse(domain *Category) CategoryResponse {
	return CategoryResponse{
		ID:   domain.ID,
		Name: domain.Name,
	}
}

// Mapper Function untuk List (dipakai di GetAll)
func ToCategoryResponseList(domains []Category) []CategoryResponse {
    responses := []CategoryResponse{} // Inisialisasi slice kosong, bukan nil
    for _, domain := range domains {
        responses = append(responses, ToCategoryResponse(&domain))
    }
    return responses
}