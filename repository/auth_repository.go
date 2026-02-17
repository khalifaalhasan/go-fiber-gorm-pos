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

// CreateStoreAndUser menggunakan GORM Transaction untuk jaminan integritas data
func (r *authRepository) CreateStoreAndUser(store *model.Store, user *model.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan Toko dulu
		if err := tx.Create(store).Error; err != nil {
			return err // Otomatis Rollback
		}

		// 2. Assign ID toko yang baru jadi ke User
		user.StoreID = store.ID

		// 3. Simpan User
		if err := tx.Create(user).Error; err != nil {
			return err // Otomatis Rollback beserta toko di atas
		}

		return nil // Otomatis Commit
	})
}

func (r *authRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}