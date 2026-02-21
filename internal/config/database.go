package config

import (
	"fmt"
	"go-fiber-pos/internal/core"
	"go-fiber-pos/pkg/logger"

	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbname, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	DB = database
	logger.Log.Info("Berhasil terhubung ke database PostgreSQL via GORM!")

	err = database.AutoMigrate(
		&core.StoreProfile{},
		&core.User{},
		&core.Category{},
		&core.Product{},
		&core.Voucher{},
		&core.DailyCounter{},
		&core.Order{},
		&core.OrderItem{},
		&core.Payment{},
	)

	if err != nil {
		logger.Log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}
	
	logger.Log.Info("Migrasi database berhasil! Tabel sudah siap digunakan.")
	// ----------------------------------------------

	DB = database
}