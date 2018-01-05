package sender

import (
	"foodcourt/logger"
	"foodcourt/orders"
	"foodcourt/database"
	"foodcourt/environment"
	"foodcourt/settings"
	"html/template"
	"time"
	"fmt"
	"errors"
	"github.com/op/go-logging"
	"gopkg.in/gomail.v2"
	"bytes"
)

const dateFmt = "20060102"
const datePrettyFmt = "Monday, 02 January 2006"

const errTplParse = "failed to parse template '%s': %v"
const errSmtpPort = "SMTP port must be an integer"
const errSendMail = "failed to send order mail to %s: %v"

const emailSubjectTemplate = "Your order for %s"

const logMsgSent = "Message send successfully to %s"

var orderTemplatePath = environment.GetResourcePath("order-mail.html")
var orderTemplate *template.Template
var log *logging.Logger

func SendLunchOrders() (bool, error) {

	log = logger.GetLogger()
	cfgPtr, err := settings.GetSettings()

	if err != nil {
		log.Error(fmt.Sprintf("Failed to get settings: %v", err))
		return false, err
	}

	// System settings
	cfg := *cfgPtr

	if !cfg.Sender.Enable {
		log.Warning(msgSenderDisabled)
		return false, nil
	}

	today := time.Now().Format(dateFmt)

	log.Info(fmt.Sprintf("Processing orders for date: %s", today))

	db := database.GetInstance()
	defer db.Close()

	var ordersList []orders.OrderSummary

	err = orders.GetOrderSummary(&ordersList, today, db)

	if err != nil {
		log.Error(fmt.Sprintf("Failed to fetch orders for %s: %s", today, err.Error()))
		return false, err
	}

	if len(ordersList) == 0 {
		log.Info(fmt.Sprintf("No orders for today (%s)", today))
		return true, nil
	}

	failedMails := make([]string, 0)

	// Get mail sender
	sender, senderErr := getMailSender(cfgPtr)

	if senderErr != nil {
		log.Error(senderErr.Error())
		return false, senderErr
	}

	orderTemplate, err = template.ParseFiles(orderTemplatePath)

	if err != nil {
		log.Error(fmt.Sprintf(errTplParse, orderTemplatePath, err))
		return false, err
	}

	ptrGroup := *orders.GroupOrders(&ordersList)

	for _, orderGroup := range ptrGroup {
		if success := sendLunchOrder(&orderGroup, sender, cfgPtr); !success {
			failedMails = append(failedMails, orderGroup.Email)
		}
	}

	if len(failedMails) > 0 {
		return false, errors.New(fmt.Sprintf("failed to deliver emails %v", failedMails))
	}

	return true, nil
}

// Sends lunch order to specified client
func sendLunchOrder(orderGroup *orders.OrderGroup, sender *gomail.Sender, configPtr *settings.Settings) bool {
	cfg := *configPtr
	vm := OrderMailData{
		DisplayedDate: time.Now().Format(datePrettyFmt),
		Group: *orderGroup,
		BaseURL: cfg.BaseURL,
	}

	userEmail := vm.Group.Email

	// Compose email using html template
	var tplBuff bytes.Buffer

	if err := orderTemplate.Execute(&tplBuff, vm); err != nil {
		log.Error(fmt.Sprintf(errTplParse, orderTemplatePath, err.Error()))
		return false
	}

	// Parse email template
	htmlText := tplBuff.String()

	// Compose email
	mail := gomail.NewMessage()
	mail.SetHeader("From", fmt.Sprintf("FoodCourt <%s>", cfg.Sender.Email))
	mail.SetAddressHeader("To", userEmail, vm.Group.FullName)
	mail.SetHeader("Subject", fmt.Sprintf(emailSubjectTemplate, vm.DisplayedDate))
	mail.SetBody("text/html", htmlText)

	if err := gomail.Send(*sender, mail); err != nil {
		log.Error(fmt.Sprintf(errSendMail, userEmail, err))
		return false
	}

	log.Info(fmt.Sprintf(logMsgSent, userEmail))

	return true
}
