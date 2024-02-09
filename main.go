package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

    img "image"
)

const SpriteSheet = "assets/imgs/spritesheet.png"

type Game struct{
    assetloader AssetLoader
    message string
}

func (g *Game)Init() {
    g.assetloader.Init()

    err := g.assetloader.LoadSpriteSheet(SpriteSheet)
    if err != nil {
		log.Fatal(err)
        return
    }

    g.assetloader.LoadImage("player", img.Rectangle{Min: img.Point{X: 0, Y: 36}, Max: img.Point{X: 11, Y: 8}}, 1)

    g.assetloader.LoadImage("enemy1", img.Rectangle{Min: img.Point{X: 0, Y: 0}, Max: img.Point{X: 8, Y: 8}}, 3)
    g.assetloader.LoadImage("enemy2", img.Rectangle{Min: img.Point{X: 0, Y: 9}, Max: img.Point{X: 8, Y: 8}}, 3)
    g.assetloader.LoadImage("enemy3", img.Rectangle{Min: img.Point{X: 0, Y: 18}, Max: img.Point{X: 18, Y: 8}}, 3)
    g.assetloader.LoadImage("enemy4", img.Rectangle{Min: img.Point{X: 0, Y: 27}, Max: img.Point{X: 18, Y: 8}}, 3)
}

func (g *Game) HandleInputs() {
    if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
        g.message = "left"
    }

    if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
        g.message = "right"
    }
}

func (g *Game) Update() error {
    g.HandleInputs()

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    screen.Clear()
	ebitenutil.DebugPrint(screen, g.message)
    //player := g.assetloader.get("player", 0)
    enemy1_1, err := g.assetloader.get("enemy2", 2)
    if err != nil {
        log.Fatal(err)
        return
    }

    screen.DrawImage(enemy1_1, &ebiten.DrawImageOptions{})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
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
