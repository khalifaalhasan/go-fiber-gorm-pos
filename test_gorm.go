package main

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func main() {
	db := &gorm.DB{Error: errors.New("dummy error")}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panicked!", r)
		}
	}()

	tx := db.Begin()
	fmt.Println("Result:", tx.Error)
}
