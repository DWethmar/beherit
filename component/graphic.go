package component

const GraphicType Type = "graphic"

type Graphic struct {
	Name    string `json:"name"`
	Resize  bool   `json:"resize"`
	ResizeW int    `json:"resize_w"`
	ResizeH int    `json:"resize_h"`
	OffSetX int    `json:"offset_x"`
	OffSetY int    `json:"offset_y"`
}

func NewGraphicComponent() *Component {
	return &Component{
		Type: GraphicType,
		Data: &Graphic{},
	}
}
