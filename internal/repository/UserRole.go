package repository

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
	"gorm.io/gorm"
)

type UserRoleRepository interface {
	IsRoleExist(userID uint, roleID uint, scope utils.RoleScope, scopeID *uint) (bool, error)
	AddRole(string, dto.AddRoleReq) (*model.UserRole, error)
	GetUserRoles(string) ([]model.UserRole, error)
	GetRoleByUserIDs([]uint) ([]model.UserRole, error)
}

type userRoleRepository struct {
	DB *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) UserRoleRepository {
	return &userRoleRepository{
		DB: db,
	}
}

func (r *userRoleRepository) IsRoleExist(userID uint, roleID uint, scope utils.RoleScope, scopeID *uint) (bool, error) {
	var userRole model.UserRole
	query := r.DB.Where("user_id = ? AND role_id = ? AND scope = ? AND is_deleted = ?", userID, roleID, scope, false)
	if scopeID != nil {
		query = query.Where("scope_id = ?", *scopeID)
	}
	if err := query.First(&userRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *userRoleRepository) AddRole(createdBy string, req dto.AddRoleReq) (*model.UserRole, error) {
	role := &model.UserRole{
		UserID:  req.UserID,
		RoleID:  req.RoleID,
		Role:    req.Role,
		Scope:   req.Scope,
		ScopeID: req.ScopeID,
		AuditFields: model.AuditFields{
			CreatedBy: createdBy,
		},
	}
	if err := r.DB.Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *userRoleRepository) GetUserRoles(userID string) ([]model.UserRole, error) {
	var roles []model.UserRole
	if err := r.DB.Where("user_id = ? AND is_deleted = ?", userID, false).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *userRoleRepository) GetRoleByUserIDs(userIDs []uint) ([]model.UserRole, error) {
	var roles []model.UserRole
	if err := r.DB.Where("user_id IN ? AND is_deleted = ?", userIDs, false).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
