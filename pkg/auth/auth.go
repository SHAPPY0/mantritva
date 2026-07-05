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

func (au *AuthUser) IsStateAdmin() bool {
	return au.HasRole(utils.STATE_ADMIN)
}
func (au *AuthUser) IsDistrictAdmin() bool {
	return au.HasRole(utils.DISTRICT_ADMIN)
}
func (au *AuthUser) IsLokSabhaAdmin() bool {
	return au.HasRole(utils.LOKSABHA_ADMIN)
}
func (au *AuthUser) IsVidhanSabhaAdmin() bool {
	return au.HasRole(utils.VIDHANSABHA_ADMIN)
}
func (au *AuthUser) IsLocalBodyAdmin() bool {
	return au.HasRole(utils.LOCALBODY_ADMIN)
}
func (au *AuthUser) IsWardAdmin() bool {
	return au.HasRole(utils.WARD_ADMIN)
}
func (au *AuthUser) IsAnyAdmin() bool {
	return au.IsAdminOrSuperAdmin() || au.IsStateAdmin() || au.IsDistrictAdmin() || au.IsLokSabhaAdmin() || au.IsVidhanSabhaAdmin() || au.IsLocalBodyAdmin() || au.IsWardAdmin()
}
func (au *AuthUser) IsUser() bool {
	return au.HasRole(utils.USER)
}
func (au *AuthUser) IsAnonymous() bool {
	return au.Anonymous
}

func (au *AuthUser) GetAssignedStates() []uint {
	var states []uint
	for _, role := range au.Roles {
		if role.Role == utils.STATE_ADMIN && role.Scope == utils.ROLE_SCOPE_STATE {
			states = append(states, role.ScopeID)
		}
	}
	return states
}

func (au *AuthUser) GetAssignedDistricts() []uint {
	var districts []uint
	for _, role := range au.Roles {
		if role.Role == utils.DISTRICT_ADMIN && role.Scope == utils.ROLE_SCOPE_DISTRICT {
			districts = append(districts, role.ScopeID)
		}
	}
	return districts
}

func (au *AuthUser) GetAssignedLokSabhas() []uint {
	var lokSabhas []uint
	for _, role := range au.Roles {
		if role.Role == utils.LOKSABHA_ADMIN && role.Scope == utils.ROLE_SCOPE_LOKSABHA {
			lokSabhas = append(lokSabhas, role.ScopeID)
		}
	}
	return lokSabhas
}

func (au *AuthUser) GetAssignedVidhanSabhas() []uint {
	var vidhanSabhas []uint
	for _, role := range au.Roles {
		if role.Role == utils.VIDHANSABHA_ADMIN && role.Scope == utils.ROLE_SCOPE_VIDHANSABHA {
			vidhanSabhas = append(vidhanSabhas, role.ScopeID)
		}
	}
	return vidhanSabhas
}

func (au *AuthUser) GetAssignedLocalBodies() []uint {
	var localBodies []uint
	for _, role := range au.Roles {
		if role.Role == utils.LOCALBODY_ADMIN && role.Scope == utils.ROLE_SCOPE_LOCALBODY {
			localBodies = append(localBodies, role.ScopeID)
		}
	}
	return localBodies
}

func (au *AuthUser) GetAssignedWards() []uint {
	var wards []uint
	for _, role := range au.Roles {
		if role.Role == utils.WARD_ADMIN && role.Scope == utils.ROLE_SCOPE_WARD {
			wards = append(wards, role.ScopeID)
		}
	}
	return wards
}
