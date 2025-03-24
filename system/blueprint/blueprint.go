package blueprint

import (
	"fmt"
	"log/slog"

	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
	"github.com/dwethmar/beherit/entity"
)

const (
	CreateEntityCommandType = "create.blueprint.entity"
)

type CreateEntityCommand struct {
	EntityID   entity.Entity            `json:"entity_id"`
	Components []map[string]interface{} `json:"components"`
}

func (c CreateEntityCommand) Type() string { return CreateEntityCommandType }

func NewCreateEntityCommand() command.Command {
	return &CreateEntityCommand{}
}

type System struct {
	logger *slog.Logger
	em     *entity.Manager
	cm     *component.Manager
}

func New(logger *slog.Logger, em *entity.Manager, cm *component.Manager) *System {
	return &System{
		logger: logger,
		em:     em,
		cm:     cm,
	}
}

func (s *System) Handle(c command.Command) error {
	switch c := c.(type) {
	case *CreateEntityCommand:
		s.logger.Info("handle command", slog.Any("command", c))
		return nil
	default:
		return fmt.Errorf("unknown command: %T", c)
	}
}
