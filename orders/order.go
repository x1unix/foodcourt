package orders

import (
	// "../../shared/dishes"
)

type Order struct {
	UserId string `db:"user_id"`
	Date string `db:"date"`
	Label string `db:"label"`
	Type int `db:"type"`
}
//
//func SumOrders(ordersPtr *[]Order) *map[string] int {
//	result := make(map[string] int)
//
//	for _, order := range *ordersPtr {
//
//	}
//}
