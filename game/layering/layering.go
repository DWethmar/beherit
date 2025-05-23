package layering

import (
	"fmt"

	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
	"github.com/dwethmar/beherit/entity"
	"github.com/dwethmar/beherit/game"
)

func New(cm *component.Manager) command.HandlerFunc {
	return func(cmd *command.Command) error {
		switch c := cmd.Data.(type) {
		case *game.SortGraphicsCommand:
			return SortGraphics(cm, c.GraphicComponents)
		default:
			return fmt.Errorf("unknown command: %T", c)
		}
	}
}

func SortGraphics(cm *component.Manager, graphics []component.Component) error {

	positions := make(map[entity.Entity]int)
	for _, g := range graphics {
		p, ok := cm.FirstByEntity(g.Entity, component.PositionType)
		if !ok {
			return fmt.Errorf("could not find position component for entity %d", g.Entity)
		}
		pData, ok := p.Data.(*component.Position)
		if !ok {
			return fmt.Errorf("expected component type: %T, want: %T", p.Data, &component.Position{})
		}
		positions[g.Entity] = pData.Y
	}

	// Sort the graphics by their Y position
	for i := 0; i < len(graphics)-1; i++ {
		for j := i + 1; j < len(graphics); j++ {
			if positions[graphics[i].Entity] < positions[graphics[j].Entity] {
				graphics[i], graphics[j] = graphics[j], graphics[i]
			}
		}
	}

	// Update the graphics components with the new order
	for i, g := range graphics {
		g.Data.(*component.Graphic).Layer = i
		if err := cm.Update(g); err != nil {
			return fmt.Errorf("could not update graphic component: %w", err)
		}
	}

	return nil
}
