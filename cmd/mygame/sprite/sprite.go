package sprite

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprites map[Image]*Sprite

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
func Load(s *Spritesheet) Sprites {
	return Sprites{
		// Skeleton Death Sprites
		SkeletonDeath1: bottomCenteredAlignedSprite(s.GetImage(SkeletonDeath1)),
		SkeletonDeath2: bottomCenteredAlignedSprite(s.GetImage(SkeletonDeath2)),
		SkeletonDeath3: bottomCenteredAlignedSprite(s.GetImage(SkeletonDeath3)),
		SkeletonDeath4: bottomCenteredAlignedSprite(s.GetImage(SkeletonDeath4)),
		SkeletonDeath5: bottomCenteredAlignedSprite(s.GetImage(SkeletonDeath5)),
		SkeletonDeath6: bottomCenteredAlignedSprite(s.GetImage(SkeletonDeath6)),
		// Skeleton Move Up Sprites
		SkeletonMoveUp1: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveUp1)),
		SkeletonMoveUp2: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveUp2)),
		SkeletonMoveUp3: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveUp3)),
		SkeletonMoveUp4: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveUp4)),
		SkeletonMoveUp5: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveUp5)),
		SkeletonMoveUp6: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveUp6)),
		SkeletonMoveUp7: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveUp7)),
		SkeletonMoveUp8: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveUp8)),
		SkeletonMoveUp9: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveUp9)),
		// Skeleton Move Left Sprites
		SkeletonMoveLeft1: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveLeft1)),
		SkeletonMoveLeft2: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveLeft2)),
		SkeletonMoveLeft3: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveLeft3)),
		SkeletonMoveLeft4: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveLeft4)),
		SkeletonMoveLeft5: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveLeft5)),
		SkeletonMoveLeft6: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveLeft6)),
		SkeletonMoveLeft7: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveLeft7)),
		SkeletonMoveLeft8: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveLeft8)),
		SkeletonMoveLeft9: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveLeft9)),
		// Skeleton Move Down Sprites
		SkeletonMoveDown1: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveDown1)),
		SkeletonMoveDown2: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveDown2)),
		SkeletonMoveDown3: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveDown3)),
		SkeletonMoveDown4: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveDown4)),
		SkeletonMoveDown5: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveDown5)),
		SkeletonMoveDown6: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveDown6)),
		SkeletonMoveDown7: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveDown7)),
		SkeletonMoveDown8: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveDown8)),
		SkeletonMoveDown9: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveDown9)),
		// Skeleton Move Right Sprites
		SkeletonMoveRight1: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveRight1)),
		SkeletonMoveRight2: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveRight2)),
		SkeletonMoveRight3: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveRight3)),
		SkeletonMoveRight4: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveRight4)),
		SkeletonMoveRight5: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveRight5)),
		SkeletonMoveRight6: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveRight6)),
		SkeletonMoveRight7: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveRight7)),
		SkeletonMoveRight8: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveRight8)),
		SkeletonMoveRight9: bottomCenteredAlignedSprite(s.GetImage(SkeletonMoveRight9)),
	}
}
