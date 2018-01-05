package orders

import (
	"fmt"
	"foodcourt/dishes"
)

type OrderGroup struct {
	FullName string
	Email    string
	Items    []dishes.Dish
}

// Groups list of order summaries to order groups
func GroupOrders(orders *[]OrderSummary) *[]OrderGroup {
	output := make([]OrderGroup, 0)
	groupIndexes := make(map[string]int)
	list := *orders

	for _, order := range list {
		if grpIdx, idxExists := groupIndexes[order.Email]; idxExists {
			output[grpIdx].Items = append(output[grpIdx].Items, *order.ToDish())
		} else {
			grpIdx = len(output)
			groupIndexes[order.Email] = grpIdx
			output = append(output, OrderGroup{
				FullName: fmt.Sprintf("%s %s", order.FirstName, order.LastName),
				Email:    order.Email,
				Items:    make([]dishes.Dish, 0),
			})

			output[grpIdx].Items = append(output[grpIdx].Items, *order.ToDish())
		}
	}

	groupIndexes = make(map[string]int)
	return &output
}
