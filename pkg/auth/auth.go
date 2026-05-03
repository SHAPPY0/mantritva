package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/pkg/utils"
)

type AuthUser struct {
	ID        int
	Name      string
	UserName  string
	Email     string
	Role      string
	Roles     []dto.RoleRes
	Found     bool
	Anonymous bool
}

func GetAuthUser(ctx *gin.Context) AuthUser {
	authUser := AuthUser{Found: false, Anonymous: false}
	anonymous, exist := ctx.Get("anonymous")
	if exist && anonymous.(bool) {
		authUser.Anonymous = true
		return authUser
	}
	user, exist := ctx.Get("user")
	if !exist {
		return authUser
	}
	userClaim, ok := user.(*UserClaims)
	if ok && userClaim != nil {
		authUser.ID = int(userClaim.UserID)
		authUser.Name = userClaim.Name
		authUser.UserName = userClaim.UserName
		authUser.Email = userClaim.Email
		authUser.Role = userClaim.Role
		authUser.Roles = userClaim.Roles
		authUser.Found = true
	}
	return authUser
}

func (au *AuthUser) IsSuperAdmin() bool {
	return au.HasRole(utils.SUPERADMIN)
}

func (au *AuthUser) HasRole(roleName utils.Roles) bool {
	if roleName == "" {
		return false
	}
	for _, role := range au.Roles {
		if role.Role == roleName {
			return true
		}
	}
	return false
}

func (au *AuthUser) IsAdmin() bool {
	return au.HasRole(utils.ADMIN)
}

func (au *AuthUser) IsAdminOrSuperAdmin() bool {
	return au.HasRole(utils.ADMIN) || au.HasRole(utils.SUPERADMIN)
}

func (au *AuthUser) IsModerator() bool {
	return au.Role == "MODERATOR"
}

func (au *AuthUser) IsUser() bool {
	return au.HasRole(utils.USER)
}

func (au *AuthUser) IsAnonymous() bool {
	return au.Anonymous
}
