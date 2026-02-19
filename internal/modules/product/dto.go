package product

import "github.com/google/uuid"

// Request DTO (Dari Frontend ke Server)
type CreateProductRequest struct {
	CategoryID     uuid.UUID `json:"category_id" validate:"required"`
	Name           string    `json:"name" validate:"required,min=3"`
	Description    string    `json:"description" validate:"required,min=10"`
	ImageURL       string    `json:"image_url"`
	NormalPrice    int       `json:"normal_price" validate:"required,min=0"`
	IsAvailable    bool      `json:"is_available"`

	IsPromoActive  bool   `json:"is_promo_active"`
	PromoPrice     int    `json:"promo_price"`
	PromoStartTime string `json:"promo_start_time"`
	PromoEndTime   string `json:"promo_end_time"`
}

// Response DTO (Dari Server ke Frontend)
type ProductResponse struct {
	ID             uuid.UUID `json:"id"`
	CategoryID     uuid.UUID `json:"category_id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	Description    string    `json:"description"`
	ImageURL       string    `json:"image_url"`
	NormalPrice    int       `json:"normal_price"`
	IsAvailable    bool      `json:"is_available"`
	IsPromoActive  bool      `json:"is_promo_active"`
	PromoPrice     int       `json:"promo_price"`
	PromoStartTime string    `json:"promo_start_time"`
	PromoEndTime   string    `json:"promo_end_time"`
}