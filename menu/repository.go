package menu

import (
	"fmt"
	"foodcourt/cache"
	"foodcourt/dishes"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// Table name
const Table = "menu"

// Table rows
const RowId, DishId, Date = "row_id", "dish_id", "date"

// SQL query for selecting dishes in menu
const sqlQueryMenuDishes = "select d.label, d.description, d.photo_url, d.type, m.dish_id as id from dishes d inner join menu m on d.id = m.dish_id where m.date = ?"

// Redis key prefix for menus
const rdKeyMenuPrefix = "menu_lock_%d"

// Gets list of dishes in menu at specific date
func GetDishesInMenu(output *[]dishes.Dish, date int, db *sqlx.DB) error {
	return db.Select(output, sqlQueryMenuDishes, date)
}

// Delete all dishes from the menu
func ClearMenu(date int, db *sqlx.DB) error {
	_, err := squirrel.Delete(Table).Where(squirrel.Eq{Date: date}).RunWith(db.DB).Exec()
	return err
}

// Add the dish to the menu for the specific date.
// The date must be in format YYYYMMDD (20171116)
func AddDish(dishId int, date int, db *sqlx.DB) error {
	_, err := squirrel.Insert(Table).Columns(DishId, Date).Values(dishId, date).RunWith(db.DB).Exec()
	return err
}

// Set the new list of dishes for specific date.
// All previous dishes will be deleted.
func SetDishesForDate(dishesIds []int, date int, db *sqlx.DB) error {
	// SQL query to delete all previous items
	delQ, delArgs, _ := squirrel.Delete(Table).Where(squirrel.Eq{Date: date}).ToSql()
	insertBuilder := squirrel.Insert(Table).Columns(DishId, Date)

	for _, dishId := range dishesIds {
		insertBuilder = insertBuilder.Values(dishId, date)
	}

	insQ, insArgs, _ := insertBuilder.ToSql()

	// Start transaction
	tx := db.MustBegin()
	tx.MustExec(delQ, delArgs...)
	tx.MustExec(insQ, insArgs...)

	return tx.Commit()
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

// Gets item ids of dishes in menu
func GetMenuItemsIds(output *[]int, dishIds []int, date int, db *sqlx.DB) error {
	// Build sql query
	q, a, _ := squirrel.Select(RowId).From(Table).Where(squirrel.Eq{DishId: dishIds, Date: date}).ToSql()

	return db.Select(output, q, a...)
}

func getMenuLockKey(date int) string {
	return fmt.Sprintf(rdKeyMenuPrefix, date)
}

// Check if menu is locked for the new orders
func GetMenuLockStatus(date int) (bool, error) {
	key := getMenuLockKey(date)

	val, keErr := cache.Client.Exists(key).Result()
	exists := val > 0

	return exists, keErr
}

// Set menu status
func SetMenuLockStatus(lockState bool, date int) error {
	isLocked, chkErr := GetMenuLockStatus(date)

	if chkErr != nil {
		return chkErr
	}

	rdKey := getMenuLockKey(date)

	if lockState {
		if !isLocked {
			// Set lock key
			_, err := cache.Client.Set(rdKey, true, 0).Result()
			return err
		}
	} else {
		if isLocked {
			// Unset lock key
			_, err := cache.Client.Del(rdKey).Result()
			return err
		}
	}

	return nil
}
