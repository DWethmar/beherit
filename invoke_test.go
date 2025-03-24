package beherit_test

import (
	"log/slog"
	"testing"

	"github.com/dwethmar/beherit"
	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
	"github.com/dwethmar/beherit/entity"
)

type MockCommand struct {
	EntityID   entity.Entity
	Components map[component.Type]map[string]interface{}
}

func (c MockCommand) Type() string { return "testcommand" }

func TestInvoker_Invoke(t *testing.T) {
	t.Run("invoke", func(t *testing.T) {
		calledHandler := 0
		commandFactory := command.NewFactory()
		commandFactory.Register(func() command.Command {
			return &MockCommand{}
		})
		commandBus := command.NewBus(slog.Default())
		commandBus.RegisterHandler(&MockCommand{}, command.HandlerFunc(func(_ command.Command) error {
			calledHandler++
			return nil
		}))
		invoker := beherit.NewInvoker(commandBus, commandFactory)
		invoker.Config(
			beherit.Config{
				Triggers: map[string][]*beherit.TriggerCommand{
					beherit.GameCreatedTrigger: {
						{
							Command: "testcommand",
							Params: map[string]any{
								"EntityID": 1,
								"Components": map[string]any{
									"position": map[string]any{
										"X": 1,
										"Y": 2,
									},
								},
							},
						},
						{
							Command: "testcommand",
							Params: map[string]any{
								"EntityID": 2,
								"Components": map[string]any{
									"position": map[string]any{
										"X": 1,
										"Y": 2,
									},
								},
							},
						},
					},
				},
			},
		)

		if err := invoker.Invoke(beherit.GameCreatedTrigger); err != nil {
			t.Error(err)
		}

		if calledHandler != 2 {
			t.Errorf("expected handler to be called 2 times, got %d", calledHandler)
		}
	})
}
