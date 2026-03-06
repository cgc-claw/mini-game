package enemies

import (
	"game/physics"
	"game/sprites"

	"github.com/hajimehoshi/ebiten/v2"
)

type DinoRobot struct {
	X, Y       float64
	Width      float64
	Height     float64
	VX, VY     float64
	HP         int
	Sprite     *ebiten.Image
}

func NewDinoRobot(x, y float64) *DinoRobot {
	w, h := float64(sprites.SpriteSize), float64(sprites.SpriteSize)
	return &DinoRobot{
		X:      x,
		Y:      y,
		Width:  w,
		Height: h,
		VX:     -2,
		VY:     0,
		HP:     1,
		Sprite: sprites.DinoRobotSprite,
	}
}

func (e *DinoRobot) AABB() *physics.AABB {
	return physics.NewAABB(e.X, e.Y, e.Width, e.Height)
}

func (e *DinoRobot) Update(playerX float64, platforms []*physics.AABB) {
	e.VY = physics.ApplyGravity(e.VY)
	e.Y += e.VY

	for _, plat := range platforms {
		if e.AABB().Intersects(plat) {
			if e.VY > 0 {
				e.Y = plat.Y - e.Height
				e.VY = -8
				if playerX > e.X {
					e.VX = 2
				} else {
					e.VX = -2
				}
			}
		}
	}
	e.X += e.VX
}

func (e *DinoRobot) Draw(screen *ebiten.Image, camX, camY float64) {
	sx := e.X - camX
	sy := e.Y - camY
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(sx, sy)
	screen.DrawImage(e.Sprite, op)
}
