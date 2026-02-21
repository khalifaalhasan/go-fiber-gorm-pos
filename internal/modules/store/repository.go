package store

import (
	"go-fiber-pos/internal/core"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db: db}
}

// GetProfile mengambil satu-satunya baris StoreProfile yang ada.
func (r *storeRepository) GetProfile() (*core.StoreProfile, error) {
	var profile core.StoreProfile
	err := r.db.First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// Upsert menyimpan profil toko. Jika belum ada (belum punya ID), buat baru. Jika ada, update.
func (r *storeRepository) Upsert(profile *core.StoreProfile) (*core.StoreProfile, error) {
	// Cek apakah sudah ada data
	var existing core.StoreProfile
	err := r.db.First(&existing).Error
	if err != nil {
		// Belum ada — buat baru
		profile.ID = uuid.New()
		if createErr := r.db.Create(profile).Error; createErr != nil {
			return nil, createErr
		}
		return profile, nil
	}

	// Sudah ada — update kolom yang relevan saja
	existing.Name = profile.Name
	existing.Address = profile.Address
	existing.Phone = profile.Phone
	existing.MarkupFee = profile.MarkupFee
	if saveErr := r.db.Save(&existing).Error; saveErr != nil {
		return nil, saveErr
	}
	return &existing, nil
}
