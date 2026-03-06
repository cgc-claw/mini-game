package assets

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func CreateBackground() error {
	os.MkdirAll("assets/backgrounds", 0755)
	
	img := image.NewRGBA(image.Rect(0, 0, 1600, 600))
	
	// Top gradient: dark to lighter purple
	for y := 0; y < 200; y++ {
		r := uint8(30 + (y * 20 / 200))
		g := uint8(20 + (y * 15 / 200))
		b := uint8(80 + (y * 30 / 200))
		col := color.RGBA{r, g, b, 255}
		for x := 0; x < 1600; x++ {
			img.Set(x, y, col)
		}
	}
	
	// Middle: pink/cyan synthwave gradient
	for y := 200; y < 400; y++ {
		r := uint8(120 + (y-200)*10/200)
		g := uint8(40)
		b := uint8(140 + (y-200)*10/200)
		col := color.RGBA{r, g, b, 255}
		for x := 0; x < 1600; x++ {
			img.Set(x, y, col)
		}
	}
	
	// Bottom: darker purple fade
	for y := 400; y < 600; y++ {
		r := uint8(60 + (600-y)*10/200)
		g := uint8(20)
		b := uint8(100 + (600-y)*5/200)
		col := color.RGBA{r, g, b, 255}
		for x := 0; x < 1600; x++ {
			img.Set(x, y, col)
		}
	}
	
	f, _ := os.Create("assets/backgrounds/background.png")
	defer f.Close()
	return png.Encode(f, img)
}
