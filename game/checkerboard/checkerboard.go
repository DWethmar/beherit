package checkerboard

import (
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

func (s *System) HandleCommand(c *command.Command) error {
	return fmt.Errorf("unknown command: %T", c)
}
