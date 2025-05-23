package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"regexp"

	"github.com/dwethmar/beherit"
	"github.com/dwethmar/beherit/cmd/mygame/sprite"
	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
	"github.com/dwethmar/beherit/entity"
	"github.com/dwethmar/beherit/game"
	"github.com/dwethmar/beherit/game/blueprint"
	"github.com/dwethmar/beherit/game/follow"
	"github.com/dwethmar/beherit/game/layering"
	"github.com/dwethmar/beherit/game/render"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed images
var imagesFS embed.FS

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

	// Register the command handlers

	// blueprint system
	blueprintSystem := blueprint.New(logger, cf, cm)
	commandFactory.Register(game.NewCreateEntityCommand)
	commandBus.RegisterHandler(game.NewCreateEntityCommand(), blueprintSystem)

	commandFactory.Register(game.NewUpdateEntityCommand)
	commandBus.RegisterHandler(game.NewUpdateEntityCommand(), blueprintSystem)

	commandFactory.Register(game.NewRemoveComponentsCommand)
	commandBus.RegisterHandler(game.NewRemoveComponentsCommand(), blueprintSystem)

	commandFactory.Register(game.NewSortGraphicsCommand)
	commandBus.RegisterHandler(game.NewSortGraphicsCommand(), layering.New(cm))

	// follow system
	followSystem := follow.New(logger, cf, cm)
	commandFactory.Register(game.NewSetTargetCommand)
	commandBus.RegisterHandler(game.NewSetTargetCommand(), followSystem)

	commandFactory.Register(game.NewMoveTowardsTarget)
	commandBus.RegisterHandler(game.NewMoveTowardsTarget(), followSystem)

	// render system
	skeletonDeathImg, err := sprite.LoadPng(imagesFS, "images/skeleton_death.png")
	if err != nil {
		return fmt.Errorf("failed to load images/skeleton_death.png: %w", err)
	}

	skeletonMoveImg, err := sprite.LoadPng(imagesFS, "images/skeleton_move.png")
	if err != nil {
		return fmt.Errorf("failed to load images/skeleton_move.png: %w", err)
	}

	skeletonDeathSprites := sprite.LoadSkeletonDeath(ebiten.NewImageFromImage(skeletonDeathImg))
	skeletonMoveSprites := sprite.LoadSkeletonMove(ebiten.NewImageFromImage(skeletonMoveImg))

	var sprites map[string]*sprite.Sprite = make(map[string]*sprite.Sprite)
	maps.Copy(sprites, skeletonDeathSprites)
	maps.Copy(sprites, skeletonMoveSprites)

	renderSystem := render.New(logger, cm, sprites)
	commandFactory.Register(game.NewRenderCommand)
	commandBus.RegisterHandler(game.NewRenderCommand(), renderSystem)

	env := &env{logger: logger, entityManager: em, componentManager: cm}

	invoker := beherit.NewInvoker(beherit.InvokerOptions{
		Logger:         logger,
		EntityManager:  em,
		CommandBus:     commandBus,
		CommandFactory: commandFactory,
		NewEnv:         env.Create,
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
