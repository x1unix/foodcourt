package sender

import (
	"../../shared/orders"
	"../../shared/logger"
	"../../shared/database"
	"../../shared/dishes"
	. "github.com/ahmetb/go-linq"
	"gopkg.in/gomail.v2"
	"html/template"
	"time"
	"fmt"
	"strings"
)

const errFetchReportData = "failed to fetch orders from the db: %v (%s)"

type OrderReportMailData struct {
	DateString string
	Orders map[string] int
}


func SendOrderReport() (bool, error) {
	log := logger.GetLogger()
	today := time.Now().Format(dateFmt)

	db := database.GetInstance()
	defer db.Close()

	var ordersList []orders.Order

	if err := orders.GetTotalOrderedDishes(&ordersList, today, db); err != nil {
		log.Error(fmt.Sprintf(errFetchReportData, err, today))
		return false, err
	}

	merged := mergeOrders(&ordersList)

	fmt.Println(*merged)

	return true, nil
}

func sendOrderMail(date string, receivers []string, items *map[string] int) error {
	tpl :=
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