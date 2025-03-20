package component

import (
	"github.com/dwethmar/beherit"
	"github.com/dwethmar/beherit/cmd/mygame/sprite"
	"github.com/dwethmar/beherit/direction"
)

const (
	ComponentTypePosition beherit.ComponentType = "position"
	ComponentTypeSkeleton beherit.ComponentType = "skeleton"
)

type Position struct {
	X, Y int
}

type Skeleton struct {
	Facing direction.Direction
}

type Sprite struct {
	sprite.Sprite
}
