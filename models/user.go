package models

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `json:"name" validate:"required"`
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("User created")
	return
}
