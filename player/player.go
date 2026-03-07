package player

import (
	"game/assets"
	"game/physics"
	"game/sprites"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	X, Y        float64
	Width       float64
	Height      float64
	VX, VY      float64
	FacingRight bool
	Grounded    bool
	JumpCount   int
	Invincible  int
	HP          int
	AnimTimer   int
	Sprite      *ebiten.Image
}

func New(x, y float64) *Player {
	w := float64(sprites.SpriteSize)
	h := float64(sprites.SpriteSize)
	return &Player{
		X:           x,
		Y:           y,
		Width:       w,
		Height:      h,
		VX:          0,
		VY:          0,
		FacingRight: true,
		Grounded:    false,
		JumpCount:   0,
		Invincible:  0,
		HP:          3,
		AnimTimer:   0,
		Sprite:      sprites.PlayerSprite,
	}
}

func (p *Player) AABB() *physics.AABB {
	return physics.NewAABB(p.X, p.Y, p.Width, p.Height)
}

func (p *Player) Update(keys map[ebiten.Key]bool, platforms []*Platform) {
	p.AnimTimer++

	moveLeft := keys[ebiten.KeyLeft]
	moveRight := keys[ebiten.KeyRight]
	jump := keys[ebiten.KeyUp]

	if moveLeft {
		p.VX = -physics.MoveSpeed
		p.FacingRight = false
	} else if moveRight {
		p.VX = physics.MoveSpeed
		p.FacingRight = true
	} else {
		p.VX = 0
	}

	if jump && (p.Grounded || p.JumpCount < 2) {
		assets.PlaySound("jump")
		p.VY = physics.JumpForce
		p.Grounded = false
		p.JumpCount++
	}

	p.VY = physics.ApplyGravity(p.VY)

	p.X += p.VX
	p.handleHorizontalCollisions(platforms)

	p.Y += p.VY
	p.Grounded = false
	p.handleVerticalCollisions(platforms)

	if p.X < 0 {
		p.X = 0
	}

	if p.Invincible > 0 {
		p.Invincible--
	}
}

func (p *Player) handleHorizontalCollisions(platforms []*Platform) {
	for _, plat := range platforms {
		if p.AABB().Intersects(plat.AABB()) {
			if p.VX > 0 {
				p.X = plat.X - p.Width
			} else if p.VX < 0 {
				p.X = plat.X + plat.Width
			}
		}
	}
}

func (p *Player) handleVerticalCollisions(platforms []*Platform) {
	for _, plat := range platforms {
		if p.AABB().Intersects(plat.AABB()) {
			if p.VY > 0 {
				p.Y = plat.Y - p.Height
				p.Grounded = true
				p.JumpCount = 0
				p.VY = 0
			} else if p.VY < 0 {
				p.Y = plat.Y + plat.Height
				p.VY = 0
			}
		}
	}
}

func (p *Player) Shoot() *Projectile {
	assets.PlaySound("shoot")
	dir := 1.0
	if !p.FacingRight {
		dir = -1
	}
	px := p.X + p.Width/2
	if !p.FacingRight {
		px = p.X - 8
	}
	return NewProjectile(px, p.Y+p.Height/2-4, dir, true)
}

func (p *Player) Draw(screen *ebiten.Image, camX, camY float64) {
	sx := p.X - camX
	sy := p.Y - camY

	if p.Invincible > 0 && p.Invincible%4 < 2 {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// Animation Transforms
	cx := p.Width / 2
	cy := p.Height / 2
	op.GeoM.Translate(-cx, -cy)

	if p.HP <= 0 {
		// Dead: fall over
		op.GeoM.Rotate(math.Pi / 2)
	} else if !p.Grounded {
		// Jump: Stretch vertically, squish horizontally
		op.GeoM.Scale(0.8, 1.2)
	} else if p.VX != 0 {
		// Run: waddle and bounce
		wobble := math.Sin(float64(p.AnimTimer)*0.4) * 0.15
		op.GeoM.Rotate(wobble)
		bounce := math.Abs(math.Sin(float64(p.AnimTimer)*0.4)) * 3.0
		cy -= bounce
	}

	op.GeoM.Translate(cx, cy)

	if !p.FacingRight {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(p.Width, 0)
	}

	op.GeoM.Translate(sx, sy)
	screen.DrawImage(p.Sprite, op)
}

type Platform struct {
	X, Y   float64
	Width  float64
	Height float64
	Sprite *ebiten.Image
}

func NewPlatform(x, y, w float64) *Platform {
	h := float64(16)
	return &Platform{
		X:      x,
		Y:      y,
		Width:  w,
		Height: h,
		Sprite: sprites.PlatformSprite,
	}
}

func (p *Platform) AABB() *physics.AABB {
	return physics.NewAABB(p.X, p.Y, p.Width, p.Height)
}

func (p *Platform) Draw(screen *ebiten.Image, camX, camY float64) {
	sx := p.X - camX
	sy := p.Y - camY

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(p.Width/64, 1)
	op.GeoM.Translate(sx, sy)

	screen.DrawImage(p.Sprite, op)
}

type Projectile struct {
	X, Y     float64
	Width    float64
	Height   float64
	VX       float64
	IsPlayer bool
	Sprite   *ebiten.Image
}

func NewProjectile(x, y, dir float64, isPlayer bool) *Projectile {
	w, h := 8.0, 8.0
	sprite := sprites.BulletSprite
	if !isPlayer {
		sprite = sprites.BulletSprite
	}
	return &Projectile{
		X:        x,
		Y:        y,
		Width:    w,
		Height:   h,
		VX:       dir * 8,
		IsPlayer: isPlayer,
		Sprite:   sprite,
	}
}

func (p *Projectile) AABB() *physics.AABB {
	return physics.NewAABB(p.X, p.Y, p.Width, p.Height)
}

func (p *Projectile) Update() {
	p.X += p.VX
}

func (p *Projectile) Draw(screen *ebiten.Image, camX, camY float64) {
	sx := p.X - camX
	sy := p.Y - camY
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(sx, sy)
	screen.DrawImage(p.Sprite, op)
}
