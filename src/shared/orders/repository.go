package orders

import (
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"../dishes"
	"../menu"
)

const Table = "orders"

const OrderId, ItemId, UserId = "order_id", "item_id", "user_id";

func OrderDishes(dishIds []int, date int, db *sqlx.DB) error {

}