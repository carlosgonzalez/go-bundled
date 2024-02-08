package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title string `json:"title" validate:"required"`
	Body  string `json:"body" validate:"required"`
}
