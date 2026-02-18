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

	// 4. Siapkan Entity User (Store Resmi Dihapus karena Single-Tenant)
	user := &model.User{
		Name:         req.Name,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         "ADMIN", // Otomatis jadi Admin sistem
	}

	// 5. Lempar ke Repository (Panggil method CreateUser yang baru)
	return s.repo.CreateUser(user)
}

func (s *authService) Login(req model.LoginRequest) (string, error) {
	// 1. Cari user
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		// Security Best Practice: Jangan pernah kasih tahu apakah 'username' atau 'password' yang salah.
		return "", errors.New("username atau password salah")
	}

	// 2. Cocokkan Password Asli dengan Hash di DB
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("username atau password salah")
	}

	// 3. Generate JWT Token menggunakan utils
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", errors.New("gagal membuat token autentikasi")
	}

	return token, nil
}