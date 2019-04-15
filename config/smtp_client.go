package config

import (
	"fmt"
	"net/smtp"
)

var _ Client = (*SMTPClient)(nil)

type Client interface {
	Validator
	Send(from string, to []string, msg []byte) error
	Host() string
}

type SMTPClient struct {
	config SMTPEmail
	auth   smtp.Auth
	host   string
}

func NewSMAPCLient(cfg SMTPEmail) *SMTPClient {
	host := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth(ApplicationName, cfg.Username, cfg.Password, cfg.Host)
	return &SMTPClient{
		config: cfg,
		auth:   auth,
		host:   host,
	}
}

func (s *SMTPClient) Send(from string, to []string, msg []byte) error {
	return smtp.SendMail(s.host, s.auth, from, to, msg)
}

func (s *SMTPClient) Host() string {
	return s.host
}

func (s *SMTPClient) Valid() *Validation {
	return s.config.Valid()
}
