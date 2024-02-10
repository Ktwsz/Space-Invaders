package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil"

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

    g.assetloader.LoadSprite("enemy1", img.Rectangle{Min: img.Point{X: 0, Y: 0}, Max: img.Point{X: 8, Y: 8}}, 3)
    g.assetloader.LoadSprite("enemy2", img.Rectangle{Min: img.Point{X: 0, Y: 9}, Max: img.Point{X: 8, Y: 8}}, 3)
    g.assetloader.LoadSprite("enemy3", img.Rectangle{Min: img.Point{X: 0, Y: 18}, Max: img.Point{X: 8, Y: 8}}, 3)
    g.assetloader.LoadSprite("enemy4", img.Rectangle{Min: img.Point{X: 0, Y: 27}, Max: img.Point{X: 8, Y: 8}}, 3)
}

func (g *Game) HandleInputs() {
    if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
        g.gamestate.MovePlayerLeft()
    }

    if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
        g.gamestate.MovePlayerRight()
    }
}

func (g *Game) Update() error {
    g.HandleInputs()

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
        op.GeoM.Translate(entity.getDrawPosition().x, entity.getPosition().y)

        screen.DrawImage(entitySprite, op)
    }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.gamestate.bounds.x, g.gamestate.bounds.y
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
