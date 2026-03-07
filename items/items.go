package items

import (
	"game/physics"
	"game/sprites"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type ItemType int

const (
	TypeHealth ItemType = iota
	TypeAmmo
	TypeLife
)

type Item struct {
	X, Y      float64
	Width     float64
	Height    float64
	VY        float64
	Type      ItemType
	AnimTimer int
	Sprite    *ebiten.Image
}

func NewItem(x, y float64, itemType ItemType) *Item {
	var sprite *ebiten.Image
	switch itemType {
	case TypeHealth:
		sprite = sprites.ItemHealthSprite
	case TypeAmmo:
		sprite = sprites.ItemAmmoSprite
	case TypeLife:
		sprite = sprites.ItemLifeSprite
	}

	return &Item{
		X:         x,
		Y:         y,
		Width:     16,
		Height:    16,
		VY:        0,
		Type:      itemType,
		AnimTimer: 0,
		Sprite:    sprite,
	}
}

func (i *Item) AABB() *physics.AABB {
	return physics.NewAABB(i.X, i.Y, i.Width, i.Height)
}

func (i *Item) Update(platforms []*physics.AABB) {
	i.VY = physics.ApplyGravity(i.VY)
	i.Y += i.VY
	i.AnimTimer++

	// Collision with platforms
	for _, p := range platforms {
		if i.AABB().Intersects(p) {
			if i.VY > 0 { // Falling
				i.Y = p.Y - i.Height
				i.VY = 0
			}
		}
	}
}

func (i *Item) Draw(screen *ebiten.Image, camX, camY float64) {
	sx := i.X - camX
	sy := i.Y - camY

	cx := i.Width / 2
	cy := i.Height / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-cx, -cy)

	// Floating bob
	bobbing := math.Sin(float64(i.AnimTimer)*0.1) * 2.0

	op.GeoM.Translate(cx, cy+bobbing)
	op.GeoM.Translate(sx, sy)

	screen.DrawImage(i.Sprite, op)
}
