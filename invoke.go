package beherit

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"regexp"
	"sync"

	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/entity"
	"github.com/dwethmar/beherit/mapcopier"
	"github.com/dwethmar/beherit/mapquery"
	"github.com/dwethmar/beherit/mapstruct"
	"github.com/expr-lang/expr"
)

type RegexExprMatcher struct {
	KeyRegex *regexp.Regexp
}

func (m RegexExprMatcher) Match(key string) bool {
	return m.KeyRegex.MatchString(key)
}

type Invoker struct {
	logger         *slog.Logger
	configMux      sync.RWMutex
	triggers       map[string][]*TriggerCommand
	entityManager  *entity.Manager
	commandBus     *command.Bus
	commandFactory *command.Factory
	newEnv         func() map[string]any
	ec             *ExpressionCompiler
}

type InvokerOptions struct {
	Logger         *slog.Logger
	EntityManager  *entity.Manager
	CommandBus     *command.Bus
	CommandFactory *command.Factory
	NewEnv         func() map[string]any
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
		newEnv:         opt.NewEnv,
		ec:             opt.ExprComp,
	}
}

func (i *Invoker) SetTriggers(triggers map[string][]*TriggerCommand) error {
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
				return fmt.Errorf("could not compile vars: %w", err)
			}
			trigger.Vars = in

			if trigger.Set != nil {
				for _, s := range trigger.Set {
					v, err := i.ec.Compile(s.ValueStr)
					if err != nil {
						return fmt.Errorf("could not compile set value: %w", err)
					}
					s.Value = v
				}
			}

			if trigger.ConditionStr != "" {
				cond, err := i.ec.Compile(trigger.ConditionStr)
				if err != nil {
					return fmt.Errorf("could not compile condition: %w", err)
				}
				trigger.condition = cond
			}

			p, err := i.ec.CompileMap(trigger.Mapping)
			if err != nil {
				return fmt.Errorf("could not compile params: %w", err)
			}
			trigger.Mapping = p
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
	env := i.newEnv()
	for j, c := range triggers {
		localEnv := map[string]any{}
		maps.Copy(localEnv, mapcopier.Copy(globalEnv))
		maps.Copy(localEnv, mapcopier.Copy(env))
		if err := i.trigger(c, map[string]any{"env": localEnv}); err != nil {
			return fmt.Errorf("failed to trigger command %d: %w", j, err)
		}
	}
	return nil
}

func (i *Invoker) trigger(c *TriggerCommand, env map[string]any) error {
	if c.Vars != nil {
		vars, err := i.ec.Run(c.Vars, env)
		if err != nil {
			return fmt.Errorf("failed to run vars expression for command %q: %w", c.Command, err)
		}
		env["vars"] = vars
	}

	// set clause
	for _, s := range c.Set {
		if s.Value == nil {
			return fmt.Errorf("set value is nil for command %q", c.Command)
		}
		val, err := expr.Run(s.Value, env)
		if err != nil {
			return fmt.Errorf("failed to run set value expression for command %q: %w", c.Command, err)
		}
		if err := mapquery.Set(s.Path, val, &env); err != nil {
			return fmt.Errorf("failed to set value for command %q: %w", c.Command, err)
		}
	}

	// check if condition is met
	if c.condition != nil {
		cond, err := expr.Run(c.condition, env)
		if err != nil {
			return fmt.Errorf("failed to run condition expression for command %q: %w", c.Command, err)
		}
		r, ok := cond.(bool)
		if !ok {
			return fmt.Errorf("condition result must be a boolean")
		}
		if !r {
			return nil
		}
	}

	command, err := i.commandFactory.Create(c.Command)
	if err != nil {
		return fmt.Errorf("failed to create command %q: %w", c.Command, err)
	}

	data, err := i.ec.Run(c.Mapping, env)
	if err != nil {
		return fmt.Errorf("failed to run param expression for command %q: %w", c.Command, err)
	}

	if err = mapstruct.To(data, command.Data); err != nil {
		return fmt.Errorf("failed to map params to command %q: %w", c.Command, err)
	}

	if err = i.commandBus.Emit(command); err != nil {
		return err
	}
	return nil
}
