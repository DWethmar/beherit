package render

import (
	"errors"
	"fmt"
	"image/color"
	"log/slog"

	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var red = color.RGBA{0xff, 0x00, 0x00, 0xff}

type Renderer struct {
	logger *slog.Logger
	cm     *component.Manager
}

func New(logger *slog.Logger, cm *component.Manager) *Renderer {
	return &Renderer{
		logger: logger,
		cm:     cm,
	}
}

func (r *Renderer) Attach(commandFactory *command.Factory, commandBus *command.Bus) {
	commandFactory.Register(NewRenderCommand)
	commandBus.RegisterHandler(NewRenderCommand(), r)
}

func (r *Renderer) Handle(cmd *command.Command) error {
	switch cmd.Type {
	case RenderCommandType:
		c, ok := cmd.Data.(*RenderCommand)
		if !ok {
			return fmt.Errorf("could not convert to render command: %v", cmd.Data)
		}
		return r.Draw(c.Screen)
	}
	return nil
}

func (r *Renderer) Draw(dst *ebiten.Image) error {
	if dst == nil {
		return errors.New("nil destination image")
	}
	for _, g := range r.cm.List(component.PositionType) {
		p, ok := g.Data.(*component.Position)
		if !ok {
			return fmt.Errorf("could not convert to position: %v", g.Data)
		}
		vector.DrawFilledRect(dst, float32(p.X), float32(p.Y), 16, 16, red, true)
	}
	return nil
}
