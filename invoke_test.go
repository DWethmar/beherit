package beherit_test

import (
	"log/slog"
	"testing"

	"github.com/dwethmar/beherit"
	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/entity"
)

type Component struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
}

type MockCommand struct {
	EntityID   entity.Entity
	Components []Component
}

func (c MockCommand) Type() string { return "testcommand" }

type MockEnv struct {
	nextId int
}

func (e MockEnv) Create() map[string]any {
	return map[string]any{
		"CreateEntity": func() entity.Entity {
			e.nextId++
			return entity.Entity(e.nextId)
		},
	}
}

func TestInvoker_Invoke(t *testing.T) {
	type MockCommand struct {
		EntityID   entity.Entity
		Components []map[string]interface{}
	}
	t.Run("invoke", func(t *testing.T) {
		entityManager := entity.NewManager(1)
		calledHandler := 0
		commandFactory := command.NewFactory()
		createCMD := func() *command.Command {
			return &command.Command{
				Type: "testcommand",
				Data: &MockCommand{},
			}
		}
		commandFactory.Register(createCMD)
		commandBus := command.NewBus(slog.Default())
		commandBus.RegisterHandler(createCMD(), command.HandlerFunc(func(_ *command.Command) error {
			calledHandler++
			return nil
		}))
		invoker := beherit.NewInvoker(beherit.InvokerOptions{
			Logger:         slog.Default(),
			EntityManager:  entityManager,
			CommandBus:     commandBus,
			CommandFactory: commandFactory,
			Env:            MockEnv{},
			ExprComp: beherit.NewExpressionCompiler(&beherit.RegexExprMatcher{
				KeyRegex: MatchByEndingDollarRegex,
			}),
		})
		invoker.Config(
			beherit.Config{
				Triggers: map[string][]*beherit.TriggerCommand{
					beherit.GameCreatedTrigger: {
						{
							Command: "testcommand",
							Params: map[string]any{
								"EntityID$": "CreateEntity()",
								"Components": []map[string]any{
									{
										"type": "position",
										"params": map[string]any{
											"X$": "1 + 2",
											"Y":  2,
										},
									},
								},
							},
						},
						{
							Command: "testcommand",
							Params: map[string]any{
								"EntityID": 1,
								"Components": []map[string]any{
									{
										"type": "position",
										"params": map[string]any{
											"X": 1,
											"Y": 2,
										},
									},
								},
							},
						},
					},
				},
			},
		)

		if err := invoker.Invoke(beherit.GameCreatedTrigger, map[string]any{}); err != nil {
			t.Error(err)
		}

		if calledHandler != 2 {
			t.Errorf("expected handler to be called 2 times, got %d", calledHandler)
		}
	})
}
