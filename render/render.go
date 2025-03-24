package render

import (
	"fmt"
	"image/color"

	"github.com/dwethmar/beherit/component"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var red = color.RGBA{0xff, 0x00, 0x00, 0xff}

type Renderer struct {
	cm *component.Manager
}

func NewRenderer(cm *component.Manager) *Renderer {
	return &Renderer{
		cm: cm,
	}
}

func (r *Renderer) Draw(dst *ebiten.Image) error {
	// Draw all entities
	for _, g := range r.cm.List(component.PositionType) {
		p, ok := g.Data.(*component.Position)
		if !ok {
			return fmt.Errorf("could not convert to position: %v", g.Data)
		}
		vector.DrawFilledRect(dst, float32(p.X), float32(p.Y), 16, 16, red, true)
	}
	return nil
}
