package dishes

import (
	"../database"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const Table = "dishes"
const Id = "id"

var Columns = []string{"label", "description", "type", "photo_url"}

func Select(cols string, destination interface{}, where sq.Eq, db *sqlx.DB) error {
	query, args, _ := sq.Select(cols).From(Table).Where(where).ToSql()

	return db.Get(destination, query, args...)
}

// Get all items
func All(dest *[]Dish, querySelector *database.QuerySelector, db *sqlx.DB) error {
	query, args, _ := querySelector.ApplyOnSelect(sq.Select("*").From(Table)).ToSql()

	if len(args) > 0 {
		return db.Select(dest, query, args)
	} else {
		return db.Select(dest, query)
	}
}

// Find dish by id
func FindById(id string, dest *Dish, db *sqlx.DB) error {
	return Select("*", dest, sq.Eq{Id: id}, db)
}

// Delete an item
func Delete(id string, db *sqlx.DB) error {
	_, err := sq.Delete(Table).Where(sq.Eq{Id: id}).RunWith(db.DB).Exec()
	return err
}

// Delete multiple items by id
func DeleteIds(ids []int, db *sqlx.DB) error {
	_, err := sq.Delete(Table).Where(sq.Eq{Id: ids}).RunWith(db.DB).Exec()
	return err
}

// Check if dish exists
func Exists(id string, db *sqlx.DB) (error, bool) {
	q, args, _ := sq.Select("COUNT(*)").From(Table).Where(sq.Eq{Id: id}).ToSql()
	var count int

	err := db.Get(&count, q, args...)

	return err, count > 0
}

// Add new dish
func Add(dish *Dish, db *sqlx.DB) error {
	_, err := sq.Insert(Table).
		Columns(Columns...).
		Values(dish.Label, dish.Description, dish.Type, dish.PhotoUrl).
		RunWith(db.DB).
		Query()

	return err
}
