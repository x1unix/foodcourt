package sender

import (
	"foodcourt/settings"
	"gopkg.in/gomail.v2"
	"fmt"
)

const errSmtpDial = "failed to contact with SMTP server: %s"
const msgSenderDisabled = "Email notifications are disabled. Please enable notifications in System settings";


func getMailSender(configPtr *settings.Settings) (*gomail.Sender, error) {
	// Get configuration from pointer
	cfg := *configPtr

	// Create mail sender instance
	dialer := gomail.NewDialer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.User, cfg.SMTP.Password)

	var sender gomail.Sender
	var err error

	// Contact with SMTP server
	if sender, err = dialer.Dial(); err != nil {
		log.Error(fmt.Sprintf(errSmtpDial, err.Error()))
		return nil, err
	}

	return &sender, nil
}
