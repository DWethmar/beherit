package beherit

import (
	"github.com/hajimehoshi/ebiten/v2"
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

type System interface {
	Init(*EntityManager, *ComponentManager) error
	Update() error
	Draw(screen *ebiten.Image)
}

type Game struct {
	state            GameState
	entityManager    *EntityManager
	componentManager *ComponentManager
	systems          []System
}

func NewGame(
	systems []System,
) *Game {
	return &Game{
		state:            GameStateInitializing,
		entityManager:    NewEntityManager(),
		componentManager: NewComponentManager(),
		systems:          systems,
	}
}

func (g *Game) Update() error {
	switch g.state {
	case GameStateInitializing:
		for _, system := range g.systems {
			if err := system.Init(g.entityManager, g.componentManager); err != nil {
				return err
			}
		}
		g.state = GameStateRunning
	case GameStateRunning:
		for _, system := range g.systems {
			if err := system.Update(); err != nil {
				return err
			}
		}
	case GameStatePaused:
	case GameStateStopped:
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, system := range g.systems {
		system.Draw(screen)
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
