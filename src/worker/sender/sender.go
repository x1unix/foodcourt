package sender

import (
	"../../shared/config"
	"errors"
	"gopkg.in/gomail.v2"
	"strconv"
	"fmt"
)

const errSmtpDial = "failed to contact with SMTP server: %s"


func getMailSender() (*gomail.Sender, error) {

	// Read SMTP configuration
	params := make(map[string] string)
	err := config.GetMultiple(&params, config.SmtpHost, config.SmtpPort, config.SmtpUser, config.SmtpPass)

	if err != nil {
		return nil, err
	}

	// Parse SMTP port (string -> int)
	var smtpPort int
	smtpPort, err = strconv.Atoi(params[config.SmtpPort])

	if err != nil {
		return nil, errors.New(errSmtpPort)
	}

	// Create mail sender instance
	dialer := gomail.NewDialer(params[config.SmtpHost], smtpPort, params[config.SmtpUser], params[config.SmtpPass])

	var sender gomail.Sender

	// Contact with SMTP server
	if sender, err = dialer.Dial(); err != nil {
		log.Error(fmt.Sprintf(errSmtpDial, err.Error()))
		return nil, err
	}

	return &sender, nil
}
