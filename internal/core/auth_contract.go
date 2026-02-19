package core

// DTO untuk validasi input JSON dari Postman/React
type RegisterRequest struct {
	// StoreName dan Subdomain DIHAPUS karena Single-Tenant
	Name     string `json:"name" validate:"required,min=3"`
	Username string `json:"username" validate:"required,alphanum,min=4"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Interface Repository (Hanya ngobrol ke DB)
type AuthRepository interface {
	// Berubah dari CreateStoreAndUser menjadi CreateUser saja
	CreateUser(user *User) error
	FindByUsername(username string) (*User, error)
}

// Interface Service (Otak Bisnis)
type AuthService interface {
	Register(req RegisterRequest) error
	Login(req LoginRequest) (string, error)
}