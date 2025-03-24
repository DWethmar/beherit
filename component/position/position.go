package position

import "github.com/dwethmar/beherit/component"

const Type = "position"

type Position struct {
	X, Y int
}

func NewComponent() *component.Component {
	return &component.Component{
		Type: Type,
		Data: &Position{},
	}
}
