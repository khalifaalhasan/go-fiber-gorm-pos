package service

import (
	"errors"
	"go-fiber-pos/model"
	"go-fiber-pos/utils"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo model.AuthRepository
}

func NewAuthService(repo model.AuthRepository) model.AuthService {
	return &authService{repo: repo}
}

func (s *authService) Register(req model.RegisterRequest) error {
	// 1. Validasi Input
	if err := utils.Validate.Struct(req); err != nil {
		return errors.New("validasi gagal, pastikan semua data terisi dengan benar (min. 6 karakter untuk password)")
	}

	// 2. Cek apakah username sudah ada
	_, err := s.repo.FindByUsername(req.Username)
	if err == nil {
		return errors.New("username sudah digunakan, silakan pilih yang lain")
	}

	// 3. Hash Password (Best Practice)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal memproses password")
	}

	// 4. Siapkan Entity Store & User
	store := &model.Store{
		Name:      req.StoreName,
		Subdomain: req.Subdomain,
	}

	user := &model.User{
		Name:         req.Name,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         "ADMIN", // Orang yang daftar pertama otomatis jadi Admin Kafe
	}

	// 5. Lempar ke Repository
	return s.repo.CreateStoreAndUser(store, user)
}

func (s *authService) Login(req model.LoginRequest) (string, error) {
	// 1. Cari user
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		// Security Best Practice: Jangan pernah kasih tahu apakah 'username' yang salah atau 'password' yang salah.
		return "", errors.New("username atau password salah") 
	}

	// 2. Cocokkan Password Asli dengan Hash di DB
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", errors.New("username atau password salah")
	}

	// 3. Buat JWT Token pakai fungsi yang kita buat sebelumnya
	token, err := utils.GenerateToken(user.ID, user.StoreID, user.Role)
	if err != nil {
		return "", errors.New("gagal membuat sesi token")
	}

	return token, nil
}