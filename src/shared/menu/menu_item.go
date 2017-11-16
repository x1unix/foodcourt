package menu

type MenuItem struct {
	ItemId int 		`db:"row_id" json:"itemId,omitempty"`
	DishId int 		`db:"dish_id" json:"dishId" validate:"required"`
	Date   int 		`db:"date" json:"date" validate:"required"`
}