package beherit

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"gopkg.in/yaml.v3"
)

const (
	windowWidth  = 640
	windowHeight = 480
)

const (
	GameCreatedTrigger string = "game.created"
	UpdateTrigger      string = "game.updated"
	GameInputCursor    string = "game.input.cursor"
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

type Cursor struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	globalEnv := map[string]any{
		"cursor": Cursor{X: x, Y: y},
	}

	switch g.state {
	case GameStateInitializing:
		if err := g.loadConfig(); err != nil {
			return err
		}
		if err := g.Invoker.Invoke(GameCreatedTrigger, globalEnv); err != nil {
			return fmt.Errorf("could not invoke %s trigger: %w", GameCreatedTrigger, err)
		}
		g.state = GameStateRunning
	case GameStateRunning:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if err := g.Invoker.Invoke(GameInputCursor, globalEnv); err != nil {
				return fmt.Errorf("could not invoke %s trigger: %w", GameInputCursor, err)
			}
		}
		if err := g.Invoker.Invoke(UpdateTrigger, map[string]any{}); err != nil {
			return fmt.Errorf("could not invoke %s trigger: %w", UpdateTrigger, err)
		}
	case GameStatePaused:
	case GameStateStopped:
	}
	return nil
}

func (g *Game) WatchConfig(ctx context.Context) error {
	lastModTime := time.Now()
	for {
		select {
		case <-time.After(1 * time.Second):
			info, err := os.Stat(g.configFile)
			if err != nil {
				g.logger.Error("could not stat file", slog.String("file", g.configFile), slog.String("error", err.Error()))
				continue
			}
			if info.ModTime().After(lastModTime) {
				lastModTime = info.ModTime()
				if err := g.loadConfig(); err != nil {
					return fmt.Errorf("could not load config: %w", err)
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (g *Game) loadConfig() error {
	file, err := os.Open(g.configFile)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()
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
