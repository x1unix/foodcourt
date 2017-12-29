package orders

type OrderSummary struct {
	Email string `db:"email"`
	Label string `db:"label"`
	Description string `db:"description"`
	PhotoUrl string `db:"photo_url"`
	Type int `db:"type"`
}
