package main

import (
	"flag"
	"log/slog"

	"github.com/dwethmar/beherit"
	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
	"github.com/dwethmar/beherit/entity"
	"github.com/dwethmar/beherit/render"
	"github.com/dwethmar/beherit/system/blueprint"
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
	flag.Parse()
	logger := slog.Default()

	logger.Info("loading spritesheets")
	// sheet, err := sprite.NewSpritesheet()
	// if err != nil {
	// 	panic(err)
	// }

	// sprites := sprite.Load(sheet)
	// eventBus := event.NewBus()
	em := entity.NewManager(1)
	cm := component.NewManager(1)
	cf := component.NewFactory()
	cf.Register(component.NewGraphicComponent)
	cf.Register(component.NewPositionComponent)

	commandBus := command.NewBus(logger)
	commandFactory := command.NewFactory()

	commandFactory.Register(blueprint.NewCreateEntityCommand)
	commandFactory.Register(blueprint.NewUpdateEntityCommand)

	blueprintSystem := blueprint.New(logger, em, cf, cm)
	commandBus.RegisterHandler(blueprint.NewCreateEntityCommand(), blueprintSystem)
	commandBus.RegisterHandler(blueprint.NewUpdateEntityCommand(), blueprintSystem)

	if err := beherit.NewGame(beherit.Options{
		Logger:     logger,
		ConfigFile: *configFileFlag,
		Invoker: beherit.NewInvoker(em, commandBus, commandFactory, &env{
			logger:           logger,
			entityManager:    em,
			componentManager: cm,
		}),
		Renderer: []beherit.Renderer{render.NewRenderer(cm)},
	}).Run(); err != nil {
		panic(err)
	}
}
