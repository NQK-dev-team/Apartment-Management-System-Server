package config

import (
	"strconv"

	"gopkg.in/gomail.v2"
)

var Mailer *gomail.Dialer
var MailFromAddress string
var MailFromName string

func InitMailer() error {
	mailHost, err := GetEnv("MAIL_HOST")
	if err != nil {
		return err
	}

	mailPortString, err := GetEnv("MAIL_PORT")
	if err != nil {
		return err
	}

	mailPort, err := strconv.Atoi(mailPortString)
	if err != nil {
		return err
	}

	mailUsername, err := GetEnv("MAIL_USERNAME")
	if err != nil {
		return err
	}

	mailPassword, err := GetEnv("MAIL_PASSWORD")
	if err != nil {
		return err
	}

	MailFromAddress, _ = GetEnv("MAIL_FROM_ADDRESS")
	MailFromName, _ = GetEnv("MAIL_FROM_NAME")

	Mailer = gomail.NewDialer(mailHost, mailPort, mailUsername, mailPassword)

	return nil
}
