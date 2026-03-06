package assets

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func CreateSprites() error {
	os.MkdirAll("assets/sprites", 0755)
	
	// Player sprite (32x32, red)
	if err := createPlayerSprite(); err != nil {
		return err
	}
	// Dino sprite (48x32, green)
	if err := createDinoSprite(); err != nil {
		return err
	}
	// BigRobot sprite (48x48, blue)
	if err := createBigRobotSprite(); err != nil {
		return err
	}
	// Boss sprite (64x64, purple)
	if err := createBossSprite(); err != nil {
		return err
	}
	// Platform texture (64x16, gray)
	if err := createPlatformTexture(); err != nil {
		return err
	}
	
	return nil
}

func createPlayerSprite() error {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	fillColor(img, 0, 0, 32, 32, color.RGBA{255, 0, 0, 255})
	fillColor(img, 10, 10, 22, 22, color.RGBA{50, 50, 50, 255})
	f, _ := os.Create("assets/sprites/player.png")
	defer f.Close()
	return png.Encode(f, img)
}

func createDinoSprite() error {
	img := image.NewRGBA(image.Rect(0, 0, 48, 32))
	fillColor(img, 0, 0, 48, 32, color.RGBA{0, 200, 0, 255})
	fillColor(img, 40, 4, 48, 12, color.RGBA{255, 100, 0, 255})
	f, _ := os.Create("assets/sprites/dino.png")
	defer f.Close()
	return png.Encode(f, img)
}

func createBigRobotSprite() error {
	img := image.NewRGBA(image.Rect(0, 0, 48, 48))
	fillColor(img, 4, 4, 44, 44, color.RGBA{0, 100, 255, 255})
	fillColor(img, 20, 0, 28, 8, color.RGBA{255, 50, 50, 255})
	f, _ := os.Create("assets/sprites/big_robot.png")
	defer f.Close()
	return png.Encode(f, img)
}

func createBossSprite() error {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	fillColor(img, 8, 8, 56, 56, color.RGBA{200, 0, 200, 255})
	fillColor(img, 24, 24, 40, 40, color.RGBA{255, 0, 255, 255})
	f, _ := os.Create("assets/sprites/boss.png")
	defer f.Close()
	return png.Encode(f, img)
}

func createPlatformTexture() error {
	img := image.NewRGBA(image.Rect(0, 0, 64, 16))
	fillColor(img, 0, 0, 64, 16, color.RGBA{100, 100, 120, 255})
	f, _ := os.Create("assets/sprites/platform.png")
	defer f.Close()
	return png.Encode(f, img)
}

func fillColor(img *image.RGBA, x0, y0, x1, y1 int, col color.Color) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			img.Set(x, y, col)
		}
	}
}
