package service

import (
	"errors"
	"go-fiber-pos/model"
	"go-fiber-pos/utils"

	"github.com/google/uuid"
)

type categoryService struct {
	repo model.CategoryRepository
}

func NewCategoryService(repo model.CategoryRepository) model.CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(req model.CreateCategoryRequest) (*model.Category, error) {
	// 1. Validasi Input
	if err := utils.Validate.Struct(req); err != nil {
		return nil, errors.New("validasi gagal: " + err.Error())
	}

	// 2. Cek Duplikat (Panggil FindByName yang baru)
	existing, _ := s.repo.FindByName(req.Name)
	if existing != nil {
		return nil, errors.New("kategori dengan nama tersebut sudah ada")
	}

	// 3. Map ke Domain (Tanpa StoreID)
	category := &model.Category{
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

func (s *categoryService) GetAllCategories() ([]model.Category, error) {
	return s.repo.GetAll()
}