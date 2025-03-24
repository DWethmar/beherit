package render

import (
	"github.com/dwethmar/beherit/component"
	"github.com/hajimehoshi/ebiten/v2"
)

type Renderer struct {
	cm *component.Manager
}

func NewRenderer(cm *component.Manager) *Renderer {
	return &Renderer{}
}

func (r *Renderer) Draw(s *ebiten.Image) error {
	// Draw all entities
	// for _, g := range r.cm.List(component.TypeGraphic) {
	// 	// Draw the entity
	// 	// g.Draw(s)
	// }
	return nil
}
