package users

import (
	"api/src/modules/users/entities"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entities.User) error
	FindAll() ([]entities.User, error)
	FindById(id uint) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	Update(user *entities.User) error
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindAll() ([]entities.User, error) {
	var users []entities.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) FindById(id uint) (*entities.User, error) {
	var user entities.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) Update(user *entities.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&entities.User{}, id).Error
}
