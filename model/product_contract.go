package model

import "github.com/google/uuid"

// ==========================================
// DTO (REQUEST & RESPONSE)
// ==========================================

// CreateProductRequest: Ini yang akan ditangkap dari JSON Body Postman/Frontend
type CreateProductRequest struct {
	CategoryID     uuid.UUID `json:"category_id" validate:"required"`
	Name           string    `json:"name" validate:"required,min=3"`
	Description    string    `json:"description" validate:"required,min=10"`
	ImageURL       string    `json:"image_url"`
	NormalPrice    int       `json:"normal_price" validate:"required,min=0"`
	IsAvailable	   bool      `json:"is_available"`
	
	// Opsional
	IsPromoActive  bool      `json:"is_promo_active"`
	PromoPrice     int       `json:"promo_price"`
	PromoStartTime string    `json:"promo_start_time"`
	PromoEndTime   string    `json:"promo_end_time"`
}

// ProductResponse: Kardus rapi untuk dikirim ke JSON Frontend (menyembunyikan data rahasia DB)
type ProductResponse struct {
	ID             uuid.UUID `json:"id"`
	CategoryID     uuid.UUID `json:"category_id"`
	Name           string    `json:"name"`
	Slug		   string    `json:"slug"`
	Description    string    `json:"description"`
	ImageURL       string    `json:"image_url"`
	NormalPrice    int       `json:"normal_price"`
	IsAvailable    bool      `json:"is_available"`
	IsPromoActive  bool      `json:"is_promo_active"`
	PromoPrice     int       `json:"promo_price"`
	PromoStartTime string    `json:"promo_start_time"`
	PromoEndTime   string    `json:"promo_end_time"`
}

// ==========================================
// INTERFACE (KONTRAK)
// ==========================================

// Kontrak untuk Repository (ngobrol ke DB)
type ProductRepository interface {
	Create(product *Product) error
	GetAll() ([]Product, error)
	FindByName(name string) (*Product, error)
}

// Kontrak untuk Service (Otak Bisnis)
type ProductService interface {
	CreateProduct(req CreateProductRequest) (*Product, error)
	GetAllProducts() ([]Product, error)

}

// ==========================================
// MAPPER FUNCTIONS (BEST PRACTICE)
// ==========================================

// ToProductResponse: Mengubah Domain GORM menjadi Response DTO
func ToProductResponse(domain *Product) ProductResponse {
	return ProductResponse{
		ID:             domain.ID,
		CategoryID:     domain.CategoryID,
		Name:           domain.Name,
		Slug:		    domain.Slug,
		Description:    domain.Description,
		ImageURL:       domain.ImageURL,
		NormalPrice:    domain.NormalPrice,
		IsAvailable:    domain.IsAvailable,
		IsPromoActive:  domain.IsPromoActive,
		PromoPrice:     domain.PromoPrice,
		PromoStartTime: domain.PromoStartTime,
		PromoEndTime:   domain.PromoEndTime,
	}
}

// ToProductResponseList: Mapper Function untuk data array/list (dipakai di endpoint GetAll)
func ToProductResponseList(domains []Product) []ProductResponse {
    // Inisialisasi dengan slice kosong, jangan nil
    responses := []ProductResponse{} 
    
    for _, domain := range domains {
        responses = append(responses, ToProductResponse(&domain))
    }
    return responses
}