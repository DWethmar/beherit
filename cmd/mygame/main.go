package main

import (
	"log/slog"

	"github.com/dwethmar/beherit"
	"github.com/dwethmar/beherit/cmd/mygame/sprite"
	"github.com/dwethmar/beherit/cmd/mygame/system/skeleton"
)

func main() {
	logger := slog.Default()

	logger.Info("loading spritesheets")
	sheet, err := sprite.NewSpritesheet()
	if err != nil {
		panic(err)
	}

	sprites := sprite.Load(sheet)

	if err = beherit.NewGame([]beherit.System{
		skeleton.New(logger, sprites),
	}).Run(); err != nil {
		panic(err)
	}
}
