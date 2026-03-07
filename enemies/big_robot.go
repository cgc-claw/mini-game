package enemies

import (
	"game/physics"
	"game/player"
	"game/sprites"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type BigRobot struct {
	X, Y       float64
	Width      float64
	Height     float64
	VX, VY     float64
	HP         int
	AnimTimer  int
	ShootTimer int
	Sprite     *ebiten.Image
}

func NewBigRobot(x, y float64) *BigRobot {
	w, h := float64(sprites.SpriteSize*2), float64(sprites.SpriteSize*2)
	return &BigRobot{
		X:          x,
		Y:          y,
		Width:      w,
		Height:     h,
		VX:         0,
		VY:         0,
		HP:         3,
		AnimTimer:  0,
		ShootTimer: 0,
		Sprite:     sprites.BigRobotSprite,
	}
}

func (e *BigRobot) AABB() *physics.AABB {
	return physics.NewAABB(e.X, e.Y, e.Width, e.Height)
}

func (e *BigRobot) Update(playerX float64, platforms []*physics.AABB) *player.Projectile {
	e.VY = physics.ApplyGravity(e.VY)
	e.Y += e.VY
	e.AnimTimer++

	onGround := false
	var currentPlat *physics.AABB
	for _, plat := range platforms {
		if e.AABB().Intersects(plat) {
			if e.VY > 0 {
				e.Y = plat.Y - e.Height
				e.VY = 0
				onGround = true
				currentPlat = plat
			}
		}
	}

	if onGround && currentPlat != nil {
		if playerX > e.X {
			e.VX = 0.5
		} else {
			e.VX = -0.5
		}

		if e.X+e.VX < currentPlat.X || e.X+e.Width+e.VX > currentPlat.X+currentPlat.Width {
			e.VX = 0
		}
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

	cx := e.Width / 2
	cy := e.Height / 2
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-cx, -cy)
	wobble := math.Sin(float64(e.AnimTimer)*0.2) * 0.05
	op.GeoM.Rotate(wobble)
	op.GeoM.Translate(cx, cy)

	if e.VX > 0 {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(e.Width, 0)
	}

	op.GeoM.Translate(sx, sy)
	screen.DrawImage(e.Sprite, op)
}
