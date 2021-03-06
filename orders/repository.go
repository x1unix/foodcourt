package orders

import (
	"fmt"

	"foodcourt/dishes"
	"foodcourt/menu"

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

// Bulk order
func BulkOrderDishes(orders *BulkOrderBundle, userId int, db *sqlx.DB) error {
	// Start transaction
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	for date, orderIds := range *orders {
		// Remove order for a day
		if _, err = tx.Exec(sqOrdersPurge, date, userId); err != nil {
			return fmt.Errorf("failed to remove orders of uid %d on date %d - %v", userId, date, err)
		}

		// Build insertion que`ry
		q, args, _ := squirrel.Insert(Table).Columns(ItemId, UserId).
			Select(squirrel.Select("m.row_id, u.id").From(menu.Table + " m").Join("users u").
			Where(squirrel.Eq{"m.dish_id": orderIds, "m.date": date, "u.id": userId})).ToSql()

		// Execute insert
		if _, err = tx.Exec(q, args...); err != nil {
			return fmt.Errorf("failed to insert order of uid %d for %d - %v", userId, date, err)
		}
	}

	// Execute transaction and pray
	if err = tx.Commit(); err != nil {

		// Try to rollback
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error("failed to rollback failed bulk order transaction (uid:%d) - %v", userId, rollbackErr)
		}

		return fmt.Errorf("failed to commit bulk order of uid %d", userId, err)
	}

	return err
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

// GetUserOrdersForPeriod gets ids of ordered items for specified period
func GetUserOrdersForPeriod(userID int, dateFrom int, dateTill int, db *sqlx.DB) (*map[int][]int, error) {
	out := make(map[int][]int)

	q, a, _ := squirrel.Select("m.dish_id, m.date").From(menu.Table + " m").
		Join(Table + " o on o.item_id = m.row_id").
		Where("m.date >= ? and m.date <= ?", dateFrom, dateTill).
		Where(squirrel.Eq{"o.user_id": userID}).
		ToSql()

	// Run query
	rows, err := db.Query(q, a...)

	if err != nil {
		return nil, err
	}

	// Iterate through each dish and group by date
	for rows.Next() {
		var dishId int
		var date int

		if err = rows.Scan(&dishId, &date); err != nil {
			return nil, err
		}

		out[date] = append(out[date], dishId)
	}


	return &out, err
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

	q, a, _ := squirrel.Select("o.user_id, m.date").From(Table+" o").
		Join(menu.Table+" m on o.item_id = m.row_id").
		Where("m.date >= ? and m.date <= ?", dateFrom, dateTill).
		GroupBy("o.user_id, m.date").
		ToSql()

	return db.Select(output, q, a...)
}

// Gets order summary for specified date
func GetOrderSummary(output *[]OrderSummary, orderDate string, db *sqlx.DB) error {
	/**
	select u.email,
		d.label, d.photo_url, d.type
	from orders o
		join menu m
			on o.item_id = m.row_id
		join dishes d
			on m.dish_id = d.id
		join users u
			on o.user_id = u.id
		where m.date = 20171229;
	*/
	q, a, _ := squirrel.Select("u.email, u.firstName, u.lastName, d.label, d.description, d.photo_url, d.type").From(Table + " o").
		Join(menu.Table + " m on o.item_id = m.row_id").
		Join(dishes.Table + " d on m.dish_id = d.id").
		Join("users u on o.user_id = u.id").
		Where(squirrel.Eq{"m.date": orderDate}).
		ToSql()

	return db.Select(output, q, a...)
}

// Gets
func GetTotalOrderedDishes(output *[]Order, orderDate string, db *sqlx.DB) error {
	/*
		select o.user_id,
				m.date,
				d.label,
				d.type
		from orders o
			join menu m
				on o.item_id = m.row_id
			join dishes d
				on m.dish_id = d.id
		where m.date = 20180102;
	*/
	q, a, _ := squirrel.Select("o.user_id, m.date, d.label, d.type").From(Table + " o").
		Join(menu.Table + " m on o.item_id = m.row_id").
		Join(dishes.Table + " d on m.dish_id = d.id").
		Where(squirrel.Eq{"m.date": orderDate}).
		ToSql()

	return db.Select(output, q, a...)
}
