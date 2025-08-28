package config

import (
	"errors"
	"strconv"

	"gopkg.in/gomail.v2"
)

var Mailer *gomail.Dialer
var MailFromAddress string
var MailFromName string

func InitMailer() error {
	mailHost := GetEnv("MAIL_HOST")
	if mailHost == "" {
		return errors.New("MAIL_HOST environment variable is not set")
	}

	mailPortString := GetEnv("MAIL_PORT")
	if mailPortString == "" {
		return errors.New("MAIL_PORT environment variable is not set")
	}

	mailPort, err := strconv.Atoi(mailPortString)
	if err != nil {
		return err
	}

	mailUsername := GetEnv("MAIL_USERNAME")
	if mailUsername == "" {
		return errors.New("MAIL_USERNAME environment variable is not set")
	}

	mailPassword := GetEnv("MAIL_PASSWORD")
	if mailPassword == "" {
		return errors.New("MAIL_PASSWORD environment variable is not set")
	}

	MailFromAddress = GetEnv("MAIL_FROM_ADDRESS")
	if MailFromAddress == "" {
		return errors.New("MAIL_FROM_ADDRESS environment variable is not set")
	}

	MailFromName = GetEnv("MAIL_FROM_NAME")
	if MailFromName == "" {
		MailFromName = "Apartment Management System"
	}

	Mailer = gomail.NewDialer(mailHost, mailPort, mailUsername, mailPassword)

	return nil
}
