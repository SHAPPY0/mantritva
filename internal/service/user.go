package service

import (
	"time"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/auth"
	"github.com/shappy0/mantritva/pkg/config"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
)

type UserService interface {
	Create(dto.UserReq) (*model.User, error)
	GetByID(int) (*model.User, error)
	GetByUserID(string) (*model.User, error)
	UpdateTokenVersion(int, int) (bool, error)
	ValidateTokenVersion(int, int) (bool, error)
	SendOTP([]dto.SendOtpReq, string) (uint, error)
	VerifyOTP(dto.VerifyOtpReq) (bool, error)
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
		Name:            req.Name,
		UserID:          req.UserID,
		Email:           req.Email,
		Mobile:          req.Mobile,
		Password:        &req.Password,
		Role:            req.Role,
		AccountVerified: req.AccountVerified,
	}

	if userExist, err := s.Repository.UserExist(user.Mobile, user.Email); err != nil {
		log.Errorf("Error while checking if user exist %v", err)
		return nil, utils.NewError("Internal server error")
	} else if userExist != nil {
		log.Errorf("User already exist with mobile %s, email %s", userExist.Mobile, user.Email)
		return nil, utils.NewError("User already exist with given details")
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

func (s *userService) SendOTP(req []dto.SendOtpReq, ipAddress string) (uint, error) {
	for _, r := range req {
		now := time.Now()
		otp := utils.GenrateOTP(4)
		otpRequest := model.OtpRequest{
			Identifier: r.Identifier,
			Purpose:    r.Purpose,
			Mode:       r.Mode,
			OtpHash:    utils.HashOTP(otp, config.Env.GetString("OTP_SECRET")),
			ExpiresAt:  now.Add(time.Minute * 5),
			IpAddress:  ipAddress,
			AuditFields: model.AuditFields{
				CreatedAt: now,
				CreatedBy: "SYSTEM",
			},
		}
		var otpHandler auth.OtpHandler
		switch r.Mode {
		case utils.OTP_MODE_EMAIL:
			if !utils.IsValidEmail(r.Identifier) {
				return 0, utils.NewError("Invalid email address")
			}
			if err := s.Repository.SaveOTPRequest(&otpRequest); err != nil {
				return 0, err
			}
			otpHandler = auth.NewEmailOtpHandler()
		case utils.OTP_MODE_SMS:
			if !utils.IsValidMobile(r.Identifier) {
				return 0, utils.NewError("Invalid mobile number")
			}
			if err := s.Repository.SaveOTPRequest(&otpRequest); err != nil {
				return 0, err
			}
			otpHandler = auth.NewSMSOtpHandler()
		case utils.OTP_MODE_WHATSAPP:
			if !utils.IsValidMobile(r.Identifier) {
				return 0, utils.NewError("Invalid mobile number")
			}
			if err := s.Repository.SaveOTPRequest(&otpRequest); err != nil {
				return 0, err
			}
			otpHandler = auth.NewWhatsappOtpHandler()
		default:
			return 0, utils.NewError("Invalid OTP mode")
		}

		if otpHandler == nil {
			return 0, utils.NewError("Failed to create OTP handler")
		}

		if err := otpHandler.SendOTP(otpRequest, otp); err != nil {
			log.Errorf("Failed to send OTP: %v", err)
			s.deleteOtpRequest(otpRequest.ID)
			return 0, err
		} else {
			log.Infof("OTP sent successfully to %s for purpose %s", r.Identifier, r.Purpose)
			return otpRequest.ID, nil
		}
	}
	return 0, nil
}

func (s *userService) deleteOtpRequest(otpID uint) error {
	if err := s.Repository.DeleteOTPRequest(otpID); err != nil {
		return err
	}
	return nil
}

func (s *userService) VerifyOTP(req dto.VerifyOtpReq) (bool, error) {
	if otpRequest, err := s.Repository.GetOTPRequest(req.Identifier, req.Mode, req.Purpose); err != nil {
		return false, utils.NewError("Invalid otp verification request")
	} else if otpRequest != nil {
		if otpRequest.ConsumedAt != nil {
			return false, utils.NewError("OTP Expired")
		}
		if time.Now().After(otpRequest.ExpiresAt) {
			return false, utils.NewError("OTP Expired")
		}
		if otpRequest.Attempts >= otpRequest.MaxAttempts {
			return false, utils.NewError("OTP Expired")
		}
		hashedOtp := utils.HashOTP(req.Otp, config.Env.GetString("OTP_SECRET"))
		otpRequest.Attempts++
		if !utils.SecureCompare(hashedOtp, otpRequest.OtpHash) {
			if _, err := s.Repository.UpdateOTPAttempts(otpRequest.ID, otpRequest.Attempts); err != nil {
				log.Errorf("Error while updating otp attempt %v", err)
			}
			return false, utils.NewError("Incorrect OTP")
		}
		now := time.Now()
		if ok, err := s.Repository.UpdateOTPConsumed(otpRequest.ID, otpRequest.Attempts, now); err != nil {
			log.Errorf("Error while updating otp request, but otp matched")
		} else if !ok {
			log.Errorf("Error while updating otp request, but otp matched")
		}
		return true, nil
	} else {
		return false, utils.NewError("Invalid otp verification request")
	}
}
