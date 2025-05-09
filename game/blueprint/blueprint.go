package blueprint

import (
	"fmt"
	"log/slog"

	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
	"github.com/dwethmar/beherit/entity"
	"github.com/dwethmar/beherit/game"
	"github.com/dwethmar/beherit/mapstruct"
)

type System struct {
	logger *slog.Logger
	em     *entity.Manager
	cf     *component.Factory
	cm     *component.Manager
}

func New(logger *slog.Logger, cf *component.Factory, cm *component.Manager) *System {
	return &System{
		logger: logger,
		cf:     cf,
		cm:     cm,
	}
}

func (s *System) Handle(c *command.Command) error {
	switch c := c.Data.(type) {
	case *game.CreateEntity:
		return s.createEntity(c)
	case *game.UpdateEntity:
		return s.updateEntity(c)
	default:
		return fmt.Errorf("unknown command: %T", c)
	}
}

func (s *System) createEntity(cmd *game.CreateEntity) error {
	for _, n := range cmd.Components {
		// check if data is an map. This happens if the command is from a trigger
		if data, ok := n.Data.(map[string]any); ok {
			c, err := s.cf.Create(n.Type)
			if err != nil {
				return fmt.Errorf("failed to create component of type %s: %w", n.Type, err)
			}
			if err := mapstruct.To(data, c.Data); err != nil {
				return fmt.Errorf("failed to map data to component: %w", err)
			}
			n.Data = c.Data
		}
		if err := s.cm.Add(*n); err != nil {
			return fmt.Errorf("failed to add component: %w", err)
		}
		s.logger.Info("created component", slog.Int("entity", int(n.Entity)), slog.String("type", string(n.Type)))
	}
	return nil
}

func (s *System) updateEntity(cmd *game.UpdateEntity) error {
	for _, n := range cmd.Components {
		c, ok := s.cm.FirstByEntity(n.Entity, n.Type)
		if !ok {
			return fmt.Errorf("could not find component %s for entity %d", n.Type, n.Entity)
		}
		// check if data is an map. This happens if the command is from a trigger
		if data, ok := n.Data.(map[string]any); ok {
			if err := mapstruct.To(data, c.Data); err != nil {
				return fmt.Errorf("failed to map data to component: %w", err)
			}
			n.Data = c.Data
		}
		if err := s.cm.Update(*n); err != nil {
			return fmt.Errorf("failed to update component: %w", err)
		}
		s.logger.Info("updated component", slog.Int("entity", int(n.Entity)), slog.String("type", string(n.Type)))
	}
	return nil
}
