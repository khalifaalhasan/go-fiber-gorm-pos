package store

import (
	"errors"
	"go-fiber-pos/internal/core"

	"github.com/go-playground/validator/v10"
)

type storeService struct {
	repo StoreRepository
	v    *validator.Validate
}

func NewStoreService(repo StoreRepository, v *validator.Validate) StoreService {
	return &storeService{repo: repo, v: v}
}

func (s *storeService) GetProfile() (*core.StoreProfile, error) {
	profile, err := s.repo.GetProfile()
	if err != nil {
		// Profil belum dibuat adalah kondisi valid untuk toko baru
		if errors.Is(err, core.ErrNotFound) {
			return nil, core.ErrNotFound
		}
		return nil, core.ErrInternalServer
	}
	return profile, nil
}

func (s *storeService) UpdateProfile(req UpdateStoreRequest) (*core.StoreProfile, error) {
	if err := s.v.Struct(req); err != nil {
		return nil, err
	}

	profile := &core.StoreProfile{
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		MarkupFee: req.MarkupFee,
	}

	result, err := s.repo.Upsert(profile)
	if err != nil {
		return nil, core.ErrInternalServer
	}
	return result, nil
}
