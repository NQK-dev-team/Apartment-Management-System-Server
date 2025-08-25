package services

import (
	"api/config"
	"api/models"
	"api/repositories"

	"gopkg.in/gomail.v2"
)

type EmailQueueService struct {
	emailQueueRepository     *repositories.EmailQueueRepository
	emailQueueFailRepository *repositories.EmailQueueFailRepository
}

func NewEmailQueueService() *EmailQueueService {
	return &EmailQueueService{
		emailQueueRepository:     repositories.NewEmailQueueRepository(),
		emailQueueFailRepository: repositories.NewEmailQueueFailRepository(),
	}
}

func (s *EmailQueueService) SendEmail() {
	var emails []models.EmailQueueModel
	err := s.emailQueueRepository.Get(&emails)

	if err != nil {
		return
	}

	for _, email := range emails {
		receiverEmail := email.ReceiverEmail
		testEmail := config.GetEnv("TEST_MAIL_TO")
		if testEmail != "" {
			receiverEmail = testEmail
		}

		// Send email
		message := gomail.NewMessage()
		message.SetHeader("From", config.MailFromAddress, config.MailFromName)
		message.SetHeader("To", receiverEmail)
		message.SetHeader("Subject", email.Subject)
		message.SetBody("text/html", email.Body)

		if err := config.Mailer.DialAndSend(message); err != nil {
			// Add to fail queue
			s.emailQueueFailRepository.Create(&models.EmailQueueFailModel{
				ReceiverEmail: receiverEmail,
				Subject:       email.Subject,
				Body:          email.Body,
				Error:         err.Error(),
			})

		}
		s.emailQueueRepository.Delete(email.ID)
	}
}
