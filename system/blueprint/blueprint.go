package blueprint

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
	"github.com/dwethmar/beherit/entity"
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

func (s *System) Attach(commandFactory *command.Factory, commandBus *command.Bus) {
	commandFactory.Register(NewCreateEntityCommand)
	commandFactory.Register(NewUpdateEntityCommand)
	commandFactory.Register(NewRemoveComponentsCommand)
	commandBus.RegisterHandler(NewCreateEntityCommand(), s)
	commandBus.RegisterHandler(NewUpdateEntityCommand(), s)
}

func (s *System) Handle(c *command.Command) error {
	switch c := c.Data.(type) {
	case *CreateEntity:
		return s.createEntity(c)
	case *UpdateEntity:
		return s.updateEntity(c)
	default:
		return fmt.Errorf("unknown command: %T", c)
	}
}

func (s *System) createEntity(cmd *CreateEntity) error {
	for _, n := range cmd.Components {
		c, err := s.cf.Create(component.Type(n.Type))
		if err != nil {
			return fmt.Errorf("create component: %w", err)
		}
		b, err := json.Marshal(n.Data)
		if err != nil {
			return fmt.Errorf("marshal params: %w", err)
		}
		if err := json.Unmarshal(b, c.Data); err != nil {
			return fmt.Errorf("unmarshal params: %w", err)
		}
		n.Data = c.Data
		if err := s.cm.Add(*n); err != nil {
			return fmt.Errorf("update component: %w", err)
		}
		s.logger.Info("created component", slog.Int("entity", int(c.Entity)), slog.String("type", string(c.Type)))
	}
	return nil
}

func (s *System) updateEntity(cmd *UpdateEntity) error {
	for _, n := range cmd.Components {
		c, err := s.cf.Create(component.Type(n.Type))
		if err != nil {
			return fmt.Errorf("create component: %w", err)
		}
		b, err := json.Marshal(n.Data)
		if err != nil {
			return fmt.Errorf("marshal params: %w", err)
		}
		if err := json.Unmarshal(b, c.Data); err != nil {
			return fmt.Errorf("unmarshal params: %w", err)
		}
		n.Data = c.Data
		if err := s.cm.Update(*n); err != nil {
			return fmt.Errorf("update component: %w", err)
		}
		s.logger.Info("updated component", slog.Int("entity", int(c.Entity)), slog.String("type", string(c.Type)))
	}
	return nil
}
