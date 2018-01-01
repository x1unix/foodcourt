package sender

import (
	"../../shared/orders"
)

type OrderMailData struct {
	Group orders.OrderGroup
	BaseURL string
}
