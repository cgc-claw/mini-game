package main

import (
	"fmt"
	"game/assets"
	"game/camera"
	"game/enemies"
	"game/items"
	"game/level"
	"game/physics"
	"game/player"
	"game/sprites"
	"image/color"
	"io/ioutil"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type DifficultyMode int

const (
	ModeEasy DifficultyMode = iota
	ModeMedium
	ModeHard
	ModeHardcore
)

type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StateGameOver
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	FPS          = 60
)

type Game struct {
	Player        *player.Player
	Level         *level.Generator
	Camera        *camera.Camera
	Projectiles   []*player.Projectile
	Score         int
	GameOver      bool
	State         GameState
	HighScore     int
	FrameCount    int
	BgFixed       *ebiten.Image
	BgMoving      *ebiten.Image
	Mode          DifficultyMode
	Stage         int
	MovementTimer int
}

func loadHighScore() int {
	data, err := ioutil.ReadFile("highscore.txt")
	if err == nil {
		hs, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err == nil {
			return hs
		}
	}
	return 0
}

func saveHighScore(score int) {
	_ = ioutil.WriteFile("highscore.txt", []byte(strconv.Itoa(score)), 0644)
}

func (g *Game) Restart() {
	g.State = StatePlaying
	g.Score = 0
	g.FrameCount = 0
	g.Player = player.New(100, 400)
	g.Level = level.New(ScreenWidth, ScreenHeight)
	g.Level.Init(ScreenWidth, ScreenHeight)
	g.Camera = camera.New(ScreenWidth, ScreenHeight)
	g.Projectiles = make([]*player.Projectile, 0)
	g.Stage = 1
	g.MovementTimer = 0

	// Set baseline difficulty based on mode
	switch g.Mode {
	case ModeEasy:
		g.Level.Difficulty = 0.5
	case ModeMedium:
		g.Level.Difficulty = 1.0
	case ModeHard:
		g.Level.Difficulty = 2.0
	case ModeHardcore:
		g.Level.Difficulty = 3.0
	}
	bgFixed, _, err1 := ebitenutil.NewImageFromFile("assets/backgrounds/bg_fixed.png")
	if err1 == nil {
		g.BgFixed = bgFixed
	}
	bgMoving, _, err2 := ebitenutil.NewImageFromFile("assets/backgrounds/bg_moving.png")
	if err2 == nil {
		g.BgMoving = bgMoving
	}
}

func NewGame() *Game {
	g := &Game{
		Player:        player.New(100, 400),
		Level:         level.New(ScreenWidth, ScreenHeight),
		Camera:        camera.New(ScreenWidth, ScreenHeight),
		Projectiles:   make([]*player.Projectile, 0),
		Score:         0,
		State:         StateMenu,
		HighScore:     loadHighScore(),
		FrameCount:    0,
		Mode:          ModeMedium,
		Stage:         1,
		MovementTimer: 0,
	}
	g.Level.Init(ScreenWidth, ScreenHeight)
	bgFixed, _, err1 := ebitenutil.NewImageFromFile("assets/backgrounds/bg_fixed.png")
	if err1 == nil {
		g.BgFixed = bgFixed
	}
	bgMoving, _, err2 := ebitenutil.NewImageFromFile("assets/backgrounds/bg_moving.png")
	if err2 == nil {
		g.BgMoving = bgMoving
	}

	return g
}

func (g *Game) Update() error {
	if g.State == StateMenu {
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			g.Mode = ModeEasy
		} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
			g.Mode = ModeMedium
		} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
			g.Mode = ModeHard
		} else if inpututil.IsKeyJustPressed(ebiten.Key4) {
			g.Mode = ModeHardcore
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			assets.PlaySound("restart")
			g.Restart()
		}
		return nil
	}

	if g.State == StateGameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			assets.PlaySound("restart")
			g.Restart()
		}
		return nil
	}

	assets.PlayMusic("music")

	g.FrameCount++

	gamepadIDs := ebiten.GamepadIDs()
	var gpadLeft, gpadRight, gpadJump, gpadShoot bool
	if len(gamepadIDs) > 0 {
		id := gamepadIDs[0]
		gpadLeft = ebiten.GamepadAxisValue(id, 0) < -0.5
		gpadRight = ebiten.GamepadAxisValue(id, 0) > 0.5
		gpadJump = inpututil.IsGamepadButtonJustPressed(id, ebiten.GamepadButton0)  // Bottom face
		gpadShoot = inpututil.IsGamepadButtonJustPressed(id, ebiten.GamepadButton2) // Left face
	}

	keys := map[ebiten.Key]bool{
		ebiten.KeyLeft:  ebiten.IsKeyPressed(ebiten.KeyLeft) || gpadLeft,
		ebiten.KeyRight: ebiten.IsKeyPressed(ebiten.KeyRight) || gpadRight,
		ebiten.KeyUp:    inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeySpace) || gpadJump,
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyZ) || gpadShoot {
		proj := g.Player.Shoot()
		if proj != nil {
			g.Projectiles = append(g.Projectiles, proj)
		}
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
							if rand.Float64() < 0.2 {
								g.Level.AddItem(ent.X, ent.Y, items.TypeAmmo)
							} else if rand.Float64() < 0.05 {
								g.Level.AddItem(ent.X, ent.Y, items.TypeHealth)
							}
							g.Level.RemoveEnemy(j)
						}
					case *enemies.BigRobot:
						ent.HP--
						if ent.HP <= 0 {
							g.Score += 50
							if rand.Float64() < 0.4 {
								g.Level.AddItem(ent.X, ent.Y, items.TypeAmmo)
							} else if rand.Float64() < 0.15 {
								g.Level.AddItem(ent.X, ent.Y, items.TypeHealth)
							}
							g.Level.RemoveEnemy(j)
						}
					case *enemies.AlienBoss:
						ent.HP--
						if ent.HP <= 0 {
							g.Score += 500
							g.Level.AddItem(ent.X, ent.Y, items.TypeLife)
							g.Level.RemoveEnemy(j)
						}
					}
					g.Projectiles = append(g.Projectiles[:i], g.Projectiles[i+1:]...)
					break
				}
			}
		} else {
			if p.AABB().Intersects(g.Player.AABB()) && g.Player.Invincible == 0 {
				bulletDmg := 10
				if g.Mode == ModeMedium {
					bulletDmg = 15
				} else if g.Mode == ModeHard {
					bulletDmg = 25
				} else if g.Mode == ModeHardcore {
					bulletDmg = 50
				}
				g.Player.HP -= bulletDmg
				g.Player.Invincible = 60
				g.Projectiles = append(g.Projectiles[:i], g.Projectiles[i+1:]...)
			}
		}
	}

	for _, e := range g.Level.Enemies {
		touchDmg := 25
		if g.Mode != ModeEasy {
			touchDmg = 50
		}
		switch ent := e.(type) {
		case *enemies.DinoRobot:
			ent.Update(g.Player.X, g.Level.GetPlatformPhysicsAABB())
			if ent.AABB().Intersects(g.Player.AABB()) && g.Player.Invincible == 0 {
				g.Player.HP -= touchDmg
				g.Player.Invincible = 60
			}
		case *enemies.BigRobot:
			proj := ent.Update(g.Player.X, g.Level.GetPlatformPhysicsAABB())
			if proj != nil {
				g.Projectiles = append(g.Projectiles, proj)
			}
			if ent.AABB().Intersects(g.Player.AABB()) && g.Player.Invincible == 0 {
				g.Player.HP -= touchDmg
				g.Player.Invincible = 60
			}
		case *enemies.AlienBoss:
			proj := ent.Update(g.Player.X, g.Player.Y)
			if proj != nil {
				g.Projectiles = append(g.Projectiles, proj)
			}
			if ent.AABB().Intersects(g.Player.AABB()) && g.Player.Invincible == 0 {
				g.Player.HP -= touchDmg
				g.Player.Invincible = 60
			}
		}
	}

	for i := len(g.Level.Items) - 1; i >= 0; i-- {
		it := g.Level.Items[i]
		it.Update(g.Level.GetPlatformPhysicsAABB())
		if it.AABB().Intersects(g.Player.AABB()) {
			switch it.Type {
			case items.TypeHealth:
				g.Player.HP += 25
				if g.Player.HP > 50 {
					g.Player.HP = 50
				}
			case items.TypeAmmo:
				g.Player.Bullets += 10
			case items.TypeLife:
				g.Player.Lives++
			}
			assets.PlaySound("hit") // Reusing hit sound for pickup
			g.Level.RemoveItem(i)
		}
	}

	if g.Player.Y > float64(ScreenHeight)+50 { // Touch Lava
		g.Player.HP = 0
		g.Player.Lives = -1 // Instant Game Over
	}

	if g.Player.HP <= 0 {
		assets.PlaySound("death")
		g.Player.Lives--
		if g.Player.Lives < 0 {
			g.State = StateGameOver
			if g.Score > g.HighScore {
				g.HighScore = g.Score
				saveHighScore(g.HighScore)
			}
		} else {
			// Respawn
			g.Player.HP = 50
			g.Player.Bullets = 10
			g.Player.Invincible = 120                   // 2 seconds of iframes
			g.Player.X = g.Level.LastPlatformX[0] - 400 // Safeish spawn
			g.Player.Y = 200
			g.Player.VX = 0
			g.Player.VY = 0
		}
	}

	if g.Player.VX != 0 {
		g.MovementTimer++
	}

	if g.MovementTimer >= 3*60*60 { // 3 minutes at 60 FPS
		g.Stage++
		g.MovementTimer = 0
		assets.PlaySound("hit") // Sound for stage up

		var multiplier float64
		switch g.Mode {
		case ModeEasy:
			multiplier = 1.0
		case ModeMedium:
			multiplier = 1.5
		case ModeHard:
			multiplier = 2.0
		case ModeHardcore:
			multiplier = 2.0 * math.Log2(float64(g.Stage+1))
		}
		g.Level.Difficulty *= multiplier
	}

	// if g.FrameCount%600 == 0 { // Removed as per instruction
	// 	g.Level.IncreaseDifficulty()
	// }

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.BgFixed != nil {
		bgW := g.BgFixed.Bounds().Dx()
		bgH := g.BgFixed.Bounds().Dy()
		scale := float64(ScreenHeight) / float64(bgH)
		scaledW := float64(bgW) * scale

		offsetX := math.Mod(g.Camera.X*0.05, scaledW)
		if offsetX < 0 {
			offsetX += scaledW
		}

		for i := -1; i <= int(ScreenWidth/scaledW)+2; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(-offsetX+float64(i)*scaledW, 0)
			screen.DrawImage(g.BgFixed, op)
		}
	} else {
		screen.Fill(color.RGBA{10, 10, 20, 255})
	}

	if g.BgMoving != nil {
		bgW := g.BgMoving.Bounds().Dx()
		bgH := g.BgMoving.Bounds().Dy()
		scale := float64(ScreenHeight) / float64(bgH)
		scaledW := float64(bgW) * scale

		offsetX := math.Mod(g.Camera.X*0.2, scaledW)
		if offsetX < 0 {
			offsetX += scaledW
		}

		for i := -1; i <= int(ScreenWidth/scaledW)+2; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(-offsetX+float64(i)*scaledW, 0)
			screen.DrawImage(g.BgMoving, op)
		}
	}

	g.Level.Draw(screen, g.Camera.X, g.Camera.Y)

	g.Player.Draw(screen, g.Camera.X, g.Camera.Y)

	// Draw Lava floor
	lavaY := float64(ScreenHeight + 30)
	offsetX := math.Mod(g.Camera.X*1.0, 32.0)
	if offsetX < 0 {
		offsetX += 32.0
	}

	for yOff := lavaY; yOff < float64(ScreenHeight+200); yOff += 32.0 {
		for i := -1; i <= int(ScreenWidth/32)+2; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(-offsetX+float64(i)*32.0, yOff-g.Camera.Y)
			screen.DrawImage(sprites.LavaSprite, op)
		}
	}

	for _, p := range g.Projectiles {
		p.Draw(screen, g.Camera.X, g.Camera.Y)
	}

	modeStr := "Medium"
	switch g.Mode {
	case ModeEasy:
		modeStr = "Easy"
	case ModeHard:
		modeStr = "Hard"
	case ModeHardcore:
		modeStr = "Hardcore"
	}
	scoreText := fmt.Sprintf("Stage: %d (%s)  Lives: %d  HP: %d  Ammo: %d  Score: %d  Diff: %.1f", g.Stage, modeStr, g.Player.Lives, g.Player.HP, g.Player.Bullets, g.Score, g.Level.Difficulty)
	text.Draw(screen, scoreText, basicfont.Face7x13, 10, 20, color.RGBA{255, 255, 255, 255})

	hsText := fmt.Sprintf("High Score: %d", g.HighScore)
	text.Draw(screen, hsText, basicfont.Face7x13, ScreenWidth-150, 20, color.RGBA{255, 255, 0, 255})

	if g.State == StateMenu {
		msg := "IRON PLATFORMER"
		text.Draw(screen, msg, basicfont.Face7x13, ScreenWidth/2-55, ScreenHeight/2-30, color.RGBA{255, 255, 255, 255})
		msgSelect := "Press 1-Easy 2-Medium 3-Hard 4-Hardcore"
		text.Draw(screen, msgSelect, basicfont.Face7x13, ScreenWidth/2-140, ScreenHeight/2-10, color.RGBA{200, 200, 255, 255})

		modeStr := "Medium"
		switch g.Mode {
		case ModeEasy:
			modeStr = "Easy"
		case ModeHard:
			modeStr = "Hard"
		case ModeHardcore:
			modeStr = "Hardcore"
		}
		msgMode := fmt.Sprintf("Current Mode: %s", modeStr)
		text.Draw(screen, msgMode, basicfont.Face7x13, ScreenWidth/2-70, ScreenHeight/2+10, color.RGBA{255, 255, 0, 255})

		msg2 := "Press SPACE or ENTER to Start"
		text.Draw(screen, msg2, basicfont.Face7x13, ScreenWidth/2-100, ScreenHeight/2+40, color.RGBA{200, 200, 200, 255})
	} else if g.State == StateGameOver {
		msg := "GAME OVER"
		text.Draw(screen, msg, basicfont.Face7x13, ScreenWidth/2-40, ScreenHeight/2, color.RGBA{255, 50, 50, 255})
		msg2 := fmt.Sprintf("Final Score: %d", g.Score)
		text.Draw(screen, msg2, basicfont.Face7x13, ScreenWidth/2-50, ScreenHeight/2+20, color.RGBA{255, 255, 255, 255})
		if g.Score >= g.HighScore && g.Score > 0 {
			msgNew := "NEW HIGH SCORE!"
			text.Draw(screen, msgNew, basicfont.Face7x13, ScreenWidth/2-55, ScreenHeight/2-30, color.RGBA{0, 255, 0, 255})
		}
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
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	assets.CreateSprites()
	assets.CreateBackground()
	assets.CreateSounds()
	assets.LoadAudio()

	game := NewGame()

	ebiten.RunGame(game)
}
