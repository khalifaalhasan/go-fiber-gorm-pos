package store

// UpdateStoreRequest adalah DTO untuk request update profil toko.
type UpdateStoreRequest struct {
	Name      string `json:"name" validate:"required,min=3"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	MarkupFee int    `json:"markup_fee" validate:"min=0"`
}

// StoreResponse adalah DTO untuk response profil toko.
type StoreResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	MarkupFee int    `json:"markup_fee"`
}
