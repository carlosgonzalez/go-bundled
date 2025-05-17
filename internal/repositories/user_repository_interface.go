package repositories

import (
	"github.com/carlosgonzalez/go-bundled/internal/models"
)

type UserRepositoryInterface interface {
	CreateUser(user *models.User) error
	DeleteUser(user *models.User) error
	GetAllUsers() ([]*models.User, error)
	GetUser(id string) (models.User, error)
	UpdateUser(oldUser *models.User, newUser *models.User) (*models.User, error)
}
