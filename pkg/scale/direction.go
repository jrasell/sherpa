package scale

type Direction string

const (
	DirectionIn   Direction = "in"
	DirectionOut  Direction = "out"
	DirectionNone Direction = "none"
)

func (d *Direction) String() string {
	return string(*d)
}
