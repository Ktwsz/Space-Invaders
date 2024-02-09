package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{
    message string
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
	ebitenutil.DebugPrint(screen, g.message)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Space Invaders")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
