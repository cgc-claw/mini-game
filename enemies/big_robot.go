package enemies

import (
	"game/player"
	"game/physics"
	"game/sprites"

	"github.com/hajimehoshi/ebiten/v2"
)

type BigRobot struct {
	X, Y        float64
	Width       float64
	Height      float64
	VX          float64
	HP          int
	ShootTimer  int
	Sprite      *ebiten.Image
}

func NewBigRobot(x, y float64) *BigRobot {
	w, h := float64(sprites.SpriteSize*2), float64(sprites.SpriteSize*2)
	return &BigRobot{
		X:         x,
		Y:         y,
		Width:     w,
		Height:    h,
		VX:        0,
		HP:        3,
		ShootTimer: 0,
		Sprite:    sprites.BigRobotSprite,
	}
}

func (e *BigRobot) AABB() *physics.AABB {
	return physics.NewAABB(e.X, e.Y, e.Width, e.Height)
}

func (e *BigRobot) Update(playerX float64) *player.Projectile {
	if playerX > e.X {
		e.VX = 1
	} else {
		e.VX = -1
	}
	e.X += e.VX

	e.ShootTimer++
	if e.ShootTimer >= 90 {
		e.ShootTimer = 0
		dir := 1.0
		if playerX <= e.X {
			dir = -1
		}
		return player.NewProjectile(e.X+e.Width/2, e.Y+e.Height/2-4, dir, false)
	}
	return nil
}

func (e *BigRobot) Draw(screen *ebiten.Image, camX, camY float64) {
	sx := e.X - camX
	sy := e.Y - camY
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(sx, sy)
	screen.DrawImage(e.Sprite, op)
}
