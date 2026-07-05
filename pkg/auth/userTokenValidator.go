package auth

import "github.com/shappy0/mantritva/internal/model"

type UserAuthorization interface {
	ValidateTokenVersion(int, int) (bool, error)
	GetUserRoles(uint) ([]model.UserRole, error)
}
