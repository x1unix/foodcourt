package sender

import (
	"../../shared/logger"
	"../../shared/orders"
	"../../shared/database"
	"../../shared/environment"
	"../../shared/config"
	"html/template"
	"time"
	"fmt"
	"errors"
	"github.com/op/go-logging"
	"gopkg.in/gomail.v2"
	"bytes"
)

const dateFmt = "20060102"

const errTplParse = "failed to parse template '%s': %v"
const errSmtpPort = "SMTP port must be an integer"
const errSendMail = "failed to send order mail to %s: %v"

const emailSubjectTemplate = "Your order for today"

const logMsgSent = "Message send successfully to %s"

var orderTemplatePath = environment.GetResourcePath("order-mail.html")
var orderTemplate *template.Template
var log *logging.Logger

func SendLunchOrders() (bool, error) {

	log = logger.GetLogger()
	today := time.Now().Format(dateFmt)

	log.Info(fmt.Sprintf("Processing orders for date: %s", today))

	db := database.GetInstance()
	defer db.Close()

	var ordersList []orders.OrderSummary

	err := orders.GetOrderSummary(&ordersList, today, db)

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
	sender, senderErr := getMailSender()

	if senderErr != nil {
		log.Error(senderErr)
		return false, senderErr
	}

	orderTemplate, err = template.ParseFiles(orderTemplatePath)

	if err != nil {
		log.Error(fmt.Sprintf(errTplParse, orderTemplatePath, err))
		return false, err
	}

	ptrGroup := *orders.GroupOrders(&ordersList)

	for _, orderGroup := range ptrGroup {
		if success := sendLunchOrder(&orderGroup, sender); !success {
			failedMails = append(failedMails, orderGroup.Email)
		}
	}

	//From(ordersList).GroupBy(func(order interface{}) interface{} {
	//	return order.(orders.OrderSummary).Email
	//}, func(order interface{}) interface{} {
	//	return order.(orders.OrderSummary)
	//}).ForEach(func (i interface{}) {
	//	group := i.(Group)
	//	email := group.Key.(string)
	//	success := sendLunchOrder(sender, email, group.Group)
	//
	//	if !success {
	//		failedMails = append(failedMails, email)
	//	}
	//})

	if len(failedMails) > 0 {
		return false, errors.New(fmt.Sprintf("failed to deliver emails %v", failedMails))
	}

	return true, nil
}

// Sends lunch order to specified client
func sendLunchOrder(orderGroup *orders.OrderGroup, sender *gomail.Sender) bool {
	vm := OrderMailData{
		Group: *orderGroup,
		BaseURL: config.Get(config.BASE_URL, "#"),
	}

	userEmail := vm.Group.Email

	// Compose email using html template
	// tpl, err := orderTemplate.ParseFiles(orderTemplatePath)

	var tplBuff bytes.Buffer

	if err := orderTemplate.Execute(&tplBuff, vm); err != nil {
		log.Error(fmt.Sprintf(errTplParse, orderTemplatePath, err.Error()))
		return false
	}

	// Parse email template
	htmlText := tplBuff.String()

	// Compose email
	mail := gomail.NewMessage()
	mail.SetHeader("From", config.Get(config.SmtpFrom, "voracity"))
	mail.SetAddressHeader("To", userEmail, vm.Group.FullName)
	mail.SetHeader("Subject", emailSubjectTemplate)
	mail.SetBody("text/html", htmlText)

	if err := gomail.Send(*sender, mail); err != nil {
		log.Error(fmt.Sprintf(errSendMail, userEmail, err))
		return false
	}

	log.Info(fmt.Sprintf(logMsgSent, userEmail))

	return true
}
