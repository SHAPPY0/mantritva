package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
)

type UserHandler struct {
	Service service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		Service: svc,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.UserReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	reps, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create state", err)
		return
	}
	utils.Success(c, "User created successfully", reps)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req dto.LoginReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}

	user, err := h.Service.GetByUserID(req.UserID)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch user", err)
		return
	}
	if user == nil {
		utils.Error(c, 404, "User not found", nil)
		return
	}

	isValid := utils.ComparePassword(user.Password, req.Password)
	if !isValid {
		utils.Error(c, 401, "Invalid credentials", nil)
		return
	}

	user.Password = ""
	token, err := auth.CreateToken(auth.UserClaims{
		Name:         user.Name,
		UserName:     user.UserID,
		UserID:       user.ID,
		Email:        user.Email,
		Role:         user.Role,
		TokenVersion: user.TokenVersion,
	})
	if err != nil {
		log.Errorf("Error while creating token %v", err)
		utils.Error(c, 500, "Failed to login, Please try again.", err)
		return
	}
	utils.Success(c, "Login Successful", map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	username, err := auth.GetKeyInCtx[string](c, "user_name")
	if err != nil {
		utils.Error(c, 401, "No valid detail found", nil)
		return
	}
	user, err := h.Service.GetByUserID(username)
	if err != nil {
		utils.Error(c, 500, "Failed to get user profile", err)
		return
	}
	user.Password = ""
	utils.Success(c, "", user)
}
