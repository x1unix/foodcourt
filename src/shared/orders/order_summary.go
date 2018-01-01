package orders

type OrderSummary struct {
	Email string `db:"email"`
	Label string `db:"label"`
	FirstName string `db:"firstName"`
	LastName string `db:"lastName"`
	Description string `db:"description"`
	PhotoUrl string `db:"photo_url"`
	Type int `db:"type"`
}
