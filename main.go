package main

import (
	"fmt"
	"game/assets"
	"game/camera"
	"game/enemies"
	"game/level"
	"game/physics"
	"game/player"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	FPS          = 60
)

type Game struct {
	Player      *player.Player
	Level       *level.Generator
	Camera      *camera.Camera
	Projectiles []*player.Projectile
	Score       int
	GameOver    bool
	FrameCount  int
	Background  *ebiten.Image
}

func (g *Game) Restart() {
	g.GameOver = false
	g.Score = 0
	g.FrameCount = 0
	g.Player = player.New(100, 400)
	g.Level = level.New(ScreenWidth, ScreenHeight)
	g.Level.Init(ScreenWidth, ScreenHeight)
	g.Camera = camera.New(ScreenWidth, ScreenHeight)
	g.Projectiles = make([]*player.Projectile, 0)
	
	bgImg, err := ebiten.NewImageFromFile("assets/backgrounds/background.png")
	if err == nil {
		g.Background = bgImg
	}
}

func NewGame() *Game {
	g := &Game{
		Player:     player.New(100, 400),
		Level:      level.New(ScreenWidth, ScreenHeight),
		Camera:     camera.New(ScreenWidth, ScreenHeight),
		Projectiles: make([]*player.Projectile, 0),
		Score:      0,
		GameOver:   false,
		FrameCount: 0,
	}
	g.Level.Init(ScreenWidth, ScreenHeight)
	
	bgImg, err := ebiten.NewImageFromFile("assets/backgrounds/background.png")
	if err == nil {
		g.Background = bgImg
	}
	
	return g
}

func (g *Game) Update() error {
	if g.GameOver {
		return nil
	}

	g.FrameCount++

	keys := map[ebiten.Key]bool{
		ebiten.KeyLeft:  ebiten.IsKeyPressed(ebiten.KeyLeft),
		ebiten.KeyRight: ebiten.IsKeyPressed(ebiten.KeyRight),
		ebiten.KeyUp:    ebiten.IsKeyPressed(ebiten.KeyUp),
		ebiten.KeySpace: ebiten.IsKeyPressed(ebiten.KeySpace),
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) && g.GameOver {
		g.Restart()
		return nil
	}

	g.Player.Update(keys, g.Level.GetPlatformAABB())

	g.Camera.Follow(g.Player.X, g.Player.Y)

	g.Level.Generate(g.Camera.X)

	for _, p := range g.Projectiles {
		p.Update()
	}

	for i := len(g.Projectiles) - 1; i >= 0; i-- {
		p := g.Projectiles[i]
		if p.X < g.Camera.X-100 || p.X > g.Camera.X+float64(ScreenWidth)+100 {
			g.Projectiles = append(g.Projectiles[:i], g.Projectiles[i+1:]...)
			continue
		}

		if p.IsPlayer {
			for j := len(g.Level.Enemies) - 1; j >= 0; j-- {
				var enemyAABB *physics.AABB
				enemy := g.Level.Enemies[j]
				switch ent := enemy.(type) {
				case *enemies.DinoRobot:
					enemyAABB = ent.AABB()
				case *enemies.BigRobot:
					enemyAABB = ent.AABB()
				case *enemies.AlienBoss:
					enemyAABB = ent.AABB()
				}

				if enemyAABB != nil && p.AABB().Intersects(enemyAABB) {
					switch ent := enemy.(type) {
					case *enemies.DinoRobot:
						ent.HP--
						if ent.HP <= 0 {
							g.Score += 10
							g.Level.RemoveEnemy(j)
						}
					case *enemies.BigRobot:
						ent.HP--
						if ent.HP <= 0 {
							g.Score += 50
							g.Level.RemoveEnemy(j)
						}
					case *enemies.AlienBoss:
						ent.HP--
						if ent.HP <= 0 {
							g.Score += 500
							g.Level.RemoveEnemy(j)
						}
					}
					g.Projectiles = append(g.Projectiles[:i], g.Projectiles[i+1:]...)
					break
				}
			}
		} else {
			if p.AABB().Intersects(g.Player.AABB()) && g.Player.Invincible == 0 {
				g.Player.HP--
				g.Player.Invincible = 60
				g.Projectiles = append(g.Projectiles[:i], g.Projectiles[i+1:]...)
			}
		}
	}

	for _, e := range g.Level.Enemies {
		switch ent := e.(type) {
		case *enemies.DinoRobot:
			ent.Update(g.Player.X, g.Level.GetPlatformPhysicsAABB())
			if ent.AABB().Intersects(g.Player.AABB()) && g.Player.Invincible == 0 {
				g.Player.HP--
				g.Player.Invincible = 60
			}
		case *enemies.BigRobot:
			proj := ent.Update(g.Player.X)
			if proj != nil {
				g.Projectiles = append(g.Projectiles, proj)
			}
			if ent.AABB().Intersects(g.Player.AABB()) && g.Player.Invincible == 0 {
				g.Player.HP--
				g.Player.Invincible = 60
			}
		case *enemies.AlienBoss:
			proj := ent.Update(g.Player.X, g.Player.Y)
			if proj != nil {
				g.Projectiles = append(g.Projectiles, proj)
			}
			if ent.AABB().Intersects(g.Player.AABB()) && g.Player.Invincible == 0 {
				g.Player.HP--
				g.Player.Invincible = 60
			}
		}
	}

	if g.Player.Y > float64(ScreenHeight)+100 {
		g.Player.HP = 0
	}

	if g.Player.HP <= 0 {
		g.GameOver = true
	}

	if g.FrameCount%600 == 0 {
		g.Level.IncreaseDifficulty()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.Background != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-g.Camera.X/2, 0)
		screen.DrawImage(g.Background, op)
	} else {
		screen.Fill(color.RGBA{30, 30, 50, 255})
	}

	g.Level.Draw(screen, g.Camera.X, g.Camera.Y)

	g.Player.Draw(screen, g.Camera.X, g.Camera.Y)

	for _, p := range g.Projectiles {
		p.Draw(screen, g.Camera.X, g.Camera.Y)
	}

	scoreText := fmt.Sprintf("HP: %d  Score: %d  Diff: %.1f", g.Player.HP, g.Score, g.Level.Difficulty)
	text.Draw(screen, scoreText, basicfont.Face7x13, 10, 20, color.RGBA{255, 255, 255, 255})

	if g.GameOver {
		msg := "GAME OVER"
		text.Draw(screen, msg, basicfont.Face7x13, ScreenWidth/2-40, ScreenHeight/2, color.RGBA{255, 50, 50, 255})
		msg2 := fmt.Sprintf("Final Score: %d", g.Score)
		text.Draw(screen, msg2, basicfont.Face7x13, ScreenWidth/2-50, ScreenHeight/2+20, color.RGBA{255, 255, 255, 255})
		msg3 := "Press R to Restart"
		text.Draw(screen, msg3, basicfont.Face7x13, ScreenWidth/2-55, ScreenHeight/2+40, color.RGBA{200, 200, 200, 255})
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return ScreenWidth, ScreenHeight
}

var whiteColor = func(x, y int) color.Color { return color.White }

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Iron Platformer")

	assets.CreateSprites()
	assets.CreateBackground()
	assets.CreateSounds()

	game := NewGame()

	ebiten.RunGame(game)
}
