package main

import (
	"flag"
	"log/slog"

	"github.com/dwethmar/beherit"
	"github.com/dwethmar/beherit/command"
	"github.com/dwethmar/beherit/component"
	"github.com/dwethmar/beherit/component/graphic"
	"github.com/dwethmar/beherit/component/position"
	"github.com/dwethmar/beherit/entity"
	"github.com/dwethmar/beherit/system/blueprint"
)

var configFileFlag = flag.String("config", "config.yaml", "config file")

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
	em := entity.NewManager(0)
	cm := component.NewManager()
	cm.Register(position.NewComponent)
	cm.Register(graphic.NewComponent)

	commandBus := command.NewBus(logger)
	commandFactory := command.NewFactory()

	commandFactory.Register(blueprint.NewCreateEntityCommand)

	blueprintSystem := blueprint.New(logger, em, cm)
	commandBus.RegisterHandler(&blueprint.CreateEntityCommand{}, blueprintSystem)

	if err := beherit.NewGame(beherit.Options{
		Logger:     logger,
		ConfigFile: *configFileFlag,
		Invoker:    beherit.NewInvoker(commandBus, commandFactory),
	}).Run(); err != nil {
		panic(err)
	}
}
