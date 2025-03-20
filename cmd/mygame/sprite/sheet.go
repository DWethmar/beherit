package sprite

import (
	"embed"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	skeletonDeathHeight = 64
	skeletonDeathWidth  = 64
	skeletonMoveHeight  = 64
	skeletonMoveWidth   = 64
)

const (
	skeletonDeathCols = 6
	skeletonDeathRows = 1
)

const (
	skeletonMoveCols = 9
	skeletonMoveRows = 4
)

//go:embed images
var imagesFS embed.FS

type Image uint

const (
	SkeletonDeath1 Image = iota
	SkeletonDeath2
	SkeletonDeath3
	SkeletonDeath4
	SkeletonDeath5
	SkeletonDeath6
	SkeletonMoveUp1
	SkeletonMoveUp2
	SkeletonMoveUp3
	SkeletonMoveUp4
	SkeletonMoveUp5
	SkeletonMoveUp6
	SkeletonMoveUp7
	SkeletonMoveUp8
	SkeletonMoveUp9
	SkeletonMoveLeft1
	SkeletonMoveLeft2
	SkeletonMoveLeft3
	SkeletonMoveLeft4
	SkeletonMoveLeft5
	SkeletonMoveLeft6
	SkeletonMoveLeft7
	SkeletonMoveLeft8
	SkeletonMoveLeft9
	SkeletonMoveDown1
	SkeletonMoveDown2
	SkeletonMoveDown3
	SkeletonMoveDown4
	SkeletonMoveDown5
	SkeletonMoveDown6
	SkeletonMoveDown7
	SkeletonMoveDown8
	SkeletonMoveDown9
	SkeletonMoveRight1
	SkeletonMoveRight2
	SkeletonMoveRight3
	SkeletonMoveRight4
	SkeletonMoveRight5
	SkeletonMoveRight6
	SkeletonMoveRight7
	SkeletonMoveRight8
	SkeletonMoveRight9
)

// Spritesheet is a collection of sprites.
type Spritesheet struct {
	Skeleton map[Image]*ebiten.Image
}

func (s *Spritesheet) GetImage(img Image) *ebiten.Image {
	return s.Skeleton[img]
}

func loadSkeletonDeath(s *Spritesheet) error {
	img, err := LoadPng(imagesFS, "images/skeleton_death.png")
	if err != nil {
		return fmt.Errorf("failed to load images/skeleton_death.png: %w", err)
	}
	cells := CreateRectangleGrid(skeletonDeathCols, skeletonDeathRows, skeletonDeathWidth, skeletonDeathHeight)
	skeletonKillImg := ebiten.NewImageFromImage(img)
	s.Skeleton[SkeletonDeath1] = ebiten.NewImageFromImage(skeletonKillImg.SubImage(cells[0][0]))
	s.Skeleton[SkeletonDeath2] = ebiten.NewImageFromImage(skeletonKillImg.SubImage(cells[1][0]))
	s.Skeleton[SkeletonDeath3] = ebiten.NewImageFromImage(skeletonKillImg.SubImage(cells[2][0]))
	s.Skeleton[SkeletonDeath4] = ebiten.NewImageFromImage(skeletonKillImg.SubImage(cells[3][0]))
	s.Skeleton[SkeletonDeath5] = ebiten.NewImageFromImage(skeletonKillImg.SubImage(cells[4][0]))
	s.Skeleton[SkeletonDeath6] = ebiten.NewImageFromImage(skeletonKillImg.SubImage(cells[5][0]))
	return nil
}

func loadSkeletonMove(s *Spritesheet) error {
	img, err := LoadPng(imagesFS, "images/skeleton_move.png")
	if err != nil {
		return fmt.Errorf("failed to load images/skeleton_move.png: %w", err)
	}
	cells := CreateRectangleGrid(skeletonMoveCols, skeletonMoveRows, skeletonMoveWidth, skeletonMoveHeight)
	skeletonMoveImg := ebiten.NewImageFromImage(img)
	// Skeleton moving up
	s.Skeleton[SkeletonMoveUp1] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[0][0]))
	s.Skeleton[SkeletonMoveUp2] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[1][0]))
	s.Skeleton[SkeletonMoveUp3] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[2][0]))
	s.Skeleton[SkeletonMoveUp4] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[3][0]))
	s.Skeleton[SkeletonMoveUp5] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[4][0]))
	s.Skeleton[SkeletonMoveUp6] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[5][0]))
	s.Skeleton[SkeletonMoveUp7] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[6][0]))
	s.Skeleton[SkeletonMoveUp8] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[7][0]))
	s.Skeleton[SkeletonMoveUp9] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[8][0]))
	// Skeleton moving left
	s.Skeleton[SkeletonMoveLeft1] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[0][1]))
	s.Skeleton[SkeletonMoveLeft2] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[1][1]))
	s.Skeleton[SkeletonMoveLeft3] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[2][1]))
	s.Skeleton[SkeletonMoveLeft4] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[3][1]))
	s.Skeleton[SkeletonMoveLeft5] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[4][1]))
	s.Skeleton[SkeletonMoveLeft6] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[5][1]))
	s.Skeleton[SkeletonMoveLeft7] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[6][1]))
	s.Skeleton[SkeletonMoveLeft8] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[7][1]))
	s.Skeleton[SkeletonMoveLeft9] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[8][1]))
	// Skeleton moving down
	s.Skeleton[SkeletonMoveDown1] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[0][2]))
	s.Skeleton[SkeletonMoveDown2] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[1][2]))
	s.Skeleton[SkeletonMoveDown3] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[2][2]))
	s.Skeleton[SkeletonMoveDown4] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[3][2]))
	s.Skeleton[SkeletonMoveDown5] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[4][2]))
	s.Skeleton[SkeletonMoveDown6] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[5][2]))
	s.Skeleton[SkeletonMoveDown7] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[6][2]))
	s.Skeleton[SkeletonMoveDown8] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[7][2]))
	s.Skeleton[SkeletonMoveDown9] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[8][2]))
	// Skeleton moving right
	s.Skeleton[SkeletonMoveRight1] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[0][3]))
	s.Skeleton[SkeletonMoveRight2] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[1][3]))
	s.Skeleton[SkeletonMoveRight3] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[2][3]))
	s.Skeleton[SkeletonMoveRight4] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[3][3]))
	s.Skeleton[SkeletonMoveRight5] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[4][3]))
	s.Skeleton[SkeletonMoveRight6] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[5][3]))
	s.Skeleton[SkeletonMoveRight7] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[6][3]))
	s.Skeleton[SkeletonMoveRight8] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[7][3]))
	s.Skeleton[SkeletonMoveRight9] = ebiten.NewImageFromImage(skeletonMoveImg.SubImage(cells[8][3]))
	return nil
}

func NewSpritesheet() (*Spritesheet, error) {
	s := &Spritesheet{
		Skeleton: make(map[Image]*ebiten.Image),
	}
	if err := loadSkeletonDeath(s); err != nil {
		return nil, fmt.Errorf("failed to skeleton death sprites: %w", err)
	}
	if err := loadSkeletonMove(s); err != nil {
		return nil, fmt.Errorf("failed to load skeleton move sprites: %w", err)
	}
	return s, nil
}
