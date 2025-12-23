package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(*model.User) error
	GetByID(int) (*model.User, error)
	GetByUserID(string) (*model.User, error)
	ConsoleGetUsersList() (*[]model.User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) Create(user *model.User) error {
	if err := r.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetByUserID(userID string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("user_id = ? AND is_deleted = ?", userID, false).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ConsoleGetUsersList() (*[]model.User, error) {
	var users []model.User
	if err := r.DB.Order("created_at desc").Find(&users).Error; err != nil {
		return nil, err
	}
	for i := range users {
		users[i].Password = ""
	}
	return &users, nil
}

func (r *userRepository) GetByID(id int) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("id = ? AND is_deleted = ?", id, false).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
