package direction

//go:generate go tool stringer -type=Direction
type Direction int

const (
	None Direction = iota
	North
	South
	West
	East
	NorthWest
	NorthEast
	SouthWest
	SouthEast
)
