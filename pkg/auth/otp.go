package auth

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/config"
	emailTemplates "github.com/shappy0/mantritva/pkg/config/emailTemplates"
	"github.com/shappy0/mantritva/pkg/log"
)

type OtpHandler interface {
	SendOTP(model.OtpRequest, string) error
}

type SMSOtpHandler struct {
}

type EmailOtpHandler struct {
	OtpTemplate string
}

func NewSMSOtpHandler() *SMSOtpHandler {
	return &SMSOtpHandler{}
}

func NewEmailOtpHandler() *EmailOtpHandler {
	eoh := &EmailOtpHandler{
		OtpTemplate: emailTemplates.OTP_TEMPLATE,
	}
	return eoh
}

func (h *EmailOtpHandler) SendOTP(otpRequest model.OtpRequest, otp string) error {
	from := config.Env.GetString("SMTP_FROM")
	to := []string{otpRequest.Identifier}
	subject := "One time password for Mantritva"

	smtpHost := config.Env.GetString("SMTP_HOST")
	smtpPort := fmt.Sprintf("%d", config.Env.GetInt("SMTP_PORT"))
	log.Infof("From %s, host %s, port %s", from, smtpHost, smtpPort)
	msg := fmt.Sprintf("From: %s\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", otpRequest.Identifier)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	msg += fmt.Sprintf("Subject: %s\r\n\r\n", subject)
	// Use strings.Replace to avoid fmt parsing of other % characters in the HTML template
	msg += strings.Replace(h.OtpTemplate, "{{otp}}", otp, 1)
	auth := smtp.PlainAuth("", config.Env.GetString("SMTP_USERNAME"), config.Env.GetString("SMTP_PASSWORD"), smtpHost)
	log.Infof("Body: %s, username %s, password %s", msg, config.Env.GetString("SMTP_USERNAME"), config.Env.GetString("SMTP_PASSWORD"))
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
	return err
}

func (h *SMSOtpHandler) SendOTP(otpRequest model.OtpRequest, otp string) error {
	// Implement SMS sending logic here, e.g., using Twilio or any other SMS service provider
	return nil
}
