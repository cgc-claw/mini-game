package enemies

import (
	"game/physics"
	"game/player"
	"game/sprites"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type AlienBoss struct {
	X, Y       float64
	Width      float64
	Height     float64
	VX, VY     float64
	HP         int
	Phase      int
	TargetY    float64
	AnimTimer  int
	ShootTimer int
	Sprite     *ebiten.Image
}

func NewAlienBoss(x, y float64) *AlienBoss {
	w, h := float64(sprites.SpriteSize*3), float64(sprites.SpriteSize*3)
	return &AlienBoss{
		X:          x,
		Y:          y,
		Width:      w,
		Height:     h,
		VX:         0,
		VY:         0,
		HP:         10,
		ShootTimer: 0,
		Phase:      0,
		TargetY:    y,
		AnimTimer:  0,
		Sprite:     sprites.AlienBossSprite,
	}
}

func (e *AlienBoss) AABB() *physics.AABB {
	return physics.NewAABB(e.X, e.Y, e.Width, e.Height)
}

func (e *AlienBoss) Update(playerX, playerY float64) *player.Projectile {
	e.Phase = (e.Phase + 1) % 240
	e.AnimTimer++

	if e.Phase < 120 {
		// Hover and chase horizontally
		if playerX > e.X {
			e.VX = 1.2
		} else {
			e.VX = -1.2
		}
		e.X += e.VX

		// Sine wave horizontal floating
		e.Y += math.Sin(float64(e.AnimTimer)*0.05) * 1.5
	} else if e.Phase < 150 {
		// Lock Y target
		e.TargetY = playerY
	} else {
		// Swoop down/up
		e.Y += (e.TargetY - e.Y) * 0.08
	}

	e.ShootTimer++
	if e.ShootTimer >= 45 {
		e.ShootTimer = 0
		dir := 1.0
		if playerX <= e.X {
			dir = -1
		}
		return player.NewProjectile(e.X+e.Width/2, e.Y+e.Height/2-4, dir, false)
	}
	return nil
}

func (e *AlienBoss) Draw(screen *ebiten.Image, camX, camY float64) {
	sx := e.X - camX
	sy := e.Y - camY
	cx := e.Width / 2
	cy := e.Height / 2
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-cx, -cy)
	bobbing := math.Sin(float64(e.AnimTimer)*0.1) * 0.05
	op.GeoM.Rotate(bobbing)
	op.GeoM.Translate(cx, cy)

	if e.VX > 0 {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(e.Width, 0)
	}

	op.GeoM.Translate(sx, sy)
	screen.DrawImage(e.Sprite, op)

	tentacleColor := &ebiten.ColorM{}
	tentacleColor.ChangeHSV(0.8, 1.0, 1.0)

	for i := 0; i < 3; i++ {
		offset := float64(i-1) * 20
		tx := sx + e.Width/2 + offset
		ty := sy + e.Height - 10

		for j := 0; j < 6; j++ {
			tw := 8.0
			th := 10.0
			pop := &ebiten.DrawImageOptions{}

			// Sinusoidal sway for tentacles
			sway := math.Sin(float64(e.AnimTimer)*0.1+float64(i)+float64(j)*0.5) * 4.0

			pop.GeoM.Scale(8.0/float64(sprites.SpriteSize*3), 10.0/float64(sprites.SpriteSize*3)) // scale down sprite for tentacle pieces
			pop.GeoM.Translate(tx-tw/2+sway, ty)
			pop.ColorM = *tentacleColor
			screen.DrawImage(e.Sprite, pop)
			ty += th
		}
	}
}
