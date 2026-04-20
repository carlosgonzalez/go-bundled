package handlers

import (
	"net/http"

	"github.com/carlosgonzalez/go-bundled/internal/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type postHandler struct {
	db *gorm.DB
}

func NewPostHandler(db *gorm.DB) *postHandler {
	return &postHandler{
		db: db,
	}
}

func (pHandler *postHandler) CreatePost(c echo.Context) error {
	posts, err := services.Fetcher("https://jsonplaceholder.typicode.com", "posts", 50)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "unable to fetch posts"})
	}

	pHandler.db.Create(posts)

	return c.JSON(http.StatusOK, posts)
}
