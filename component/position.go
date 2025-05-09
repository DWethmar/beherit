package component

const PositionType Type = "position"

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewPositionComponent() *Component {
	return &Component{
		Type: PositionType,
		Data: &Position{},
	}
}
