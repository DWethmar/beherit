package blueprint

import (
	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
)

const (
	CreateEntityCommandType = "blueprint.create.entity"
	UpdateEntityCommandType = "blueprint.update.entity"
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
