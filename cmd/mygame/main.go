package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/dwethmar/beherit"
	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
	"github.com/dwethmar/beherit/entity"
	"github.com/dwethmar/beherit/render"
	"github.com/dwethmar/beherit/system/blueprint"
	"github.com/dwethmar/beherit/system/follow"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	flag.Parse()
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(logHandler)

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

	g := beherit.NewGame(beherit.Options{
		Logger:     logger,
		ConfigFile: *configFileFlag,
		Invoker: beherit.NewInvoker(logger, em, commandBus, commandFactory, &env{
			logger:           logger,
			entityManager:    em,
			componentManager: cm,
		}),
		Renderer: []beherit.Renderer{render.NewRenderer(cm)},
	})

	go g.WatchConfig(ctx)

	if err := g.Run(); err != nil {
		panic(err)
	}
}
