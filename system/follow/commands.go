package follow

import (
	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/entity"
)

const (
	SetTargetCommandType  = "follow.set-target"
	MoveTowardsTargetType = "follow.move-towards-target"
)

type SetTarget struct {
	Entities []entity.Entity `json:"entities"`
	X        int             `json:"x"`
	Y        int             `json:"y"`
}

func NewSetTargetCommand() *command.Command {
	return &command.Command{
		Type: SetTargetCommandType,
		Data: &SetTarget{},
	}
}

type MoveTowardsTarget struct {
	Entities []entity.Entity `json:"entities"`
}

func NewMoveTowardsTarget() *command.Command {
	return &command.Command{
		Type: MoveTowardsTargetType,
		Data: &MoveTowardsTarget{},
	}
}
