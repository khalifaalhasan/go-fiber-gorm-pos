package category

import "github.com/google/uuid"


type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=3"`
}

type CategoryResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}