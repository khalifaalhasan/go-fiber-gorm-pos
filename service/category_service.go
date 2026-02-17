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

func (s *categoryService) CreateCategory(storeID uuid.UUID, req model.CreateCategoryRequest) (*model.Category, error) {
    // 1. Validasi Input
    if err := utils.Validate.Struct(req); err != nil {
        return nil, errors.New("validasi gagal: " + err.Error())
    }

    // 2. Cek apakah Kategori sudah ada (PENTING: Cek SEBELUM Create)
    existing, err := s.repo.FindByNameAndStoreID(req.Name, storeID)
    if err == nil && existing != nil {
        return nil, errors.New("kategori dengan nama tersebut sudah ada di toko kamu")
    }

    // 3. Map DTO ke Model & Generate ID manual biar aman
    category := &model.Category{
        ID:      uuid.New(), // Pakai google/uuid buat generate ID di level aplikasi
        StoreID: storeID,
        Name:    req.Name,
    }

    // 4. SIMPAN KE DATABASE (Ini yang tadi hilang, anjir!)
    err = s.repo.Create(category)
    if err != nil {
        return nil, err
    }

    return category, nil
}

// get category
func (s *categoryService) GetCategoriesByStore(storeID uuid.UUID) ([]model.Category, error) {
	return s.repo.FindByStoreID(storeID)
}