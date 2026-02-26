package product

import model "go-fiber-pos/internal/core"

// ToProductResponse: Domain GORM -> Response DTO
func ToProductResponse(domain *model.Product) ProductResponse {
	res := ProductResponse{
		ID:             domain.ID,
		CategoryID:     domain.CategoryID,
		Name:           domain.Name,
		Slug:           domain.Slug,
		Description:    domain.Description,
		ImageURL:       domain.ImageURL,
		NormalPrice:    domain.NormalPrice,
		IsAvailable:    domain.IsAvailable,
		IsPromoActive:  domain.IsPromoActive,
		PromoPrice:     domain.PromoPrice,
		PromoStartTime: domain.PromoStartTime,
		PromoEndTime:   domain.PromoEndTime,
	}
	
	if domain.Inventory != nil {
		res.Stock = domain.Inventory.QtyAvailable
	}
	
	return res
}

// ToProductResponseList: Array Domain -> Array Response DTO
func ToProductResponseList(domains []model.Product) []ProductResponse {
	responses := []ProductResponse{}

	for _, domain := range domains {
		responses = append(responses, ToProductResponse(&domain))
	}
	return responses
}