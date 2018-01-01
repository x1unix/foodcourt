package sender

import (
	"../../shared/orders"
)

type OrderMailData struct {
	DisplayedDate string
	Group orders.OrderGroup
	BaseURL string
}
