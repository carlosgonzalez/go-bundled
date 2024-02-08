package handlers

import (
	"net/http"

	"github.com/carlosgonzalez/learning-go/models"
	"github.com/carlosgonzalez/learning-go/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type postHandler struct {
	db *gorm.DB
	models.Post
}

func NewPostHandler(db *gorm.DB) *postHandler {
	return &postHandler{
		db:   db,
		Post: models.Post{},
	}
}

func (pHandler *postHandler) CreatePost(c echo.Context) error {
	posts, err := services.Fetcher("posts", 50)
	if err != nil {
		return err
	}

	for _, post := range posts {
		pHandler.db.Create(&post)
	}

	return c.JSON(http.StatusOK, posts)
}
