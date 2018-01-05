package orders

type UserOrderCounter struct {
	UserId	int	`db:"user_id" json:"userId"`
	Date	int	`db:"date" json:"date"`
}
