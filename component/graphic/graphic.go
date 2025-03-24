package graphic

import (
	"github.com/dwethmar/beherit/cmd/mygame/sprite"
	"github.com/dwethmar/beherit/component"
)

const Type = "graphic"

type Graphic struct {
	Sprite  sprite.Sprite
	Width   int
	Height  int
	OffSetX int
	OffSetY int
}

func NewComponent() *component.Component {
	return &component.Component{
		Type: Type,
		Data: &Graphic{
			Sprite: sprite.Sprite{},
		},
	}
}
