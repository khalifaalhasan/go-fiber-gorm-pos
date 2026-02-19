package auth

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