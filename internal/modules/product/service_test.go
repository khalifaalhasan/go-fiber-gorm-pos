package product_test

import (
	"errors"
	"testing"

	"go-fiber-pos/internal/core"

	"go-fiber-pos/internal/modules/product"
	"go-fiber-pos/internal/modules/product/mocks"

	"github.com/google/uuid"

	"github.com/go-playground/validator/v10"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateProduct_Gomock(t *testing.T) {
	testCases := []struct {
		name          string
		req           product.CreateProductRequest
		buildStubs    func(mockRepo *mocks.MockProductRepository)
		expectedError error
	}{
		{
			name: "Sukses Membuat produk Baru",
			
			req: product.CreateProductRequest{
				CategoryID:  uuid.New(),
				Name:        "Americano",
				Description: "Kopi americano yang nikmat",
				NormalPrice: 15000,
			},
			buildStubs: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().
					FindByName("Americano").
					Return(nil, nil).
					Times(1)

				mockRepo.EXPECT().
					Create(gomock.Any()).
					Return(nil).
					Times(1)
			},
			expectedError: nil,
		},
		{
			name: "Gagal - Nama Produk Sudah Terdaftar",
			req: product.CreateProductRequest{
				CategoryID:  uuid.New(),
				Name:        "Americano",
				Description: "Kopi americano yang nikmat",
				NormalPrice: 15000,
			},
			buildStubs: func(mockRepo *mocks.MockProductRepository) {
				existingProduct := &core.Product{Name: "Americano"}

				mockRepo.EXPECT().
					FindByName("Americano").
					Return(existingProduct, nil).
					Times(1)
			},
			
			
			expectedError: errors.New("produk sudah ada"),
		},
		{
			name: "Gagal - Validasi Input Tidak Lengkap",
			req: product.CreateProductRequest{
				Name: "Es", // min=3, ini cuma 2 karakter
			},
			buildStubs: func(mockRepo *mocks.MockProductRepository) {
				// Tidak ada mock call karena validasi gagal duluan
			},
			expectedError: nil, // kita cek pakai assert.Error saja
		},
		{
			name: "Gagal - Error Database Saat FindByName",
			req: product.CreateProductRequest{
				CategoryID:  uuid.New(),
				Name:        "Americano",
				Description: "Kopi americano yang nikmat",
				NormalPrice: 15000,
			},
			buildStubs: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().
					FindByName("Americano").
					Return(nil, errors.New("db connection lost")).
					Times(1)
			},
			expectedError: errors.New("terjadi kesalahan pada server"),
		},
		{
			name: "Gagal - Error Database Saat Create",
			req: product.CreateProductRequest{
				CategoryID:  uuid.New(),
				Name:        "Americano",
				Description: "Kopi americano yang nikmat",
				NormalPrice: 15000,
			},
			buildStubs: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().
					FindByName("Americano").
					Return(nil, nil).
					Times(1)

				mockRepo.EXPECT().
					Create(gomock.Any()).
					Return(errors.New("gagal menyimpan data")).
					Times(1)
			},
			expectedError: errors.New("gagal menyimpan data"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockProductRepository(ctrl)

			tc.buildStubs(mockRepo)

			v := validator.New()
			service := product.NewProductService(mockRepo, v)

			_, err := service.CreateProduct(tc.req)

			// Case khusus: validasi gagal (cek error ada, tapi bukan dari mock)
			if tc.name == "Gagal - Validasi Input Tidak Lengkap" {
				assert.Error(t, err)
				return
			}

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAllProducts_Gomock(t *testing.T) {
	testCases := []struct {
		name           string
		buildStubs     func(mockRepo *mocks.MockProductRepository)
		expectedCount  int
		expectedError  error
	}{
		{
			name: "Sukses Mengambil Semua Produk",
			buildStubs: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().
					GetAll().
					Return([]core.Product{
						{Name: "Americano", NormalPrice: 15000},
						{Name: "Latte", NormalPrice: 20000},
					}, nil).
					Times(1)
			},
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name: "Sukses - Tidak Ada Produk (List Kosong)",
			buildStubs: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().
					GetAll().
					Return([]core.Product{}, nil).
					Times(1)
			},
			expectedCount: 0,
			expectedError: nil,
		},
		{
			name: "Gagal - Error Database Saat GetAll",
			buildStubs: func(mockRepo *mocks.MockProductRepository) {
				mockRepo.EXPECT().
					GetAll().
					Return(nil, errors.New("db connection lost")).
					Times(1)
			},
			expectedCount: 0,
			expectedError: errors.New("db connection lost"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockProductRepository(ctrl)

			tc.buildStubs(mockRepo)

			v := validator.New()
			service := product.NewProductService(mockRepo, v)

			products, err := service.GetAllProducts()

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, products)
			} else {
				assert.NoError(t, err)
				assert.Len(t, products, tc.expectedCount)
			}
		})
	}
}