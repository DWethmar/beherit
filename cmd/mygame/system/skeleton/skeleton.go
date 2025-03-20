package skeleton

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/dwethmar/beherit"
	"github.com/dwethmar/beherit/cmd/mygame/component"
	"github.com/dwethmar/beherit/cmd/mygame/sprite"
	"github.com/dwethmar/beherit/direction"
	"github.com/hajimehoshi/ebiten/v2"
)

var _ beherit.System = &System{}

type System struct {
	logger  *slog.Logger
	sprites sprite.Sprites
	cm      *beherit.ComponentManager
}

func New(logger *slog.Logger, sprites sprite.Sprites) *System {
	return &System{
		logger:  logger,
		sprites: sprites,
	}
}

func (s *System) Init(em *beherit.EntityManager, cm *beherit.ComponentManager) error {
	s.cm = cm
	s.cm.Register(component.ComponentTypeSkeleton, func() any { return &component.Skeleton{} })
	skeletonCmp, err := s.cm.Create(em.Create(), component.ComponentTypeSkeleton)
	if err != nil {
		return fmt.Errorf("failed to create skeleton component: %w", err)
	}
	skeleton, ok := skeletonCmp.Data.(*component.Skeleton)
	if !ok {
		return errors.New("skeleton component is not a skeleton")
	}
	skeleton.Facing = direction.South
	if err = s.cm.Update(skeletonCmp); err != nil {
		return fmt.Errorf("failed to update skeleton component: %w", err)
	}
	return nil
}

func (s *System) Update() error {
	return nil
}

func (s *System) Draw(_ *ebiten.Image) {
	for _, c := range s.cm.Components(component.ComponentTypeSkeleton) {
		skeleton, ok := c.Data.(*component.Skeleton)
		if !ok {
			s.logger.Info("skeleton component is not a skeleton", slog.String("component", c.Type.String()))
			continue
		}

		s.logger.Info("drawing skeleton", slog.String("facing", skeleton.Facing.String()))
	}
}
