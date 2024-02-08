package handlers

import (
	"net/http"

	"github.com/carlosgonzalez/learning-go/internal/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type userHandler struct {
	db *gorm.DB
	models.User
}

func NewUserHandler(db *gorm.DB) *userHandler {
	return &userHandler{
		db:   db,
		User: models.User{},
	}
}

func (uHandler *userHandler) CreateUser(c echo.Context) error {
	u := &models.User{}

	if err := c.Validate(u); err != nil {
		return err
	}

	if err := c.Bind(u); err != nil {
		return err
	}

	uHandler.db.Create(u)

	return c.JSON(http.StatusCreated, u)
}

func (uHandler *userHandler) GetUser(c echo.Context) error {
	var user models.User
	if err := uHandler.db.First(&user, c.Param("id")).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (uHandler *userHandler) UpdateUser(c echo.Context) error {
	var existingUser models.User
	if err := uHandler.db.First(&existingUser, c.Param("id")).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	user := new(models.User)
	if err := c.Bind(&user); err != nil {
		return err
	}

	uHandler.db.Model(&existingUser).Updates(models.User{Name: user.Name})
	return c.JSON(http.StatusOK, existingUser)
}

func (uHandler *userHandler) DeleteUser(c echo.Context) error {
	var existingUser models.User
	if err := uHandler.db.First(&existingUser, c.Param("id")).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	uHandler.db.Delete(&existingUser)

	return c.NoContent(http.StatusNoContent)
}

func (uHandler *userHandler) GetAllUsers(c echo.Context) error {

	users := []*models.User{}
	uHandler.db.Find(&users)

	return c.JSON(http.StatusOK, users)
}
