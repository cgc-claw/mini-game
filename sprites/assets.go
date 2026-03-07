package sprites

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const SpriteSize = 32

func loadSprite(path string, targetW, targetH int) *ebiten.Image {
	_, originalImg, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Printf("Failed to load sprite %s: %v", path, err)
		// Return empty image as fallback
		return ebiten.NewImage(targetW, targetH)
	}

	bounds := originalImg.Bounds()
	origW := bounds.Dx()
	origH := bounds.Dy()

	rgbaImg := image.NewRGBA(bounds)

	// Copy and mask white background
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := originalImg.At(x, y)
			r, g, b, _ := c.RGBA()
			// if it's very close to pure white, make it transparent
			if r > 60000 && g > 60000 && b > 60000 {
				rgbaImg.Set(x, y, color.RGBA{0, 0, 0, 0})
			} else {
				rgbaImg.Set(x, y, c)
			}
		}
	}

	maskedEbitenImg := ebiten.NewImageFromImage(rgbaImg)

	scaleX := float64(targetW) / float64(origW)
	scaleY := float64(targetH) / float64(origH)

	scaledImg := ebiten.NewImage(targetW, targetH)
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM.Scale(scaleX, scaleY)
	scaledImg.DrawImage(maskedEbitenImg, op)

	return scaledImg
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
	PlayerSprite = loadSprite("assets/sprites/player.png", SpriteSize, SpriteSize)
	DinoRobotSprite = loadSprite("assets/sprites/dino.png", 48, 32)
	BigRobotSprite = loadSprite("assets/sprites/big_robot.png", SpriteSize*2, SpriteSize*2)
	AlienBossSprite = loadSprite("assets/sprites/boss.png", SpriteSize*3, SpriteSize*3)
	PlatformSprite = loadSprite("assets/sprites/platform.png", 64, 16)
	BulletSprite = loadSprite("assets/sprites/bullet.png", 8, 8)
}
