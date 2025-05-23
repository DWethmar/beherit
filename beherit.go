package beherit

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"sync"
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
	GameUpdateTrigger  string = "game.updated"
	GameInputCursor    string = "game.input.cursor"
	GameDraw           string = "game.draw"
)

type GameState uint

const (
	GameStateInitializing GameState = iota
	GameStateRunning
	GameStatePaused
	GameStateStopped
)

type Game struct {
	mux        sync.RWMutex
	logger     *slog.Logger
	configFile string
	Config     *Config
	state      GameState
	Invoker    *Invoker
}

type Options struct {
	Logger     *slog.Logger
	ConfigFile string
	Invoker    *Invoker
}

func NewGame(opts Options) *Game {
	return &Game{
		logger:     opts.Logger,
		configFile: opts.ConfigFile,
		state:      GameStateInitializing,
		Invoker:    opts.Invoker,
	}
}

type Cursor struct {
	X int `json:"x" expr:"x"`
	Y int `json:"y" expr:"y"`
}

type Env struct {
	Screen *ebiten.Image `json:"screen" expr:"screen"`
	Cursor Cursor        `json:"cursor" expr:"cursor"`
}

func (g *Game) Update() error {
	g.mux.RLock()
	defer g.mux.RUnlock()

	x, y := ebiten.CursorPosition()
	env := map[string]any{
		"cursor": Cursor{
			X: x,
			Y: y,
		},
	}
	maps.Copy(env, g.Config.Env)
	switch g.state {
	case GameStateInitializing:
		if err := g.Invoker.Invoke(GameCreatedTrigger, env); err != nil {
			return fmt.Errorf("could not invoke %s trigger: %w", GameCreatedTrigger, err)
		}
		g.state = GameStateRunning
	case GameStateRunning:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if err := g.Invoker.Invoke(GameInputCursor, env); err != nil {
				return fmt.Errorf("could not invoke %s trigger: %w", GameInputCursor, err)
			}
		}
		if err := g.Invoker.Invoke(GameUpdateTrigger, map[string]any{}); err != nil {
			return fmt.Errorf("could not invoke %s trigger: %w", GameUpdateTrigger, err)
		}
	case GameStatePaused:
	case GameStateStopped:
	}
	return nil
}

func (g *Game) WatchConfig(ctx context.Context) {
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
					g.logger.Error("could not load config", slog.String("file", g.configFile), slog.String("error", err.Error()))
					continue
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (g *Game) loadConfig() error {
	g.mux.Lock()
	defer g.mux.Unlock()
	file, err := os.Open(g.configFile)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()
	var config Config
	if err = yaml.NewDecoder(file).Decode(&config); err != nil {
		return fmt.Errorf("could not decode yaml: %w", err)
	}
	if err = g.Invoker.SetTriggers(config.Triggers); err != nil {
		return fmt.Errorf("failed to set triggers: %w", err)
	}
	g.Config = &config
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	env := map[string]any{
		"screen": screen,
		"cursor": Cursor{
			X: x,
			Y: y,
		},
	}
	maps.Copy(env, g.Config.Env)
	if err := g.Invoker.Invoke(GameDraw, env); err != nil {
		g.logger.Error("could not invoke draw trigger", slog.String("error", err.Error()))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Run() error {
	if err := g.loadConfig(); err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Beherit")
	if err := ebiten.RunGame(g); err != nil {
		return err
	}
	return nil
}
