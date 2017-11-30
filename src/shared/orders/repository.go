package orders

import (
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
	"../menu"
	"fmt"
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
func GetUserOrderMenuItems(output *[]int, userId int, date int, db *sqlx.DB) error {
	// SELECT m.dish_id FROM `menu` m JOIN `orders` o on o.`item_id` = m.`row_id` where  WHERE m.`date`=? AND o.`user_id`=?
	q, a, _ := squirrel.Select("m.dish_id").From(menu.Table + " m").
		Join(Table + " o on o.item_id = m.row_id").
		Where(squirrel.Eq{"m.date": date, "o.user_id": userId}).
		ToSql()

	return db.Select(output, q, a...)
}

// Delete order
func DeleteOrder(date int, userId int, db *sqlx.DB) error {
	_, err := db.Query(sqOrdersPurge, date, userId)
	return err
}