package beherit

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"sync"

	"maps"

	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/entity"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type RegexExprMatcher struct {
	KeyRegex *regexp.Regexp
}

func (m RegexExprMatcher) Match(key string) bool {
	return m.KeyRegex.MatchString(key)
}

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
	logger         *slog.Logger
	configMux      sync.RWMutex
	config         *Config
	entityManager  *entity.Manager
	commandBus     *command.Bus
	commandFactory *command.Factory
	env            Env
	ec             *ExpressionCompiler
}

type InvokerOptions struct {
	Logger         *slog.Logger
	EntityManager  *entity.Manager
	CommandBus     *command.Bus
	CommandFactory *command.Factory
	Env            Env
	ExprComp       *ExpressionCompiler
}

func NewInvoker(opt InvokerOptions) *Invoker {
	return &Invoker{
		logger:         opt.Logger,
		configMux:      sync.RWMutex{},
		config:         nil,
		entityManager:  opt.EntityManager,
		commandBus:     opt.CommandBus,
		commandFactory: opt.CommandFactory,
		env:            opt.Env,
		ec:             opt.ExprComp,
	}
}

func (i *Invoker) Config(config Config) error {
	i.configMux.Lock()
	defer i.configMux.Unlock()
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

func (i *Invoker) Invoke(t string, globalEnv map[string]any) error {
	i.configMux.RLock()
	defer i.configMux.RUnlock()
	triggers, ok := i.config.Triggers[t]
	if !ok {
		return nil
	}
	for _, c := range triggers {
		env := i.env.Create()
		maps.Copy(env, globalEnv)
		if c.Include != nil {
			include, err := i.ec.Run(c.Include, env)
			if err != nil {
				return fmt.Errorf("could not run expression: %w", err)
			}
			maps.Copy(env, include)
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
		i.logger.Debug("invoking command", slog.String("command", c.Command))

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
