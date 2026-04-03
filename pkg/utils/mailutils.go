package utils

import (
	"gopkg.in/gomail.v2"
)

type MailService struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewMailService(host string, port int, username, password string) *MailService {
	return &MailService{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (s *MailService) Send(to []string, subject, body string, isHTML bool) {
	m := gomail.NewMessage()
	m.SetHeader("From", s.Username)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	contentType := "text/plain"
	if isHTML {
		contentType = "text/html"
	}
	m.SetBody(contentType, body)
	
	d := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)
	d.SSL = true
	MustNil(d.DialAndSend(m))
}
