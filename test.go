package main

import (
"fmt"
"go-fiber-pos/internal/config"
"go-fiber-pos/internal/core"
"go-fiber-pos/internal/infrastructure/provider"
"go-fiber-pos/internal/modules/inventory"
"go-fiber-pos/internal/modules/payment"

"github.com/go-playground/validator/v10"
"github.com/joho/godotenv"
)

func main() {
godotenv.Load()
config.ConnectDatabase()

db := config.DB

var order core.Order
if err := db.Order("created_at desc").First(&order).Error; err != nil {
tln("Query error:", err)

}

v := validator.New()
gateway := provider.NewMidtransAdapter()
invRepo := inventory.NewInventoryRepository(db)
invService := inventory.NewInventoryService(invRepo, v)

repo := payment.NewPaymentRepository(db)
service := payment.NewPaymentService(repo, gateway, invService, v)

req := payment.InitiatePaymentRequest{
      order.ID,
mentMethod: "QRIS",
}

resp, err := service.InitiatePayment(req)
if err != nil {
tf("Error InitiatePayment: %v\n", err)
} else {
tf("Success: %+v\n", resp)
}
}
