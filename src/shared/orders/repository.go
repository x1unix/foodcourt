package orders

import (
	"fmt"

	"../dishes"
	"../menu"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// Order dishes from the menu for specific date
func OrderDishes(dishIds []int, date int, userId int, db *sqlx.DB) error {
	// Start transaction
	tx := db.MustBegin()

	// Delete previous items
	tx.MustExec(sqOrdersPurge, date, userId)

	// Build insertion que`ry
	sqlQ, args, _ := squirrel.Insert(Table).Columns(ItemId, UserId).
		Select(squirrel.Select("m.row_id, u.id").From(menu.Table + " m").Join("users u").
			Where(squirrel.Eq{"m.dish_id": dishIds, "m.date": date, "u.id": userId})).ToSql()

	// Add insertion query
	tx.MustExec(sqlQ, args...)

	// Commit transaction
	commitErr := tx.Commit()

	// Check for transaction commit errors
	if commitErr != nil {
		// Rollback on error
		rollbackErr := tx.Rollback()

		log.Error(fmt.Sprintf(orderErrMsg, userId, date, dishIds, commitErr.Error()))

		// Handle rollback errors
		if rollbackErr != nil {
			log.Error(fmt.Sprintf(rollbackErrMsg, userId, date, dishIds, rollbackErr.Error()))
		}
	}

	return commitErr
}

// Get list of ordered dishes
func GetUserOrderMenuItems(output *[]int, userID int, date int, db *sqlx.DB) error {
	// SELECT m.dish_id FROM `menu` m JOIN `orders` o on o.`item_id` = m.`row_id` where  WHERE m.`date`=? AND o.`user_id`=?
	q, a, _ := squirrel.Select("m.dish_id").From(menu.Table + " m").
		Join(Table + " o on o.item_id = m.row_id").
		Where(squirrel.Eq{"m.date": date, "o.user_id": userID}).
		ToSql()

	return db.Select(output, q, a...)
}

// DeleteOrder deletes order
func DeleteOrder(date int, userID int, db *sqlx.DB) error {
	_, err := db.Query(sqOrdersPurge, date, userID)
	return err
}

// GetOrderedDishes returns list of ordered dishes for specific date by specific user
func GetOrderedDishes(output *[]dishes.Dish, date int, userID int, db *sqlx.DB) error {
	/*
		SELECT d.id,
			d.label,
			d.description,
			d.type,
			d.photo_url
		FROM   orders o
			JOIN menu m
				ON o.item_id = m.row_id
			JOIN dishes d
				ON d.id = m.dish_id
		WHERE  o.user_id = 1
			AND m.date = 20171130;
	*/

	q, a, _ := squirrel.Select("d.id, d.label, d.description, d.type, d.photo_url").From(Table + " o").
		Join(menu.Table + " m on o.item_id = m.row_id").
		Join(dishes.Table + " d on d.id = m.dish_id").
		Where(squirrel.Eq{"o.user_id": userID, "m.date": date}).
		ToSql()

	return db.Select(output, q, a...)
}


// Gets list of users that made orders in range ow two dates
func GetOrderStatsBetweenDates(output *[]UserOrderCounter, dateFrom int, dateTill int, db *sqlx.DB) error {
	 /**
select o.user_id,
	m.date
from orders o
	join menu m
		on o.item_id = m.row_id
where m.date >= 20171216 and m.date <= 20171219
	group by o.user_id, m.date;
	  */

	q, a, _ := squirrel.Select("o.user_id, m.date").From(Table + " o").
		Join(menu.Table + "m o.item_id = m.row_id").
		Where("m.date >= ? and m.date <= ?", dateFrom, dateTill).
		GroupBy("o,user_id, m.date").
		ToSql()

	return db.Select(output, q, a...)
}

/*
	select o.user_id,
		count(distinct o.item_id) as count,
		m.date
	from orders o
		join menu m
			on o.item_id = m.row_id
	where m.date = 20171219
		group by o.user_id;
 */
