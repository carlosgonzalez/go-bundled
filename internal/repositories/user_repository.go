package repositories

import (
	"github.com/carlosgonzalez/go-bundled/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{
		db: db,
	}
}

func (uRepo UserRepository) GetUser(id string) (error, models.User) {
	var user models.User
	if err := uRepo.db.First(&user, id).Error; err != nil {
		return err, user
	}
	return nil, user

}

func (uRepo UserRepository) GetAllUsers() (error, []*models.User) {
	users := []*models.User{}
	tx := uRepo.db.Find(&users)
	if tx.Error != nil {
		return tx.Error, nil
	}

	return nil, users
}

func (uRepo UserRepository) CreateUser(user *models.User) error {
	tx := uRepo.db.Create(user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (uRepo UserRepository) UpdateUser(oldUser *models.User, newUser *models.User) (error, *models.User) {
	tx := uRepo.db.Model(&oldUser).Updates(models.User{Name: newUser.Name})
	if tx.Error != nil {
		return tx.Error, nil
	}
	return nil, oldUser
}

func (uRepo UserRepository) DeleteUser(user *models.User) error {
	tx := uRepo.db.Delete(&user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
