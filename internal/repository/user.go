package repository

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(*model.User) error
	GetByID(int) (*model.User, error)
	GetByUserID(string) (*model.User, error)
	ConsoleGetUsersList() (*[]model.User, error)
	UpdateTokenVersion(int, int) (bool, error)
	GetTokenVersion(int) (int, error)
	UpdateRole(dto.UpdateUserReq) (bool, error)
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
		users[i].Password = nil
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

func (r *userRepository) UpdateTokenVersion(userId, version int) (bool, error) {
	query := `UPDATE users SET token_version = ? WHERE id = ? AND is_deleted = false`
	if err := r.DB.Exec(query, version, userId).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (r *userRepository) GetTokenVersion(userId int) (int, error) {
	version := 0
	query := `SELECT token_version FROM users WHERE id = ? AND is_deleted = false`
	if err := r.DB.Raw(query, userId).Scan(&version).Error; err != nil {
		return version, err
	}
	return version, nil
}

func (r *userRepository) UpdateRole(user dto.UpdateUserReq) (bool, error) {
	query := `UPDATE users SET role = ?, updated_at = ?, updated_by = ? WHERE id = ? AND is_deleted = false`
	if err := r.DB.Exec(query, string(user.Role), user.UpdatedAt, user.UpdatedBy, user.ID).Error; err != nil {
		return false, err
	}
	return true, nil
}
