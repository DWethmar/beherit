package command

import (
	"fmt"
	"log/slog"
)

// Command is a command.
type Command interface {
	Type() string
}

type Handler interface {
	Handle(Command) error
}

type HandlerFunc func(Command) error

func (f HandlerFunc) Handle(c Command) error {
	return f(c)
}

type Bus struct {
	logger   *slog.Logger
	handlers map[string]Handler
}

func NewBus(logger *slog.Logger) *Bus {
	return &Bus{
		logger:   logger,
		handlers: make(map[string]Handler),
	}
}

func (cb *Bus) RegisterHandler(c Command, handler Handler) {
	cb.handlers[c.Type()] = handler
}

func (cb *Bus) Emit(c Command) error {
	handler, ok := cb.handlers[c.Type()]
	if !ok {
		return fmt.Errorf("no handler for trigger %s", c.Type())
	}
	return handler.Handle(c)
}

type Factory struct {
	factory map[string]func() Command
}

func NewFactory() *Factory {
	return &Factory{
		factory: make(map[string]func() Command),
	}
}

func (f *Factory) Register(v func() Command) {
	f.factory[v().Type()] = v
}

func (f *Factory) Create(commandType string) (Command, error) {
	factory, ok := f.factory[commandType]
	if !ok {
		return nil, fmt.Errorf("no factory found for command type %s", commandType)
	}
	return factory(), nil
}
