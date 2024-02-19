package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	img "image"
)

const SpriteSheet = "assets/imgs/spritesheet.png"
const GAME_WIDTH = 150
const GAME_HEIGHT = 150
const H_MARGIN = 25
const V_MARGIN = 10

type Game struct{
    assetloader AssetLoader
    gamestate GameState
}

func (g *Game)Init() {
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

    g.assetloader.LoadSprite("enemy_projectile_1", img.Rectangle{Min: img.Point{X: 0, Y: 52}, Max: img.Point{X: 3, Y: 7}}, 5)
    g.assetloader.LoadSprite("enemy_projectile_2", img.Rectangle{Min: img.Point{X: 0, Y: 60}, Max: img.Point{X: 3, Y: 7}}, 5)
    g.assetloader.LoadSprite("enemy_projectile_3", img.Rectangle{Min: img.Point{X: 0, Y: 67}, Max: img.Point{X: 3, Y: 7}}, 5)

    g.assetloader.LoadSprite("wall", img.Rectangle{Min: img.Point{X: 0, Y: 75}, Max: img.Point{X: 24, Y: 24}}, 1)

    g.gamestate.Init()
}

func (g *Game)ImageToWall() WallBody {
    maskColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
    wallImg, _ := g.assetloader.get("wall", 0)

    sizeX, sizeY := wallImg.Bounds().Dx(), wallImg.Bounds().Dy()

    body := WallBody{}

    for i := range sizeX {
        for j := range sizeY {
            body[i][j] = wallImg.At(i, j) == maskColor
        }
    }

    return body
}

func WallToImage(wall *Wall) *ebiten.Image {
    COLOR_GREEN := color.RGBA{R: 0, G: 255, B: 0, A: 255}
    result := ebiten.NewImage(WALL_SIZE_X, WALL_SIZE_Y)

    for i := range WALL_SIZE_X {
        for j := range WALL_SIZE_Y {
            if wall.body[i][j] {
                result.Set(i, j, COLOR_GREEN)
            }
        }
    }

    return result
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
    if !g.gamestate.wallsBodySet {
        body := g.ImageToWall()
        g.gamestate.SetWallsBody(body)
    }

    if g.gamestate.IsGameRunning() {
        g.HandleInputs()
    }

    g.gamestate.GameLoop()

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    walls := g.gamestate.GetWalls()
    for _, w := range walls {
        wallImg := WallToImage(w)
        op := &ebiten.DrawImageOptions{}
        op.GeoM.Translate(w.position.x - w.hitbox.x/2 + H_MARGIN, w.position.y - w.hitbox.y/2)
        screen.DrawImage(wallImg, op)
    }

    objects := g.gamestate.GetObjectsToDraw()
    for _, entity := range objects {
        entitySprite, err := g.assetloader.get(entity.getId(), entity.getCurrentFrame())
        if err != nil {
            log.Fatal(err)
            return
        }

        op := &ebiten.DrawImageOptions{}
        drawPosition := getDrawPosition(entity)
        op.GeoM.Translate(drawPosition.x + H_MARGIN, drawPosition.y)

        screen.DrawImage(entitySprite, op)
    }


    g.DrawUI(screen)
}

func (g *Game)DrawUI(screen *ebiten.Image) {
    switch g.gamestate.pauseState {
    case GAME_RUNNING:
        g.DrawUIRunning(screen)
    case GAME_OVER:
        g.DrawUIOver(screen)
    case GAME_WIN:
        g.DrawUIWin(screen)
    }
}

func (g *Game)DrawUIRunning(screen *ebiten.Image) {
    livesStr := g.gamestate.GetPlayerLivesStr()
    DrawTextImage(screen, livesStr, 0, GAME_HEIGHT + 20, 1, 0.9)


    scoreStr := g.gamestate.GetScoreStr()
    DrawTextImage(screen, scoreStr, 100, GAME_HEIGHT + 20, 1, 0.9)
}

func (g *Game)DrawUIOver(screen *ebiten.Image) {

}

func (g *Game)DrawUIWin(screen *ebiten.Image) {

}

func DrawTextImage(screen *ebiten.Image, str string, posX, posY, scaleX, scaleY float64) {
    imgOp := &ebiten.DrawImageOptions{}

    imgOp.GeoM.Scale(scaleX, scaleY)
    imgOp.GeoM.Translate(posX, posY)

    text.DrawWithOptions(screen, str, bitmapfont.Face, imgOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return GAME_WIDTH + 2 * H_MARGIN, GAME_HEIGHT + 2 * V_MARGIN
}

func main() {
    factor := 3
	ebiten.SetWindowSize(factor*(GAME_WIDTH + 2 * H_MARGIN), factor*(GAME_HEIGHT + 2 * V_MARGIN))
	ebiten.SetWindowTitle("Space Invaders")
    game := Game{}
    game.Init()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
