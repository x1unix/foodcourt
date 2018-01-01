package orders

import (
	"../dishes"
)

type OrderSummary struct {
	Email string `db:"email"`
	Label string `db:"label"`
	FirstName string `db:"firstName"`
	LastName string `db:"lastName"`
	Description string `db:"description"`
	PhotoUrl string `db:"photo_url"`
	Type int `db:"type"`
}

func (order *OrderSummary) ToDish() *dishes.Dish {
	dish := dishes.Dish{
		Label: order.Label,
		Description: order.Description,
		PhotoUrl: order.PhotoUrl,
		Type: order.Type,
	}
	return &dish
}