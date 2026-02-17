package model

// DTO untuk validasi input JSON dari Postman/React
type RegisterRequest struct {
	StoreName string `json:"store_name" validate:"required,min=3"`
	Subdomain string `json:"subdomain" validate:"required,alphanum,min=3"` // alphanum agar tidak ada spasi untuk URL
	Name      string `json:"name" validate:"required,min=3"`
	Username  string `json:"username" validate:"required,alphanum,min=4"`
	Password  string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Interface Repository (Hanya ngobrol ke DB)
type AuthRepository interface {
	CreateStoreAndUser(store *Store, user *User) error
	FindByUsername(username string) (*User, error)
}

// Interface Service (Otak Bisnis)
type AuthService interface {
	Register(req RegisterRequest) error
	Login(req LoginRequest) (string, error)
}