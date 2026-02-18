package repository

import (
	"go-fiber-pos/model"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) model.AuthRepository {
	return &authRepository{db: db}
}

// CreateUser sekarang sangat simpel. Tidak perlu lagi DB Transaction 
// untuk membuat Store, karena sistem ini Single-Tenant.
func (r *authRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByUsername tidak berubah, logikanya tetap sama
func (r *authRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}