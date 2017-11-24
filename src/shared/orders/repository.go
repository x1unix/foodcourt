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

	// Build insertion query
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