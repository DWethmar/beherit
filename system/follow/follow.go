package follow

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

func (s *System) Attach(commandFactory *command.Factory, commandBus *command.Bus) {
	commandFactory.Register(NewSetTargetCommand)
	commandFactory.Register(NewMoveTowardsTarget)
	commandBus.RegisterHandler(NewSetTargetCommand(), s)
	commandBus.RegisterHandler(NewMoveTowardsTarget(), s)
}

func (s *System) Handle(c *command.Command) error {
	switch c := c.Data.(type) {
	case *SetTarget:
		return s.Target(c)
	case *MoveTowardsTarget:
		return s.MoveTowardsTarget(c)
	default:
		return fmt.Errorf("unknown command: %T", c)
	}
}

func (s *System) Target(cmd *SetTarget) error {
	s.logger.Info("follow", slog.Int("x", cmd.X), slog.Int("y", cmd.Y), slog.Any("entity", cmd.Entities))

	for _, e := range cmd.Entities {
		var f component.Component
		if c := s.cm.ListByEntity(e, component.FollowType); len(c) > 0 {
			f = c[0]
		} else {
			return fmt.Errorf("entity %d has no follow component", e)
		}

		fData, ok := f.Data.(*component.Follow)
		if !ok {
			return fmt.Errorf("could not cast follow component data")
		}

		fData.X = cmd.X
		fData.Y = cmd.Y

		if err := s.cm.Update(f); err != nil {
			return fmt.Errorf("could not update follow component: %w", err)
		}
	}

	return nil
}

func (s *System) MoveTowardsTarget(cmd *MoveTowardsTarget) error {
	s.logger.Info("move towards target", slog.Any("entity", cmd.Entities))

	for _, e := range cmd.Entities {
		var f component.Component
		if c := s.cm.ListByEntity(e, component.FollowType); len(c) > 0 {
			f = c[0]
		} else {
			return fmt.Errorf("entity %s has no %v component", component.FollowType, e)
		}

		var p component.Component
		if c := s.cm.ListByEntity(e, component.PositionType); len(c) > 0 {
			p = c[0]
		} else {
			return fmt.Errorf("entity %s has no %v component", component.PositionType, e)
		}

		fData, ok := f.Data.(*component.Follow)
		if !ok {
			return fmt.Errorf("could not cast follow component data")
		}

		pData, ok := p.Data.(*component.Position)
		if !ok {
			return fmt.Errorf("could not cast position component data")
		}

		if pData.X < fData.X {
			pData.X++
		} else if pData.X > fData.X {
			pData.X--
		}

		if pData.Y < fData.Y {
			pData.Y++
		} else if pData.Y > fData.Y {
			pData.Y--
		}

		if err := s.cm.Update(p); err != nil {
			return fmt.Errorf("could not update follow component: %w", err)
		}
	}

	return nil
}
