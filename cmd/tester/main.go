package main

import (
	"fmt"
	"go-fiber-pos/internal/core"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=go_pos_user password=go_pos_password dbname=go_pos_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("DB connection error:", err)
		os.Exit(1)
	}

	var order core.Order
	if err := db.Order("created_at desc").First(&order).Error; err != nil {
		fmt.Println("Query error:", err)
		os.Exit(1)
	}

	fmt.Println("OrderID:", order.ID)
}
