package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"regexp"

	"github.com/dwethmar/beherit"
	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
	"github.com/dwethmar/beherit/entity"
	"github.com/dwethmar/beherit/system/blueprint"
	"github.com/dwethmar/beherit/system/follow"
	"github.com/dwethmar/beherit/system/render"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var MatchByEndingDollarRegex = regexp.MustCompile(`\$$`)
var configFileFlag = flag.String("config", "config.yaml", "config file")

type env struct {
	logger           *slog.Logger
	entityManager    *entity.Manager
	componentManager *component.Manager
}

func (e *env) Create() map[string]any {
	return map[string]any{
		"CreateEntity": e.entityManager.Create,
		"ListComponents": func(t string) []component.Component {
			r := e.componentManager.List(component.Type(t))
			return r
		},
		"KeyPressed": func(key string) bool {
			k := ebiten.Key(0)
			if err := k.UnmarshalText([]byte(key)); err != nil {
				e.logger.Error("could not parse key", slog.String("key", key))
				return false
			}
			return inpututil.KeyPressDuration(k) > 0
		},
	}
}

func main() {
	flag.Parse()
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(logHandler)
	if err := run(context.Background(), logger); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	em := entity.NewManager(1)
	cm := component.NewManager(1)
	cf := component.NewFactory()
	cf.Register(component.NewGraphicComponent)
	cf.Register(component.NewPositionComponent)
	cf.Register(component.NewFollowComponent)

	commandBus := command.NewBus(logger)
	commandFactory := command.NewFactory()

	blueprintSystem := blueprint.New(logger, cf, cm)
	blueprintSystem.Attach(commandFactory, commandBus)

	followSystem := follow.New(logger, cf, cm)
	followSystem.Attach(commandFactory, commandBus)

	renderSystem := render.New(logger, cm)
	renderSystem.Attach(commandFactory, commandBus)

	invoker := beherit.NewInvoker(beherit.InvokerOptions{
		Logger:         logger,
		EntityManager:  em,
		CommandBus:     commandBus,
		CommandFactory: commandFactory,
		Env:            &env{logger: logger, entityManager: em, componentManager: cm},
		ExprComp: beherit.NewExpressionCompiler(&beherit.RegexExprMatcher{
			KeyRegex: MatchByEndingDollarRegex,
		}),
	})

	g := beherit.NewGame(beherit.Options{
		Logger:     logger,
		ConfigFile: *configFileFlag,
		Invoker:    invoker,
	})

	go g.WatchConfig(ctx)

	if err := g.Run(); err != nil {
		return fmt.Errorf("could not run game: %w", err)
	}

	return nil
}
