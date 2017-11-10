package dishes

import (
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
func All(dest *[]Dish, db *sqlx.DB) error {
	query, _, _ := sq.Select("*").From(Table).ToSql()
	return db.Select(dest, query)
}

// Find dish by id
func FindById(id string, dest *Dish, db *sqlx.DB) error {
	return Select("*", dest, sq.Eq{Id: id}, db)
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


