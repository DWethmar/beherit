package component

const GraphicType Type = "graphic"

type Graphic struct {
	Name    string
	Width   int
	Height  int
	OffSetX int
	OffSetY int
}

func NewGraphicComponent() *Component {
	return &Component{
		Type: GraphicType,
		Data: &Graphic{},
	}
}
