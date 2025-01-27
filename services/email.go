package services

import (
	"api/models"
	"api/repositories"
	"api/structs"
	"api/utils"
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type EmailService struct {
	userRepository               *repositories.UserRepository
	emailVerifyTokenRepository   *repositories.EmailVerifyTokenRepository
	passwordResetTokenRepository *repositories.PasswordResetTokenRepository
	emailQueueRepository         *repositories.EmailQueueRepository
}

func NewEmailService() *EmailService {
	return &EmailService{
		userRepository:               repositories.NewUserRepository(),
		emailVerifyTokenRepository:   repositories.NewEmailVerifyTokenRepository(),
		passwordResetTokenRepository: repositories.NewPasswordResetTokenRepository(),
		emailQueueRepository:         repositories.NewEmailQueueRepository(),
	}
}

func (s *EmailService) SendResetPasswordEmail(ctx *gin.Context, email string) {
	tokens := []models.PasswordResetTokenModel{}
	s.passwordResetTokenRepository.GetByEmail(ctx, email, &tokens)
	isSpam := true

	if len(tokens) > 0 {
		if len(tokens) >= 5 {
			// Check if the 5 most recent tokens are sent within the last 1 hour
			for i := 0; i < 5; i++ {
				if tokens[i].CreatedAt.Add(1 * time.Hour).After(time.Now()) {
					isSpam = true
					break
				}
			}
		} else {
			isSpam = false
		}
	} else {
		isSpam = false
	}

	if isSpam {
		return
	}

	// Get the current working directory
	cwd, _ := os.Getwd()
	passwordResetTemplate := filepath.Join(cwd, "mails", "password_reset.html")
	template, err := template.ParseFiles(passwordResetTemplate)
	if err != nil {
		return
	}

	var user models.UserModel
	s.userRepository.GetByEmail(ctx, &user, email)

	tokenString, err := utils.GenerateString(64)

	if err != nil {
		return
	}

	hasedToken, err := utils.HashString(tokenString)

	if err != nil {
		return
	}

	var body bytes.Buffer
	data := structs.ResetPasswordTemplateData{
		Name:              user.LastName,
		ResetPasswordLink: ctx.Request.Host + "/new-password?email=" + email + "&token=" + tokenString,
	}

	err = template.Execute(&body, data)

	if err != nil {
		return
	}

	// // Send email
	// message := gomail.NewMessage()
	// message.SetHeader("From", config.MailFromAddress, config.MailFromName)
	// message.SetHeader("To", email)
	// message.SetHeader("Subject", "Đổi mật khẩu - Reset your password")
	// message.SetBody("text/html", body.String())

	// if err := config.Mailer.DialAndSend(message); err == nil {
	// 	// Log error

	// 	return
	// }

	if err := s.emailQueueRepository.Create(&models.EmailQueueModel{
		ReceiverEmail: email,
		Subject:       "Đổi mật khẩu - Reset your password",
		Body:          body.String(),
	}); err != nil {
		return
	}

	s.passwordResetTokenRepository.Create(ctx, &models.PasswordResetTokenModel{
		Email: email,
		Token: hasedToken,
	})
}

func (s *EmailService) SendEmailVerificationEmail(ctx *gin.Context, email string) {
	tokens := []models.EmailVerifyTokenModel{}
	s.emailVerifyTokenRepository.GetByEmail(ctx, email, &tokens)
	isSpam := true

	if len(tokens) > 0 {
		if len(tokens) >= 5 {
			// Check if the 5 most recent tokens are sent within the last 1 hour
			for i := 0; i < 5; i++ {
				if tokens[i].CreatedAt.Add(1 * time.Hour).After(time.Now()) {
					isSpam = true
					break
				}
			}
		} else {
			isSpam = false
		}
	} else {
		isSpam = false
	}

	if isSpam {
		return
	}

	// Get the current working directory
	cwd, _ := os.Getwd()
	emailVerificationTemplate := filepath.Join(cwd, "mails", "email_verification.html")
	template, err := template.ParseFiles(emailVerificationTemplate)
	if err != nil {
		return
	}

	var user models.UserModel
	s.userRepository.GetByEmail(ctx, &user, email)

	tokenString, err := utils.GenerateString(64)

	if err != nil {
		return
	}

	hasedToken, err := utils.HashString(tokenString)

	if err != nil {
		return
	}

	var body bytes.Buffer
	data := structs.VerificationTemplateData{
		Name:             user.LastName,
		VerificationLink: ctx.Request.Host + "/verify-email?email=" + email + "&token=" + tokenString,
	}

	err = template.Execute(&body, data)

	if err != nil {
		return
	}

	// // Send email
	// message := gomail.NewMessage()
	// message.SetHeader("From", config.MailFromAddress, config.MailFromName)
	// message.SetHeader("To", email)
	// message.SetHeader("Subject", "Xác thực email - Verify your email")
	// message.SetBody("text/html", body.String())

	// if err := config.Mailer.DialAndSend(message); err == nil {
	// 	// Log error

	// 	return
	// }

	if err := s.emailQueueRepository.Create(&models.EmailQueueModel{
		ReceiverEmail: email,
		Subject:       "Xác thực email - Verify your email",
		Body:          body.String(),
	}); err != nil {
		return
	}

	s.emailVerifyTokenRepository.Create(ctx, &models.EmailVerifyTokenModel{
		Email: email,
		Token: hasedToken,
	})
}
