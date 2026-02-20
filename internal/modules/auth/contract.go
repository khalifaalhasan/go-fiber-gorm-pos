package auth

import "go-fiber-pos/internal/core"


type AuthRepository interface {
	// Berubah dari CreateStoreAndUser menjadi CreateUser saja
	CreateUser(user *core.User) error
	FindByUsername(username string) (*core.User, error)
}

// Interface Service (Otak Bisnis)
type AuthService interface {
	Register(req RegisterRequest) error
	Login(req LoginRequest) (string, error)
}