package sprites

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const SpriteSize = 32

func createColoredSprite(w, h int, mainColor, accentColor color.RGBA) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, mainColor)
		}
	}
	
	edgeSize := 2
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x < edgeSize || x >= w-edgeSize || y < edgeSize || y >= h-edgeSize {
				img.Set(x, y, accentColor)
			}
		}
	}
	
	eyeY := h / 4
	eyeSize := 3
	for y := eyeY; y < eyeY+eyeSize; y++ {
		for x := w/3; x < w/3+eyeSize; x++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
		for x := 2*w/3-eyeSize; x < 2*w/3; x++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}
	
	return ebiten.NewImageFromImage(img)
}

var (
	PlayerSprite    *ebiten.Image
	DinoRobotSprite *ebiten.Image
	BigRobotSprite  *ebiten.Image
	AlienBossSprite *ebiten.Image
	PlatformSprite  *ebiten.Image
	BulletSprite    *ebiten.Image
)

func init() {
	PlayerSprite = createColoredSprite(SpriteSize, SpriteSize,
		color.RGBA{192, 192, 192, 255}, // Silver
		color.RGBA{220, 20, 60, 255},   // Red accent
	)

	DinoRobotSprite = createColoredSprite(SpriteSize, SpriteSize,
		color.RGBA{34, 139, 34, 255},   // Green
		color.RGBA{0, 100, 0, 255},     // Dark green
	)

	BigRobotSprite = createColoredSprite(SpriteSize*2, SpriteSize*2,
		color.RGBA{255, 140, 0, 255},  // Orange
		color.RGBA{139, 69, 19, 255},   // Brown
	)

	AlienBossSprite = createColoredSprite(SpriteSize*3, SpriteSize*3,
		color.RGBA{128, 0, 128, 255},  // Purple
		color.RGBA{75, 0, 130, 255},   // Indigo
	)

	PlatformSprite = createColoredSprite(64, 16,
		color.RGBA{139, 119, 101, 255}, // Brown
		color.RGBA{80, 60, 40, 255},    // Dark brown
	)

	BulletSprite = createColoredSprite(8, 8,
		color.RGBA{255, 255, 0, 255},  // Yellow
		color.RGBA{255, 165, 0, 255},   // Orange
	)
}
