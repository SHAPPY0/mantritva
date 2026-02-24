package auth

import (
	"github.com/gin-gonic/gin"
)

type AuthUser struct {
	ID        int
	Name      string
	UserName  string
	Email     string
	Role      string
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
		authUser.Found = true
	}
	return authUser
}

func (au *AuthUser) IsAdmin() bool {
	return au.Role == "ADMIN"
}

func (au *AuthUser) IsModerator() bool {
	return au.Role == "MODERATOR"
}

func (au *AuthUser) IsUser() bool {
	return au.Role == "USER"
}

func (au *AuthUser) IsAnonymous() bool {
	return au.Anonymous
}
