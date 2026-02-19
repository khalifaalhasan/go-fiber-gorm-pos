package category

import (
	"errors"

	"go-fiber-pos/internal/core"
	"go-fiber-pos/pkg/validator"

	"github.com/google/uuid"
)

type categoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(req CreateCategoryRequest) (*core.Category, error) {
	// 1. Validasi Input
	if err := validator.Validate.Struct(req); err != nil {
		return nil, errors.New("validasi gagal: " + err.Error())
	}

	// 2. Cek Duplikat (Panggil FindByName yang baru)
	existing, _ := s.repo.FindByName(req.Name)
	if existing != nil {
		return nil, errors.New("kategori dengan nama tersebut sudah ada")
	}

	// 3. Map ke Domain (Tanpa StoreID)
	category := &core.Category{
		ID:   uuid.New(),
		Name: req.Name,
	}

	// 4. Simpan ke DB
	err := s.repo.Create(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetAllCategories() ([]core.Category, error) {
	return s.repo.GetAll()
}