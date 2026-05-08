package handler

import (
	"fmt"
	"regexp"

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

func validatePassword(password string) error {
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one number")
	}
	if !regexp.MustCompile(`[@$!%*?&]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one special character")
	}
	return nil
}

func validateUserReq(req dto.UserReq) error {
	if req.Name == "" || req.UserID == "" || req.Email == "" || req.Password == "" || req.Mobile == "" {
		return fmt.Errorf("Please fill required details")
	}
	if len(req.Mobile) != 10 || !regexp.MustCompile(`^\d+$`).MatchString(req.Mobile) {
		return fmt.Errorf("Invalid mobile number format")
	}
	if err := validatePassword(req.Password); err != nil {
		return fmt.Errorf("Password must be at least 6-20 characters long and include at least one uppercase letter, one lowercase letter, one digit, and one special character")
	}
	return nil
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.UserReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	if err := validateUserReq(req); err != nil {
		utils.Error(c, 400, err.Error(), nil)
		return
	}
	if !req.AccountVerified {
		utils.Error(c, 500, "Account is not verified, Please verifiy it.", nil)
		return
	}
	reps, err := h.Service.Create(req)
	if err != nil {
		log.Errorf("Error while creating user: %v", err)
		if err.Error() == "User already exist with given details" {
			utils.Error(c, 400, err.Error(), nil)
			return
		}
		utils.Error(c, 500, "Error while creating user", err)
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
		log.Errorf("Failed to fetch user: %v", err)
		utils.Error(c, 500, "Failed to fetch user", err)
		return
	}
	if user == nil {
		log.Errorf("User not found: %s", req.UserID)
		utils.Error(c, 404, "User not found", nil)
		return
	}
	if !user.AccountVerified {
		log.Errorf("Unverified account is loggin in: %s", req.UserID)
		utils.Error(c, 401, "Account is not verified, Please reach out to admin", nil)
		return
	}
	isValid := utils.ComparePassword(*user.Password, req.Password)
	if !isValid {
		log.Errorf("Invalid credentials for user %s", user.UserID)
		utils.Error(c, 401, "Invalid credentials", nil)
		return
	}

	user.Password = nil
	userClaim := auth.UserClaims{
		Name:         user.Name,
		UserName:     user.UserID,
		UserID:       user.ID,
		Email:        user.Email,
		Roles:        user.Roles,
		TokenVersion: user.TokenVersion,
	}
	token, err := auth.CreateToken(userClaim)
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
	user.Password = nil
	utils.Success(c, "", user)
}

func (h *UserHandler) SendOTP(c *gin.Context) {
	var req []dto.SendOtpReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	if len(req) == 0 {
		utils.Error(c, 400, "Request body cannot be empty", nil)
		return
	}
	userIP := utils.ClientIP(c)
	id, err := h.Service.SendOTP(req, userIP)
	if err != nil {
		utils.Error(c, 500, "Failed to send OTP", err)
		return
	}
	utils.Success(c, "OTP sent successfully", map[string]interface{}{
		"otp_id": id,
	})
}

func (h *UserHandler) VerifyOTP(c *gin.Context) {
	var req dto.VerifyOtpReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	valid, err := h.Service.VerifyOTP(req)
	if err != nil {
		utils.Error(c, 500, "Error verifying OTP", err)
		return
	}
	if valid {
		utils.Success(c, "OTP verified successfully", nil)
	} else {
		utils.Error(c, 500, "OTP Expired", nil)
	}
}

func (h *UserHandler) ListRoles(c *gin.Context) {
	roles, err := h.Service.ListRoles()
	if err != nil {
		utils.Error(c, 500, "Error fetching roles", err)
		return
	}
	utils.Success(c, "", roles)
}
