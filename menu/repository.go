package menu

import (
	"fmt"
	"foodcourt/cache"
	"foodcourt/dishes"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type MenuPeriod map[int][]dishes.Dish

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

// GetMenuForPeriod returns menu for specific period
func GetMenuForPeriod(dateFrom int, dateTill int, db *sqlx.DB) (*MenuPeriod, error) {
	// Select all dishes in menu for specific period
	q, a, _ := squirrel.Select("d.label", "d.description", "d.photo_url", "d.type", "m.dish_id as id, m.date").
		From(dishes.Table + " d").
		Join("menu m on d.id = m.dish_id").
		Where("m.date >= ? and m.date <= ?", dateFrom, dateTill).
		OrderBy("m.date ASC").
		ToSql()

	rows, err := db.Query(q, a...)

	// Check for query errors
	if err != nil {
		return nil, err
	}

	// Prepare output result
	results := make(MenuPeriod)

	// Iterate through each row to group by date
	for rows.Next() {
		var id int
		var label string
		var description string
		var photoUrl string
		var dishType int
		var date int

		// Extract row values
		if err = rows.Scan(&label, &description, &photoUrl, &dishType, &id, &date); err != nil {
			return nil, err
		}

		// Put dish
		results[date] = append(results[date], dishes.Dish{
			Id: id,
			Label: label,
			Description: description,
			PhotoUrl: photoUrl,
			Type: dishType,
		})
	}

	// Put connection back to the pool
	rows.Close()

	return &results, err
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

// GetMenusLockStatus returns lock status for each menu
func GetMenusLockStatus(dates []int) (*map[int]bool, error) {
	keys := make([]string, len(dates))
	out := make(map[int]bool)

	// Convert date to redis key
	for i, date := range dates {
		keys[i] = getMenuLockKey(date)
	}

	data, err := cache.Client.MGet(keys...).Result();

	if err != nil {
		return nil, err
	}

	for i, value := range data {
		fmt.Printf("%v %T", value, value)
		out[dates[i]] = value != nil;
	}

	return &out, nil
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
