package beherit

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"gopkg.in/yaml.v3"
)

const (
	windowWidth  = 640
	windowHeight = 480
)

type GameState uint

const (
	GameStateInitializing GameState = iota
	GameStateRunning
	GameStatePaused
	GameStateStopped
)

type Renderer interface {
	Draw(screen *ebiten.Image) error
}

type Game struct {
	logger     *slog.Logger
	configFile string
	state      GameState
	Invoker    *Invoker
	renderers  []Renderer
}

type Options struct {
	Logger     *slog.Logger
	ConfigFile string
	Invoker    *Invoker
	Renderer   []Renderer
}

func NewGame(opts Options) *Game {
	return &Game{
		logger:     opts.Logger,
		configFile: opts.ConfigFile,
		state:      GameStateInitializing,
		Invoker:    opts.Invoker,
		renderers:  opts.Renderer,
	}
}

func (g *Game) Update() error {
	switch g.state {
	case GameStateInitializing:
		if err := g.init(); err != nil {
			return err
		}
		if err := g.Invoker.Invoke(GameCreatedTrigger); err != nil {
			return fmt.Errorf("could not invoke game created trigger: %w", err)
		}
		g.state = GameStateRunning
	case GameStateRunning:
		if err := g.Invoker.Invoke(UpdateTrigger); err != nil {
			return fmt.Errorf("could not invoke update trigger: %w", err)
		}
	case GameStatePaused:
	case GameStateStopped:
	}
	return nil
}

func (g *Game) init() error {
	file, err := os.Open(g.configFile)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()
	// parse yaml
	var config Config
	if err = yaml.NewDecoder(file).Decode(&config); err != nil {
		return fmt.Errorf("could not decode yaml: %w", err)
	}
	if err = g.Invoker.Config(config); err != nil {
		return fmt.Errorf("could not configure invoker: %w", err)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, r := range g.renderers {
		if err := r.Draw(screen); err != nil {
			g.logger.Error("could not draw", slog.String("error", err.Error()))
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Run() error {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Beherit")
	if err := ebiten.RunGame(g); err != nil {
		return err
	}
	return nil
}
