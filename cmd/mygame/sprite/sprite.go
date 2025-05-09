package sprite

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite is a sprite.
type Sprite struct {
	Img     *ebiten.Image
	OffsetX int
	OffsetY int
}

// centerXOffset calculates the horizontal offset to center the image along the X-axis.
func centerXOffset(img *ebiten.Image) int {
	return -(img.Bounds().Dx() / 2)
}

// bottomYOffset calculates the vertical offset to place the image at the bottom of the screen.
func bottomYOffset(img *ebiten.Image) int {
	return -(img.Bounds().Dy())
}

// bottomCenteredAlignedSprite creates a sprite that is centered along the X-axis and aligned to the bottom.
func bottomCenteredAlignedSprite(img *ebiten.Image) *Sprite {
	return &Sprite{
		Img:     img,
		OffsetX: centerXOffset(img),
		OffsetY: bottomYOffset(img),
	}
}

// Load returns all the sprites used in the game.
func LoadSkeletonDeath(i *ebiten.Image) map[uint]*Sprite {
	cells := CreateRectangleGrid(skeletonDeathCols, skeletonDeathRows, skeletonDeathWidth, skeletonDeathHeight)

	return map[uint]*Sprite{
		// Skeleton Death Sprites
		SkeletonDeath1: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[0][0]))),
		SkeletonDeath2: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[1][0]))),
		SkeletonDeath3: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[2][0]))),
		SkeletonDeath4: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[3][0]))),
		SkeletonDeath5: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[4][0]))),
		SkeletonDeath6: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[5][0]))),
	}
}

// Load returns all the sprites used in the game.
func LoadSkeletonMove(i *ebiten.Image) map[uint]*Sprite {
	cells := CreateRectangleGrid(skeletonMoveCols, skeletonMoveRows, skeletonMoveWidth, skeletonMoveHeight)

	return map[uint]*Sprite{
		// Skeleton Move Up Sprites
		SkeletonMoveUp1: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[0][0]))),
		SkeletonMoveUp2: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[1][0]))),
		SkeletonMoveUp3: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[2][0]))),
		SkeletonMoveUp4: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[3][0]))),
		SkeletonMoveUp5: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[4][0]))),
		SkeletonMoveUp6: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[5][0]))),
		SkeletonMoveUp7: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[6][0]))),
		SkeletonMoveUp8: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[7][0]))),
		SkeletonMoveUp9: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[8][0]))),

		// Skeleton Move Left Sprites
		SkeletonMoveLeft1: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[0][1]))),
		SkeletonMoveLeft2: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[1][1]))),
		SkeletonMoveLeft3: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[2][1]))),
		SkeletonMoveLeft4: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[3][1]))),
		SkeletonMoveLeft5: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[4][1]))),
		SkeletonMoveLeft6: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[5][1]))),
		SkeletonMoveLeft7: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[6][1]))),
		SkeletonMoveLeft8: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[7][1]))),
		SkeletonMoveLeft9: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[8][1]))),

		// Skeleton Move Down Sprites
		SkeletonMoveDown1: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[0][2]))),
		SkeletonMoveDown2: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[1][2]))),
		SkeletonMoveDown3: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[2][2]))),
		SkeletonMoveDown4: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[3][2]))),
		SkeletonMoveDown5: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[4][2]))),
		SkeletonMoveDown6: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[5][2]))),
		SkeletonMoveDown7: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[6][2]))),
		SkeletonMoveDown8: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[7][2]))),
		SkeletonMoveDown9: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[8][2]))),

		// Skeleton Move Right Sprites
		SkeletonMoveRight1: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[0][3]))),
		SkeletonMoveRight2: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[1][3]))),
		SkeletonMoveRight3: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[2][3]))),
		SkeletonMoveRight4: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[3][3]))),
		SkeletonMoveRight5: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[4][3]))),
		SkeletonMoveRight6: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[5][3]))),
		SkeletonMoveRight7: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[6][3]))),
		SkeletonMoveRight8: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[7][3]))),
		SkeletonMoveRight9: bottomCenteredAlignedSprite(ebiten.NewImageFromImage(i.SubImage(cells[8][3]))),
	}
}
