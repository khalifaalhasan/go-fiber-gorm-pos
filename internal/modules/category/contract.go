package category

import "go-fiber-pos/internal/core"

type CategoryRepository interface {
	Create(category *core.Category) error
	GetAll() ([]core.Category, error)
	FindByName(name string) (*core.Category, error)
}


type CategoryService interface {
	CreateCategory(req CreateCategoryRequest) (*core.Category, error) 
	GetAllCategories() ([]core.Category, error)                       
}