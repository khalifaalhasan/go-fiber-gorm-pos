package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func auth() {
	godotenv.Load()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		fmt.Println("No JWT_SECRET found")
		os.Exit(1)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
		"username": "admin",
		"role":     "ADMIN",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("Error signing token:", err)
		os.Exit(1)
	}
	fmt.Println(tokenString)
}
