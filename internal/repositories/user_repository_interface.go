package repositories

import (
	"github.com/carlosgonzalez/go-bundled/internal/models"
)

type UserRepositoryInterface interface {
	CreateUser(user *models.User) error
	DeleteUser(user *models.User) error
	GetAllUsers() (error, []*models.User)
	GetUser(id string) (error, models.User)
	UpdateUser(oldUser *models.User, newUser *models.User) (error, *models.User)
}
