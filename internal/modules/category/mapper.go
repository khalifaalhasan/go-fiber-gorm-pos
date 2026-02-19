package category

import "go-fiber-pos/internal/core"


func ToCategoryResponse(domain *core.Category) CategoryResponse {
	return CategoryResponse{
		ID:   domain.ID,
		Name: domain.Name,
		Slug: domain.Slug,
	}
}

// Mapper Function untuk List
func ToCategoryResponseList(domains []core.Category) []CategoryResponse {
	responses := []CategoryResponse{}
	for _, domain := range domains {
		responses = append(responses, ToCategoryResponse(&domain))
	}
	return responses
}