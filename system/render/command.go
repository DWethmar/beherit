package render

import (
	"github.com/dwethmar/beherit/command"
	"github.com/hajimehoshi/ebiten/v2"
)

const RenderCommandType = "render"

type RenderCommand struct {
	Screen *ebiten.Image `json:"screen"`
}

func NewRenderCommand() *command.Command {
	return &command.Command{
		Type: RenderCommandType,
		Data: &RenderCommand{},
	}
}
