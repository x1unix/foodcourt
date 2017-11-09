package dishes

const (
	Soup = 0
	Garnish = 1
	Main = 2
	Salad = 3
)


type Dish struct {
	Id int `db:"id" json:"id"`
	Label string `db:"label" json:"label" validate:"required"`
	Description string `db:"description" json:"description,omitempty"`
	Type int `db:"type" json:"type"`
	PhotoUrl string `db:"photo_url" json:"photoUrl"`
	HasGarnish bool `db:"has_garnish" json:"hasGarnish"`
}
