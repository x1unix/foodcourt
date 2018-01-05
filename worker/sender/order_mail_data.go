package sender

import (
	"foodcourt/orders"
)

type OrderMailData struct {
	DisplayedDate string
	Group orders.OrderGroup
	BaseURL string
}
