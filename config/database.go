package config

import (
	"fmt"
	"go-fiber-pos/model"
	"go-fiber-pos/utils"
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
		utils.Log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	DB = database
	utils.Log.Info("Berhasil terhubung ke database PostgreSQL via GORM!")

	err = database.AutoMigrate(
		&model.StoreProfile{},
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.Voucher{},
		&model.Order{},
		&model.OrderItem{},
		&model.Payment{},
	)

	if err != nil {
		utils.Log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}
	
	utils.Log.Info("Migrasi database berhasil! Tabel sudah siap digunakan.")
	// ----------------------------------------------

	DB = database
}