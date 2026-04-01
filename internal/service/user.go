package service

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/auth"
	"github.com/shappy0/mantritva/pkg/utils"
)

type UserService interface {
	Create(dto.UserReq) (*model.User, error)
	GetByID(int) (*model.User, error)
	GetByUserID(string) (*model.User, error)
	UpdateTokenVersion(int, int) (bool, error)
	ValidateTokenVersion(int, int) (bool, error)
}

type userService struct {
	Repository repository.UserRepository
}

func NewUserService(repos *repository.Repos) UserService {
	return &userService{
		Repository: repos.User,
	}
}

var _ auth.TokenVersionValidator = (*userService)(nil)

func (s *userService) Create(req dto.UserReq) (*model.User, error) {
	user := &model.User{
		Name:     req.Name,
		UserID:   req.UserID,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: &req.Password,
		Role:     req.Role,
	}

	hashedPassword, err := utils.HashAndSaltPassword(*user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = &hashedPassword

	if err := s.Repository.Create(user); err != nil {
		return nil, err
	}
	user.Password = nil
	return user, nil
}

func (s *userService) GetByID(id int) (*model.User, error) {
	user, err := s.Repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetByUserID(userID string) (*model.User, error) {
	user, err := s.Repository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateTokenVersion(userId, version int) (bool, error) {
	ok, err := s.Repository.UpdateTokenVersion(userId, version)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (s *userService) GetTokenVersion(userId int) (int, error) {
	version, err := s.Repository.GetTokenVersion(userId)
	if err != nil {
		return version, err
	}
	return version, nil
}

func (s *userService) ValidateTokenVersion(userId, jwtTokenVersion int) (bool, error) {
	dbVersion, err := s.GetTokenVersion(userId)
	if err != nil {
		return false, err
	}
	if dbVersion != jwtTokenVersion {
		return false, nil
	}
	return true, nil
}
