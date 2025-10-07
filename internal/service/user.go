package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/utils"
)

type UserService interface {
	Create(dto.UserReq) (*model.User, error)
	GetByUserID(string) (*model.User, error)
}

type userService struct {
	Repository repository.UserRepository
}

func NewUserService(repos *repository.Repos) UserService {
	return &userService{
		Repository: repos.User,
	}
}

func (s *userService) Create(req dto.UserReq) (*model.User, error) {
	user := &model.User{
		Name:     req.Name,
		UserID:   req.UserID,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: req.Password,
		Role:     req.Role,
	}

	hashedPassword, err := utils.HashAndSaltPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	if err := s.Repository.Create(user); err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *userService) GetByUserID(userID string) (*model.User, error) {
	user, err := s.Repository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
