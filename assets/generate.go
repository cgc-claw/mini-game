package assets

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

// GenerateAssets creates placeholder sprite and background images
func GenerateAssets() error {
	// Create sprites directory
	os.MkdirAll("assets/sprites", 0755)
	os.MkdirAll("assets/sounds", 0755)
	os.MkdirAll("assets/backgrounds", 0755)

	// Generate player sprite (32x32 - red robot)
	if err := generatePlayerSprite(); err != nil {
		return err
	}

	// Generate enemy sprites
	if err := generateDinoSprite(); err != nil {
		return err
	}
	if err := generateBigRobotSprite(); err != nil {
		return err
	}
	if err := generateBossSprite(); err != nil {
		return err
	}

	// Generate background
	if err := generateBackground(); err != nil {
		return err
	}

	// Generate platform texture
	if err := generatePlatformTexture(); err != nil {
		return err
	}

	return nil
}

func generatePlayerSprite() error {
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	// Red robot body
	fillRect(img, 8, 8, 24, 24, color.RGBA{255, 0, 0, 255})
	// Blue eyes
	fillRect(img, 12, 12, 14, 14, color.RGBA{0, 0, 255, 255})
	fillRect(img, 18, 12, 20, 14, color.RGBA{0, 0, 255, 255})

	f, _ := os.Create("assets/sprites/player.png")
	defer f.Close()
	return png.Encode(f, img)
}

func generateDinoSprite() error {
	img := image.NewRGBA(image.Rect(0, 0, 48, 32))
	// Green dinosaur body
	fillRect(img, 10, 10, 38, 25, color.RGBA{0, 200, 0, 255})
	// Head spike
	fillTriangle(img, 35, 8, 40, 0, 45, 8, color.RGBA{100, 200, 0, 255})
	// Yellow eyes
	fillRect(img, 20, 12, 22, 14, color.RGBA{255, 255, 0, 255})
	fillRect(img, 28, 12, 30, 14, color.RGBA{255, 255, 0, 255})

	f, _ := os.Create("assets/sprites/dino.png")
	defer f.Close()
	return png.Encode(f, img)
}

func generateBigRobotSprite() error {
	img := image.NewRGBA(image.Rect(0, 0, 48, 48))
	// Large blue robot
	fillRect(img, 8, 8, 40, 40, color.RGBA{0, 100, 255, 255})
	// Cyan armor plating
	fillRect(img, 12, 12, 36, 36, color.RGBA{0, 200, 255, 255})
	// Red cannon on top
	fillRect(img, 20, 4, 28, 10, color.RGBA{255, 50, 50, 255})
	// Yellow visor
	fillRect(img, 16, 18, 32, 22, color.RGBA{255, 255, 0, 255})

	f, _ := os.Create("assets/sprites/big_robot.png")
	defer f.Close()
	return png.Encode(f, img)
}

func generateBossSprite() error {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	// Purple alien boss
	fillRect(img, 12, 12, 52, 52, color.RGBA{200, 0, 200, 255})
	// Magenta energy core
	fillRect(img, 24, 24, 40, 40, color.RGBA{255, 0, 255, 255})
	// Cyan tentacles (left)
	fillRect(img, 4, 20, 12, 44, color.RGBA{0, 255, 255, 255})
	// Cyan tentacles (right)
	fillRect(img, 52, 20, 60, 44, color.RGBA{0, 255, 255, 255})
	// Glowing eyes
	fillRect(img, 20, 20, 24, 24, color.RGBA{255, 255, 0, 255})
	fillRect(img, 40, 20, 44, 24, color.RGBA{255, 255, 0, 255})

	f, _ := os.Create("assets/sprites/boss.png")
	defer f.Close()
	return png.Encode(f, img)
}

func generateBackground() error {
	img := image.NewRGBA(image.Rect(0, 0, 1600, 600))
	// Synthwave-ish gradient background
	// Top: dark purple/blue
	for y := 0; y < 200; y++ {
		r := uint8(20 + (y * 10 / 200))
		g := uint8(20 + (y * 20 / 200))
		b := uint8(60 + (y * 30 / 200))
		for x := 0; x < 1600; x++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	// Middle: pink/cyan grid overlay
	for y := 200; y < 400; y++ {
		r := uint8(80 + (y-200)*20/200)
		g := uint8(30)
		b := uint8(100 + (y-200)*10/200)
		for x := 0; x < 1600; x++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	// Bottom: darker colors
	for y := 400; y < 600; y++ {
		r := uint8(40 + (600-y)*5/200)
		g := uint8(10)
		b := uint8(60 + (600-y)*3/200)
		for x := 0; x < 1600; x++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	f, _ := os.Create("assets/backgrounds/bg.png")
	defer f.Close()
	return png.Encode(f, img)
}

func generatePlatformTexture() error {
	img := image.NewRGBA(image.Rect(0, 0, 64, 32))
	// Metal texture
	fillRect(img, 0, 0, 64, 32, color.RGBA{100, 100, 120, 255})
	// Highlights
	for x := 0; x < 64; x += 8 {
		fillRect(img, x, 0, x+4, 4, color.RGBA{150, 150, 170, 255})
	}
	// Shadows
	for x := 0; x < 64; x += 8 {
		fillRect(img, x+4, 28, x+8, 32, color.RGBA{60, 60, 80, 255})
	}

	f, _ := os.Create("assets/sprites/platform.png")
	defer f.Close()
	return png.Encode(f, img)
}

// Helper functions
func fillRect(img *image.RGBA, x0, y0, x1, y1 int, col color.Color) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			img.Set(x, y, col)
		}
	}
}

func fillTriangle(img *image.RGBA, x1, y1, x2, y2, x3, y3 int, col color.Color) {
	// Simple triangle fill using barycentric coordinates would be complex
	// For now, just fill a rectangle approximation
	minX, maxX := x1, x1
	minY, maxY := y1, y1
	if x2 < minX {
		minX = x2
	}
	if x2 > maxX {
		maxX = x2
	}
	if x3 < minX {
		minX = x3
	}
	if x3 > maxX {
		maxX = x3
	}
	if y2 < minY {
		minY = y2
	}
	if y2 > maxY {
		maxY = y2
	}
	if y3 < minY {
		minY = y3
	}
	if y3 > maxY {
		maxY = y3
	}
	fillRect(img, minX, minY, maxX, maxY, col)
}
