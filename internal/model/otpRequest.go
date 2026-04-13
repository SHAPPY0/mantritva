package model

import (
	"time"

	"github.com/shappy0/mantritva/pkg/utils"
)

type OtpRequest struct {
	ID          uint             `json:"_id,omitempty"`
	Identifier  string           `json:"identifier"`
	Mode        utils.OTPMode    `json:"mode"`
	OtpHash     string           `json:"otp_hash"`
	Purpose     utils.OTPPurpose `json:"purpose"`
	ExpiresAt   time.Time        `json:"expires_at"`
	ConsumedAt  *time.Time       `json:"consumed_at,omitempty"`
	Attempts    uint             `json:"attempts"`
	MaxAttempts uint             `json:"max_attempts" gorm:"default:5"`
	IpAddress   string           `json:"ip_address"`
	UserAgent   string           `json:"user_agent"`
	AuditFields `json:",inline"`
}
