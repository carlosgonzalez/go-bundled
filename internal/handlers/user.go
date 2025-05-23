package handlers

import (
	"net/http"

	"github.com/carlosgonzalez/go-bundled/internal/models"
	"github.com/carlosgonzalez/go-bundled/internal/repositories"

	"github.com/labstack/echo/v4"
)

type userHandler struct {
	repo repositories.UserRepositoryInterface
}

func NewUserHandler(repo repositories.UserRepositoryInterface) userHandler {
	return userHandler{
		repo: repo,
	}
}

func (uHandler *userHandler) CreateUser(c echo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unable to bind user, check the payload"})
	}
	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unable to validate user, check the payload"})
	}
	err := uHandler.repo.CreateUser(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "unable to create user"})
	}

	return c.JSON(http.StatusCreated, u)
}

func (uHandler *userHandler) GetUser(c echo.Context) error {
	user, err := uHandler.repo.GetUser(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}
	return c.JSON(http.StatusOK, user)
}

func (uHandler *userHandler) UpdateUser(c echo.Context) error {
	existingUser, err := uHandler.repo.GetUser(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	user := new(models.User)
	if err := c.Bind(&user); err != nil {
		return err
	}

	u, err := uHandler.repo.UpdateUser(&existingUser, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "unable to update user"})
	}
	return c.JSON(http.StatusOK, u)
}

func (uHandler *userHandler) DeleteUser(c echo.Context) error {
	existingUser, err := uHandler.repo.GetUser(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	err = uHandler.repo.DeleteUser(&existingUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "unable to delete user"})
	}
	return c.NoContent(http.StatusNoContent)
}

func (uHandler *userHandler) GetAllUsers(c echo.Context) error {
	users, err := uHandler.repo.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}
