package store

import "go-fiber-pos/internal/core"

// StoreRepository mendefinisikan kontrak akses data untuk StoreProfile.
type StoreRepository interface {
	GetProfile() (*core.StoreProfile, error)
	Upsert(profile *core.StoreProfile) (*core.StoreProfile, error)
}

// StoreService mendefinisikan kontrak business logic untuk StoreProfile.
type StoreService interface {
	GetProfile() (*core.StoreProfile, error)
	UpdateProfile(req UpdateStoreRequest) (*core.StoreProfile, error)
}
