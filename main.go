package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

    img "image"
)

const SpriteSheet = "assets/imgs/spritesheet.png"

type Game struct{
    assetloader AssetLoader
    gamestate GameState
}

func (g *Game)Init() {
    g.gamestate.Init()

    g.assetloader.Init()

    err := g.assetloader.LoadSpriteSheet(SpriteSheet)
    if err != nil {
		log.Fatal(err)
        return
    }

    g.assetloader.LoadSprite("player", img.Rectangle{Min: img.Point{X: 0, Y: 36}, Max: img.Point{X: 11, Y: 8}}, 1)

    g.assetloader.LoadSpriteWithDeath("enemy1", img.Rectangle{Min: img.Point{X: 0, Y: 0}, Max: img.Point{X: 8, Y: 8}}, 4, Vec2[int]{x: 8, y: 8})
    g.assetloader.LoadSpriteWithDeath("enemy2", img.Rectangle{Min: img.Point{X: 0, Y: 9}, Max: img.Point{X: 8, Y: 8}}, 4, Vec2[int]{x: 8, y: 8})
    g.assetloader.LoadSpriteWithDeath("enemy3", img.Rectangle{Min: img.Point{X: 0, Y: 18}, Max: img.Point{X: 8, Y: 8}}, 4, Vec2[int]{x: 8, y: 8})
    g.assetloader.LoadSpriteWithDeath("enemy4", img.Rectangle{Min: img.Point{X: 0, Y: 27}, Max: img.Point{X: 8, Y: 8}}, 4, Vec2[int]{x: 8, y: 8})

    g.assetloader.LoadSpriteWithDeath("player_projectile", img.Rectangle{Min: img.Point{X: 0, Y: 45}, Max: img.Point{X: 1, Y: 6}}, 2, Vec2[int]{x: 6, y: 6})

    g.assetloader.LoadSprite("enemy_projectile_1", img.Rectangle{Min: img.Point{X: 0, Y: 52}, Max: img.Point{X: 3, Y: 7}}, 3)
    g.assetloader.LoadSprite("enemy_projectile_2", img.Rectangle{Min: img.Point{X: 0, Y: 60}, Max: img.Point{X: 3, Y: 7}}, 3)
    g.assetloader.LoadSprite("enemy_projectile_3", img.Rectangle{Min: img.Point{X: 0, Y: 67}, Max: img.Point{X: 3, Y: 7}}, 3)
}

func (g *Game) HandleInputs() {
    if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
        g.gamestate.PlayerMoveLeft()
    }

    if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
        g.gamestate.PlayerMoveRight()
    }

    if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
        g.gamestate.PlayerShoot()
    }
}

func (g *Game) Update() error {
    g.gamestate.removeDeadEnemies()
    g.gamestate.RemoveDeadProjectiles()
    g.gamestate.CheckForMissedProjectiles()
    g.gamestate.CheckEnemiesInBounds()
    g.HandleInputs()
    g.gamestate.MoveProjectiles()
    g.gamestate.HandleCollisions()

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    objects := g.gamestate.GetObjectsToDraw()
    for _, entity := range objects {
        entitySprite, err := g.assetloader.get(entity.getId(), entity.getCurrentFrame())
        if err != nil {
            log.Fatal(err)
            return
        }

        op := &ebiten.DrawImageOptions{}
        drawPosition := getDrawPosition(entity)
        op.GeoM.Translate(drawPosition.x, drawPosition.y)

        screen.DrawImage(entitySprite, op)
    }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return 150, 120
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Space Invaders")
    game := Game{}
    game.Init()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
