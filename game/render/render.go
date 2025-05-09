package render

import (
	"errors"
	"fmt"
	"image/color"
	"log/slog"

	"github.com/dwethmar/beherit/cmd/mygame/sprite"
	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
	"github.com/dwethmar/beherit/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var red = color.RGBA{0xff, 0x00, 0x00, 0xff}

type Renderer struct {
	logger  *slog.Logger
	cm      *component.Manager
	sprites map[string]*sprite.Sprite
}

func New(logger *slog.Logger, cm *component.Manager, sprites map[string]*sprite.Sprite) *Renderer {
	return &Renderer{
		logger:  logger,
		cm:      cm,
		sprites: sprites,
	}
}

func (r *Renderer) Handle(cmd *command.Command) error {
	switch cmd.Type {
	case game.RenderCommandType:
		c, ok := cmd.Data.(*game.RenderCommand)
		if !ok {
			return fmt.Errorf("could not convert to render command: %v", cmd.Data)
		}
		return r.Draw(c.Screen, c.Components)
	default:
		return fmt.Errorf("unknown command type: %s", cmd.Type)
	}
}

func (r *Renderer) Draw(dst *ebiten.Image, components []component.Component) error {
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
