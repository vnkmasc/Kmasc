package utils

import (
	"net/smtp"
)

type EmailSender interface {
	SendEmail(to, subject, body string) error
}

type smtpSender struct {
	from     string
	password string
	host     string
	port     string
}

func NewSMTPSender(from, password, host, port string) EmailSender {
	return &smtpSender{from: from, password: password, host: host, port: port}
}

func (s *smtpSender) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.from, s.password, s.host)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	addr := s.host + ":" + s.port
	return smtp.SendMail(addr, auth, s.from, []string{to}, msg)
}
