package beherit

import (
	"encoding/json"
	"fmt"

	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/entity"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type TriggerCommand struct {
	Command     string         `yaml:"command"`
	Include     map[string]any `yaml:"include"`
	Condition   string         `yaml:"condition"`
	conditionPr *vm.Program    `yaml:"-"`
	Params      map[string]any `yaml:"params"`
}

type Config struct {
	Triggers map[string][]*TriggerCommand `yaml:"triggers"`
}

type Env interface {
	Create() map[string]any
}

type Invoker struct {
	config         *Config
	entityManager  *entity.Manager
	commandBus     *command.Bus
	commandFactory *command.Factory
	env            Env
	ec             *ExpressionCompiler
}

func NewInvoker(entityManager *entity.Manager, commandBus *command.Bus, commandFactory *command.Factory, env Env) *Invoker {
	return &Invoker{
		config:         nil,
		entityManager:  entityManager,
		commandBus:     commandBus,
		commandFactory: commandFactory,
		env:            env,
		ec:             NewExpressionCompiler(),
	}
}

func (i *Invoker) Config(config Config) error {
	i.config = &config
	for _, triggers := range config.Triggers {
		for _, trigger := range triggers {
			in, err := i.ec.CompileMap(trigger.Include)
			if err != nil {
				return fmt.Errorf("could not compile expression: %w", err)
			}
			trigger.Include = in

			if trigger.Condition != "" {
				cond, err := i.ec.Compile(trigger.Condition)
				if err != nil {
					return fmt.Errorf("could not compile expression: %w", err)
				}
				trigger.conditionPr = cond
			}

			p, err := i.ec.CompileMap(trigger.Params)
			if err != nil {
				return fmt.Errorf("could not compile expression: %w", err)
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
		env := i.env.Create()

		if c.Include != nil {
			include, err := i.ec.Run(c.Include, env)
			if err != nil {
				return fmt.Errorf("could not run expression: %w", err)
			}
			// Include env
			for k, v := range include {
				env[k] = v
			}
		}

		if c.conditionPr != nil {
			cond, err := expr.Run(c.conditionPr, env)
			if err != nil {
				return fmt.Errorf("could not run expression: %w", err)
			}
			r, ok := cond.(bool)
			if !ok {
				return fmt.Errorf("condition result must be a boolean")
			}
			if !r {
				continue
			}
		}
		params, err := i.ec.Run(c.Params, env)
		if err != nil {
			return fmt.Errorf("could not run expression: %w", err)
		}
		// Convert map to JSON
		jsonData, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("could not marshal params: %w", err)
		}
		command, err := i.commandFactory.Create(c.Command)
		if err != nil {
			return fmt.Errorf("could not create command: %w", err)
		}
		// Convert JSON to struct
		if err = json.Unmarshal(jsonData, command.Data); err != nil {
			return fmt.Errorf("could not unmarshal params: %w", err)
		}
		if err = i.commandBus.Emit(command); err != nil {
			return err
		}
	}
	return nil
}
