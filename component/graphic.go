package component

import (
	"github.com/dwethmar/beherit/cmd/mygame/sprite"
)

const GraphicType Type = "graphic"

type Graphic struct {
	Sprite  sprite.Sprite
	Width   int
	Height  int
	OffSetX int
	OffSetY int
}

func NewGraphicComponent() *Component {
	return &Component{
		Type: GraphicType,
		Data: &Graphic{
			Sprite: sprite.Sprite{},
		},
	}
}
