package sender

import (
	"../../shared/orders"
	"../../shared/logger"
	"../../shared/database"
	"../../shared/dishes"
	"../../shared/settings"
	"../../shared/environment"
	. "github.com/ahmetb/go-linq"
	"gopkg.in/gomail.v2"
	"html/template"
	"time"
	"fmt"
	"strings"
	"bytes"
)

const errFetchReportData = "failed to fetch orders from the db: %v (%s)"
const orderReportTemplateName = "order-summary-mail.html"
const orderMailDateFormat = "01.02.2006"
const orderReportSubject = "Order report for %s"

type OrderReportMailData struct {
	DisplayedDate string
	Orders map[string] int
	BaseURL string
}


func SendOrderReport() (bool, error) {
	log := logger.GetLogger()

	// Read settings
	cfgPtr, err := settings.GetSettings()

	if err != nil {
		log.Error(fmt.Sprintf("Failed to get settings: %v", err))
		return false, err
	}

	// System settings
	cfg := *cfgPtr

	// Check if notifications are enabled
	if !cfg.Sender.Enable {
		log.Warning(msgSenderDisabled)
		return false, nil
	}

	if len(cfg.Sender.ReportRecipients) == 0 {
		log.Warning("No food vendor mails are defined. Please define them in System settings")
		return false, nil
	}

	today := time.Now().Format(dateFmt)

	db := database.GetInstance()
	defer db.Close()

	var ordersList []orders.Order

	if err := orders.GetTotalOrderedDishes(&ordersList, today, db); err != nil {
		log.Error(fmt.Sprintf(errFetchReportData, err, today))
		return false, err
	}

	merged := mergeOrders(&ordersList)

	if err = sendOrderReportMail(cfgPtr, merged); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func sendOrderReportMail(cfgPtr *settings.Settings, items *map[string] int) error {
	log := logger.GetLogger()

	// Extract settings
	cfg := *cfgPtr

	// Get mail sender
	sender, err := getMailSender(cfgPtr)

	if err != nil {
		return err
	}

	// Read template
	orderReportTemplatePath := environment.GetResourcePath(orderReportTemplateName)
	tpl, err := template.ParseFiles(orderReportTemplatePath)

	if err != nil {
		log.Error(fmt.Sprintf(errTplParse, orderReportTemplatePath, err.Error()))
		return err
	}

	mailData := OrderReportMailData{
		BaseURL: cfg.BaseURL,
		DisplayedDate: time.Now().Format(orderMailDateFormat),
		Orders: *items,
	}

	// Compose email using html template
	var tplBuff bytes.Buffer

	if err := tpl.Execute(&tplBuff, mailData); err != nil {
		log.Error(fmt.Sprintf(errTplParse, orderReportTemplatePath, err))
		return err
	}

	// Parse email template
	htmlText := tplBuff.String()

	// Compose email
	mail := gomail.NewMessage()
	mail.SetHeader("From", fmt.Sprintf("FoodCourt <%s>", cfg.Sender.Email))
	mail.SetHeader("To", cfg.Sender.ReportRecipients...)
	mail.SetHeader("Subject", fmt.Sprintf(orderReportSubject, mailData.DisplayedDate))
	mail.SetBody("text/html", htmlText)

	if err := gomail.Send(*sender, mail); err != nil {
		log.Error(fmt.Sprintf("Failed to send order report to %v: %v", cfg.Sender.ReportRecipients, err))
		return err
	}

	log.Info(fmt.Sprintf("Order report sent to %v", cfg.Sender.ReportRecipients))
	return nil
}

// Merges orders by dish type (main+garnish) and counts
func mergeOrders(ordersList *[]orders.Order) *map[string]int {
	result := make(map[string] int)

	From(*ordersList).GroupBy(func(order interface{}) interface{} {
		return order.(orders.Order).UserId
	}, func(order interface{}) interface{} {
		return order
	}).ForEach(func(grp interface{}) {
		group := grp.(Group)
		pairs := make([]string, 0)
		for _, item := range group.Group {
			order := item.(orders.Order)

			if (order.Type == dishes.Garnish) || (order.Type == dishes.Main) {
				pairs = append(pairs, order.Label)
			} else {
				incrementMap(&result, order.Label)
			}
		}

		if len(pairs) > 0 {
			pairLabel := strings.Join(pairs, " + ")
			incrementMap(&result, pairLabel)
			pairs = pairs[:0]
		}
	})

	return &result
}

func incrementMap(mapPtr *map[string] int, key string) {
	mapVal := *mapPtr
	if _, exists := mapVal[key]; exists {
		mapVal[key]++
	} else {
		mapVal[key] = 1
	}
}