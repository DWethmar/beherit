package blueprint

import (
	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
	"github.com/dwethmar/beherit/entity"
)

const (
	CreateEntityCommandType = "blueprint.create.entity"
	UpdateEntityCommandType = "blueprint.update.entity"
)

const (
	RemoveComponentsCommandType = "blueprint.remove.components"
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
