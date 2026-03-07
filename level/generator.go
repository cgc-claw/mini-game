package level

import (
	"game/enemies"
	"game/items"
	"game/physics"
	"game/player"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Generator struct {
	Platforms     []*player.Platform
	Enemies       []interface{}
	Items         []*items.Item
	LastPlatformX [3]float64 // 3 layers
	ScreenWidth   int
	ScreenHeight  int
	Difficulty    float64
}

func New(w, h int) *Generator {
	return &Generator{
		Platforms:     make([]*player.Platform, 0),
		Enemies:       make([]interface{}, 0),
		Items:         make([]*items.Item, 0),
		LastPlatformX: [3]float64{100, 400, 700},
		ScreenWidth:   w,
		ScreenHeight:  h,
		Difficulty:    1.0,
	}
}

func (g *Generator) Init(screenW, screenH int) {
	g.Platforms = append(g.Platforms, player.NewPlatform(0, float64(screenH-40), 400))
	g.LastPlatformX[0] = 400
	g.LastPlatformX[1] = 600
	g.LastPlatformX[2] = 800
}

func (g *Generator) Generate(cameraX float64) {
	generateUntil := cameraX + float64(g.ScreenWidth) + 500

	for layer := 0; layer < 3; layer++ {
		for g.LastPlatformX[layer] < generateUntil {
			gap := float64(rand.Intn(80) + 60) // smaller gaps since there are many layers
			if layer == 0 && rand.Float64() < 0.2 {
				gap = float64(rand.Intn(100) + 150) // occasional large gap on bottom layer
			}

			width := float64(rand.Intn(200) + 120)

			// Base heights for each layer: Bottom(~450), Middle(~300), Top(~150)
			baseY := float64(g.ScreenHeight - 150 - (layer * 150))

			// Fluctuation
			yDiff := float64(rand.Intn(80) - 40)
			y := baseY + yDiff

			g.Platforms = append(g.Platforms, player.NewPlatform(g.LastPlatformX[layer]+gap, y, width))
			g.LastPlatformX[layer] += gap + width

			if rand.Float64() < 0.25*g.Difficulty {
				enemyX := g.LastPlatformX[layer] - width/2
				r := rand.Float64()
				if r < 0.85 {
					g.Enemies = append(g.Enemies, enemies.NewDinoRobot(enemyX, y-40))
				} else {
					g.Enemies = append(g.Enemies, enemies.NewBigRobot(enemyX, y-80))
				}
			}

			// Bosses spawn much higher up normally or rarely anywhere
			if rand.Float64() < 0.005 {
				g.Enemies = append(g.Enemies, enemies.NewAlienBoss(g.LastPlatformX[layer]-100, y-120))
			}
		}
	}

	var newPlatforms []*player.Platform
	for _, p := range g.Platforms {
		if p.X+p.Width > cameraX-200 {
			newPlatforms = append(newPlatforms, p)
		}
	}
	g.Platforms = newPlatforms

	var newEnemies []interface{}
	for _, e := range g.Enemies {
		var ex float64
		switch ent := e.(type) {
		case *enemies.DinoRobot:
			ex = ent.X
		case *enemies.BigRobot:
			ex = ent.X
		case *enemies.AlienBoss:
			ex = ent.X
		}
		if ex > cameraX-200 {
			newEnemies = append(newEnemies, e)
		}
	}
	g.Enemies = newEnemies

	var newItems []*items.Item
	for _, it := range g.Items {
		if it.X+it.Width > cameraX-200 {
			newItems = append(newItems, it)
		}
	}
	g.Items = newItems
}

func (g *Generator) IncreaseDifficulty() {
	if g.Difficulty < 3.0 {
		g.Difficulty += 0.05
	}
}

func (g *Generator) GetPlatformAABB() []*player.Platform {
	return g.Platforms
}

func (g *Generator) GetPlatformPhysicsAABB() []*physics.AABB {
	aabbs := make([]*physics.AABB, len(g.Platforms))
	for i, p := range g.Platforms {
		aabbs[i] = p.AABB()
	}
	return aabbs
}

func (g *Generator) Draw(screen *ebiten.Image, camX, camY float64) {
	for _, p := range g.Platforms {
		p.Draw(screen, camX, camY)
	}

	for _, it := range g.Items {
		it.Draw(screen, camX, camY)
	}

	for _, e := range g.Enemies {
		switch ent := e.(type) {
		case *enemies.DinoRobot:
			ent.Draw(screen, camX, camY)
		case *enemies.BigRobot:
			ent.Draw(screen, camX, camY)
		case *enemies.AlienBoss:
			ent.Draw(screen, camX, camY)
		}
	}
}

func (g *Generator) GetEnemies() []interface{} {
	return g.Enemies
}

func (g *Generator) RemoveEnemy(idx int) {
	g.Enemies = append(g.Enemies[:idx], g.Enemies[idx+1:]...)
}

func (g *Generator) AddItem(x, y float64, itemType items.ItemType) {
	g.Items = append(g.Items, items.NewItem(x, y, itemType))
}

func (g *Generator) GetItems() []*items.Item {
	return g.Items
}

func (g *Generator) RemoveItem(idx int) {
	g.Items = append(g.Items[:idx], g.Items[idx+1:]...)
}
