package beherit

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"sync"

	"maps"

	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/entity"
	"github.com/dwethmar/beherit/mapstruct"
	"github.com/expr-lang/expr"
)

type RegexExprMatcher struct {
	KeyRegex *regexp.Regexp
}

func (m RegexExprMatcher) Match(key string) bool {
	return m.KeyRegex.MatchString(key)
}

type Env interface {
	Create() map[string]any
}

type Invoker struct {
	logger         *slog.Logger
	configMux      sync.RWMutex
	triggers       map[string][]*TriggerCommand
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
		triggers:       make(map[string][]*TriggerCommand),
		entityManager:  opt.EntityManager,
		commandBus:     opt.CommandBus,
		commandFactory: opt.CommandFactory,
		env:            opt.Env,
		ec:             opt.ExprComp,
	}
}

func (i *Invoker) Configure(triggers map[string][]*TriggerCommand) error {
	i.configMux.Lock()
	defer i.configMux.Unlock()
	if triggers == nil {
		return errors.New("triggers cannot be nil")
	}
	i.triggers = triggers
	for _, triggers := range triggers {
		for _, trigger := range triggers {
			in, err := i.ec.CompileMap(trigger.Vars)
			if err != nil {
				return fmt.Errorf("could not compile expression: %w", err)
			}
			trigger.Vars = in

			if trigger.ConditionStr != "" {
				cond, err := i.ec.Compile(trigger.ConditionStr)
				if err != nil {
					return fmt.Errorf("could not compile expression: %w", err)
				}
				trigger.condition = cond
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
	triggers, ok := i.triggers[t]
	if !ok {
		return nil
	}
	env := i.env.Create()
	maps.Copy(env, globalEnv)
	for _, c := range triggers {
		if err := i.trigger(c, env); err != nil {
			return err
		}
	}
	return nil
}

func (i *Invoker) trigger(c *TriggerCommand, globalEnv map[string]any) error {
	localEnv := map[string]any{}
	maps.Copy(localEnv, globalEnv)
	if c.Vars != nil {
		vars, err := i.ec.Run(c.Vars, localEnv)
		if err != nil {
			return fmt.Errorf("could not run expression for command '%s': %w", c.Command, err)
		}
		maps.Copy(localEnv, vars)
	}
	// check if condition is met
	if c.condition != nil {
		cond, err := expr.Run(c.condition, localEnv)
		if err != nil {
			return fmt.Errorf("could not run expression for command '%s': %w", c.Command, err)
		}
		r, ok := cond.(bool)
		if !ok {
			return fmt.Errorf("condition result must be a boolean")
		}
		if !r {
			return nil
		}
	}

	params, err := i.ec.Run(c.Params, localEnv)
	if err != nil {
		return fmt.Errorf("could not run expression for command '%s': %w", c.Command, err)
	}

	command, err := i.commandFactory.Create(c.Command)
	if err != nil {
		return fmt.Errorf("could not create command: %w", err)
	}

	if err = mapstruct.To(params, command.Data); err != nil {
		return fmt.Errorf("could not map params to command '%s': %w", c.Command, err)
	}

	if err = i.commandBus.Emit(command); err != nil {
		return err
	}
	return nil
}
