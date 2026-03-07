package level

import (
	"game/enemies"
	"game/physics"
	"game/player"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Generator struct {
	Platforms     []*player.Platform
	Enemies       []interface{}
	LastPlatformX float64
	ScreenWidth   int
	ScreenHeight  int
	Difficulty    float64
}

func New(w, h int) *Generator {
	return &Generator{
		Platforms:     make([]*player.Platform, 0),
		Enemies:       make([]interface{}, 0),
		LastPlatformX: 100,
		ScreenWidth:   w,
		ScreenHeight:  h,
		Difficulty:    1.0,
	}
}

func (g *Generator) Init(screenW, screenH int) {
	g.Platforms = append(g.Platforms, player.NewPlatform(0, float64(screenH-40), 400))
	g.LastPlatformX = 400
}

func (g *Generator) Generate(cameraX float64) {
	generateUntil := cameraX + float64(g.ScreenWidth) + 500

	for g.LastPlatformX < generateUntil {
		// Gap size based on randomness: small, medium, or large (requiring double jump)
		gapType := rand.Float64()
		var gap float64
		if gapType < 0.6 {
			gap = float64(rand.Intn(50) + 50) // small gap
		} else if gapType < 0.9 {
			gap = float64(rand.Intn(80) + 100) // medium gap
		} else {
			gap = float64(rand.Intn(60) + 180) // large gap (double jump)
		}

		width := float64(rand.Intn(200) + 120)

		// Coherent height: relative to last platform
		var lastY float64
		if len(g.Platforms) > 0 {
			lastY = g.Platforms[len(g.Platforms)-1].Y
		} else {
			lastY = float64(g.ScreenHeight - 100)
		}

		// Change height by -80 to +80
		yDiff := float64(rand.Intn(160) - 80)
		y := lastY + yDiff

		// Keep within screen bounds
		if y < 200 {
			y = 200
		} else if y > float64(g.ScreenHeight-100) {
			y = float64(g.ScreenHeight - 100)
		}

		g.Platforms = append(g.Platforms, player.NewPlatform(g.LastPlatformX+gap, y, width))
		g.LastPlatformX += gap + width

		if rand.Float64() < 0.25*g.Difficulty {
			enemyX := g.LastPlatformX - width/2
			r := rand.Float64()
			if r < 0.7 {
				g.Enemies = append(g.Enemies, enemies.NewDinoRobot(enemyX, y-40))
			} else {
				g.Enemies = append(g.Enemies, enemies.NewBigRobot(enemyX, y-80))
			}
		}

		if rand.Float64() < 0.01 {
			g.Enemies = append(g.Enemies, enemies.NewAlienBoss(g.LastPlatformX-100, y-120))
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
