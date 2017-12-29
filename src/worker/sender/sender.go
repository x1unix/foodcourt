package sender

import (
	"../../shared/logger"
	"../../shared/orders"
	"../../shared/database"
	"../../shared/config"
	. "github.com/ahmetb/go-linq"
	"time"
	"fmt"
	"strings"
	"errors"
)

const dateFmt = "20060102"
const logMsgSendOk = "Send success to %s"
const logMsgSendFail = "Send fail to %s: %s"

func SendLunchOrders() (bool, error) {

	log := logger.GetLogger()
	today := time.Now().Format(dateFmt)

	log.Info(fmt.Sprintf("processing orders for date: %s", today))

	db := database.GetInstance()
	defer db.Close()

	var ordersList []orders.OrderSummary
	failedMails := make([]string, 0)

	err := orders.GetOrderSummary(&ordersList, today, db)

	if err != nil {
		log.Error(fmt.Sprintf("failed to fetch orders for %s: %s", today, err.Error()))
		return false, err
	}

	From(ordersList).GroupBy(func(order interface{}) interface{} {
		return order.(orders.OrderSummary).Email
	}, func(order interface{}) interface{} {
		return order.(orders.OrderSummary)
	}).ForEach(func (i interface{}) {
		group := i.(Group)
		email := group.Key.(string)
		success := sendLunchOrder(email, group.Group)

		if !success {
			failedMails = append(failedMails, email)
		}
	})

	if len(failedMails) > 0 {
		return false, errors.New(fmt.Sprintf("failed to deliver emails: %s", strings.Join(failedMails, ", ")))
	}

	return true, nil
}


func sendLunchOrder(email string, items []interface{}) bool {
	vm := OrderMailData{
		Email: email,
		BaseURL: config.Get(config.BASE_URL, "#"),
		Orders: items,
	}


	return false
}