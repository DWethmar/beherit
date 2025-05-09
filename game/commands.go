package game

import (
	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
	"github.com/dwethmar/beherit/entity"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	CreateEntityCommandType     = "blueprint.create.entity"
	UpdateEntityCommandType     = "blueprint.update.entity"
	RemoveComponentsCommandType = "blueprint.remove.components"

	SetTargetCommandType  = "follow.set-target"
	MoveTowardsTargetType = "follow.move-towards-target"

	RenderCommandType = "render"
)

type CreateEntity struct {
	Components []*component.Component `json:"components"`
}

func NewCreateEntityCommand() *command.Command {
	return &command.Command{
		Type: CreateEntityCommandType,
		Data: &CreateEntity{
			Components: []*component.Component{},
		},
	}
}

type UpdateEntity struct {
	Components []*component.Component `json:"components"`
}

func NewUpdateEntityCommand() *command.Command {
	return &command.Command{
		Type: UpdateEntityCommandType,
		Data: &UpdateEntity{
			Components: []*component.Component{},
		},
	}
}

type RemoveComponents struct {
	Entity       entity.Entity
	ComponentIDs []uint
}

func NewRemoveComponentsCommand() *command.Command {
	return &command.Command{
		Type: RemoveComponentsCommandType,
		Data: &RemoveComponents{},
	}
}

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

type RenderCheckerBoardCommand struct {
	Screen *ebiten.Image `json:"screen"`
}

func NewRenderCheckerBoardCommand() *command.Command {
	return &command.Command{
		Type: RenderCommandType,
		Data: &RenderCheckerBoardCommand{
			Screen: nil,
		},
	}
}

type RenderComponentsCommand struct {
	Screen     *ebiten.Image         `json:"screen"`
	Components []component.Component `json:"components"`
}

func NewRenderCommand() *command.Command {
	return &command.Command{
		Type: RenderCommandType,
		Data: &RenderComponentsCommand{
			Screen:     nil,
			Components: []component.Component{},
		},
	}
}
