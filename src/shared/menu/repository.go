package menu

import (
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
)

// Table name
const Table = "menu"

// Table rows
const RowId, DishId, Date = "row_id", "dish_id", "date"

// Add the dish to the menu for the specific date.
// The date must be in format YYYYMMDD (20171116)
func AddDish(dishId int, date int, db *sqlx.DB) error {
	_, err := squirrel.Insert(Table).Columns(DishId, Date).Values(dishId, date).RunWith(db.DB).Exec()
	return err
}

// Check if specified dish exists in the menu for specific date.
func DishExistsInMenu(dishId int, date int, db *sqlx.DB) (error, bool) {
	matcher := squirrel.Eq{DishId: dishId, Date: date}
	q, args, _ := squirrel.Select("COUNT(*)").From(Table).Where(matcher).ToSql()

	var count int

	err := db.Get(&count, q, args...)

	return err, count > 0
}

// Get menu items for specific date.
// The date must be in format YYYYMMDD (20171116)
func GetMenuItems(output *[]MenuItem, date int, db *sqlx.DB) error {
	query, args, _ := squirrel.Select("*").From(Table).Where(squirrel.Eq{Date: date}).ToSql()
	return db.Select(output, query, args)
}

// Delete an specific menu item from collection
func DeleteItem(itemId int, db *sqlx.DB) error {
	_, err := squirrel.Delete(Table).Where(squirrel.Eq{RowId: itemId}).RunWith(db.DB).Exec()
	return err
}

// Delete all menu items for specific date.
// The date must be in format YYYYMMDD (20171116)
func DeleteMenu(date int, db *sqlx.DB) error {
	_, err := squirrel.Delete(Table).Where(squirrel.Eq{Date: date}).RunWith(db.DB).Exec()
	return err
}