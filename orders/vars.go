package orders

import (
	"foodcourt/logger"
)

var log = logger.GetLogger()

/*
  === Table and columns ===
*/
const Table = "orders"

const OrderId, ItemId, UserId = "order_id", "item_id", "user_id"

/*
  === SQL Queries ===
*/

// SQL query for purging previous order for menu
const sqOrdersPurge = "DELETE o FROM orders o INNER JOIN menu m ON o.item_id = m.row_id WHERE m.date = ? AND o.user_id = ?;"

// SQL query for ordering dishes
const sqOrdersGetId = "INSERT INTO orders (item_id, user_id) SELECT m.row_id, u.id FROM menu m JOIN users u WHERE m.dish_id = ? AND m.date = ? AND u.id = ?;"

/*
  === Error messages ===
*/

// Rollback error
const rollbackErrMsg = "Order error rollback fail: (user: %d) (date: %d) of dishes: %v; Error: %s"

// Order error
const orderErrMsg = "Order make fail: (user: %d) (date: %d) of dishes: %v; Error: %s"
