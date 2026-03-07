package enemies

import (
	"game/physics"
	"game/sprites"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type DinoRobot struct {
	X, Y      float64
	Width     float64
	Height    float64
	VX, VY    float64
	HP        int
	Direction float64
	AnimTimer int
	Sprite    *ebiten.Image
}

func NewDinoRobot(x, y float64) *DinoRobot {
	w, h := float64(sprites.SpriteSize), float64(sprites.SpriteSize)
	return &DinoRobot{
		X:         x,
		Y:         y,
		Width:     w,
		Height:    h,
		VX:        -1.5,
		VY:        0,
		HP:        1,
		Direction: -1,
		AnimTimer: 0,
		Sprite:    sprites.DinoRobotSprite,
	}
}

func (e *DinoRobot) AABB() *physics.AABB {
	return physics.NewAABB(e.X, e.Y, e.Width, e.Height)
}

func (e *DinoRobot) Update(playerX float64, platforms []*physics.AABB) {
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
		distance := math.Abs(playerX - e.X)
		shouldTrack := distance < 128 // Two blocks

		if shouldTrack {
			if playerX > e.X {
				e.Direction = 1
			} else {
				e.Direction = -1
			}
		}

		e.VX = e.Direction * 1.5

		// Reverse direction at edge if patrolling or hit bounds
		if e.X+e.VX < currentPlat.X || e.X+e.Width+e.VX > currentPlat.X+currentPlat.Width {
			if !shouldTrack {
				e.Direction *= -1
				e.VX = e.Direction * 1.5
			} else {
				e.VX = 0 // Stop at edge if tracking
			}
		}
	}

	e.X += e.VX
}

func (e *DinoRobot) Draw(screen *ebiten.Image, camX, camY float64) {
	sx := e.X - camX
	sy := e.Y - camY
	cx := e.Width / 2
	cy := e.Height / 2
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-cx, -cy)
	wobble := math.Sin(float64(e.AnimTimer)*0.3) * 0.1
	op.GeoM.Rotate(wobble)
	op.GeoM.Translate(cx, cy)

	if e.VX > 0 {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(e.Width, 0)
	}

	op.GeoM.Translate(sx, sy)
	screen.DrawImage(e.Sprite, op)
}
