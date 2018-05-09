package mails

import (
	"foodcourt/settings"
	"gopkg.in/gomail.v2"
	"fmt"
	"html/template"
	"bytes"
)

type MailSender struct {
	Settings settings.SMTPSettings
	TemplatePath string
	from string
	subject string
	template *template.Template
	connected bool
	sender gomail.SendCloser
}

type MailRecipient struct {
	Email string
	FullName string
}

func (m *MailSender) Init() error {
	dialer := gomail.NewDialer(m.Settings.Host, m.Settings.Port, m.Settings.User, m.Settings.Password)

	if sender, err := dialer.Dial(); err != nil {
		return err
	} else {
		m.sender = sender
		return nil
	}

	m.connected = true

	tpl, err := template.ParseFiles(m.TemplatePath)

	if err != nil {
		return fmt.Errorf("Failed to read template '%s': %v", m.TemplatePath, err)
	}

	m.template = tpl
	return nil
}

func (m *MailSender) SetSubject(subject string) *MailSender {
	m.subject = subject
	return m
}

func (m *MailSender) From(name string, email string) *MailSender {
	m.from = fmt.Sprintf("%s <%s>", name, email)
	return m
}

func (m *MailSender) Send(contents interface{}, recipient *MailRecipient) error {
	var tplBuff bytes.Buffer

	// Compose email contents using html template
	if err := m.template.Execute(&tplBuff, contents); err != nil {
		return err
	}

	// Parse email template
	htmlText := tplBuff.String()

	// Compose email
	mail := gomail.NewMessage()
	mail.SetHeader("From", m.from)
	mail.SetAddressHeader("To", recipient.Email, recipient.FullName)
	mail.SetHeader("Subject", m.subject)
	mail.SetBody("text/html", htmlText)

	return gomail.Send(m.sender, mail);
}


func NewRecipient(email string, name string) *MailRecipient {
	return &MailRecipient{
		email,
		name,
	}
}


func NewMailSender(templatePath string) (*MailSender, error) {
	cfg, err := settings.GetSettings()

	if err != nil {
		return nil, err
	}

	var sender MailSender
	sender.Settings = cfg.SMTP

	return sender.From("FoodCourt", cfg.Sender.Email), nil
}
