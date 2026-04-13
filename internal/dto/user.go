package dto

import (
	"time"

	"github.com/shappy0/mantritva/pkg/utils"
)

type UserReq struct {
	Name            string `json:"name" binding:"required"`
	UserID          string `json:"user_id" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Mobile          string `json:"mobile,omitempty"`
	Password        string `json:"password" binding:"required"`
	AccountVerified bool   `json:"account_verified" binding:"required"`
	Role            string `json:"role" binding:"required"`
}

type LoginReq struct {
	UserID   string `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserClaim struct {
	Name   string         `json:"name"`
	UserID string         `json:"user_id"`
	Email  string         `json:"email"`
	Role   utils.UserRole `json:"role"`
}

type UpdateUserReq struct {
	ID        int `json:"id"`
	Name      string
	Email     string
	Mobile    string
	Password  string
	Role      utils.UserRole
	UpdatedAt time.Time
	UpdatedBy string
}

type SendOtpReq struct {
	Identifier string           `json:"identifier" binding:"required"`
	Purpose    utils.OTPPurpose `json:"purpose" binding:"required"`
	Mode       utils.OTPMode    `json:"mode" binding:"required,oneof=EMAIL SMS"`
}

type VerifyOtpReq struct {
	Identifier string           `json:"identifier" binding:"required"`
	Otp        string           `json:"otp" binding:"required"`
	Mode       utils.OTPMode    `json:"mode" binding:"required,oneof=EMAIL SMS"`
	Purpose    utils.OTPPurpose `json:"purpose"`
}
