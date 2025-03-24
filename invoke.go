package beherit

import (
	"encoding/json"
	"fmt"

	"github.com/dwethmar/beherit/command"
)

const (
	GameCreatedTrigger string = "create.game"
	UpdateTrigger      string = "update"
)

type TriggerCommand struct {
	Command string         `yaml:"command"`
	Params  map[string]any `yaml:"params"`
}

type Config struct {
	Triggers map[string][]*TriggerCommand `yaml:"triggers"`
}

type Invoker struct {
	config         *Config
	commandBus     *command.Bus
	commandFactory *command.Factory
	ec             *ExpressionCompiler
}

func NewInvoker(commandBus *command.Bus, commandFactory *command.Factory) *Invoker {
	return &Invoker{
		config:         nil,
		commandBus:     commandBus,
		commandFactory: commandFactory,
		ec:             NewExpressionCompiler(),
	}
}

func (i *Invoker) Config(config Config) error {
	i.config = &config
	for _, triggers := range config.Triggers {
		for _, trigger := range triggers {
			p, err := i.ec.Compile(trigger.Params)
			if err != nil {
				return err
			}
			trigger.Params = p
		}
	}
	return nil
}

func (i *Invoker) Invoke(t string) error {
	triggers, ok := i.config.Triggers[t]
	if !ok {
		return nil
	}
	for _, c := range triggers {
		command, err := i.commandFactory.Create(c.Command)
		if err != nil {
			return fmt.Errorf("could not create command: %w", err)
		}
		params, err := i.ec.Run(c.Params, nil)
		if err != nil {
			return fmt.Errorf("could not run expression: %w", err)
		}
		// Convert map to JSON
		jsonData, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("could not marshal params: %w", err)
		}
		// Convert JSON to struct
		if err = json.Unmarshal(jsonData, command); err != nil {
			return fmt.Errorf("could not unmarshal params: %w", err)
		}
		if err = i.commandBus.Emit(command); err != nil {
			return err
		}
	}
	return nil
}
