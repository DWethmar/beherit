package component

const PositionType Type = "position"

type Position struct {
	X, Y int
}

func NewPositionComponent() *Component {
	return &Component{
		Type: PositionType,
		Data: &Position{},
	}
}
