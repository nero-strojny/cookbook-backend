package email

import (
	"net/smtp"
	"server/config"
)

var (
	from     = "tasty.boi.shopping.list@gmail.com"
	smtpHost = "smtp.gmail.com"
	smtpPort = "587"
)

func Send(subject string, mime string, body string, recipients []string) error {
	auth := smtp.PlainAuth("", from, config.GetConfig().EmailPassword, smtpHost)
	email := []byte(subject + mime + body)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, recipients, email)
}
