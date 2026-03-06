package enemies

import (
	"game/player"
	"game/physics"
	"game/sprites"

	"github.com/hajimehoshi/ebiten/v2"
)

type AlienBoss struct {
	X, Y        float64
	Width       float64
	Height      float64
	VX, VY      float64
	HP          int
	ShootTimer  int
	Phase       int
	TargetY     float64
	Sprite      *ebiten.Image
}

func NewAlienBoss(x, y float64) *AlienBoss {
	w, h := float64(sprites.SpriteSize*3), float64(sprites.SpriteSize*3)
	return &AlienBoss{
		X:         x,
		Y:         y,
		Width:     w,
		Height:    h,
		VX:        0,
		VY:        0,
		HP:        10,
		ShootTimer: 0,
		Phase:     0,
		TargetY:   y,
		Sprite:    sprites.AlienBossSprite,
	}
}

func (e *AlienBoss) AABB() *physics.AABB {
	return physics.NewAABB(e.X, e.Y, e.Width, e.Height)
}

func (e *AlienBoss) Update(playerX, playerY float64) *player.Projectile {
	e.Phase = (e.Phase + 1) % 180

	if e.Phase < 60 {
		if playerX > e.X {
			e.VX = 1.5
		} else {
			e.VX = -1.5
		}
		e.X += e.VX
	} else if e.Phase < 90 {
		e.TargetY = playerY
	} else {
		e.Y += (e.TargetY - e.Y) * 0.05
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
	
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(sx, sy)
	screen.DrawImage(e.Sprite, op)
	
	tentacleColor := &ebiten.ColorM{}
	tentacleColor.ChangeHSV(0.8, 1.0, 1.0)
	op.ColorM = *tentacleColor
	
	for i := 0; i < 3; i++ {
		offset := float64(i-1) * 10
		tx := sx + e.Width/2 + offset
		ty := sy + e.Height
		
		for j := 0; j < 8; j++ {
			tw := 4.0
			th := 6.0
			pop := &ebiten.DrawImageOptions{}
			pop.GeoM.Translate(tx-tw/2, ty)
			pop.ColorM = *tentacleColor
			screen.DrawImage(e.Sprite, pop)
			ty += th
		}
	}
}
