package main

import (
	"image/color"
	"log"
    "fmt"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
    "github.com/hajimehoshi/ebiten/v2/audio"
    "github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/vector"

	img "image"
    
    "bytes"
)

const SpriteSheet = "assets/imgs/spritesheet.png"
const GAME_WIDTH = 150
const GAME_HEIGHT = 150
const H_MARGIN = 25
const V_MARGIN = 10

const SCREEN_WIDTH = GAME_WIDTH + 2 * H_MARGIN
const SCREEN_HEIGHT = GAME_HEIGHT + 2 * V_MARGIN

const sampleRate = 44000

const (
    SCREEN_MAIN = iota
    SCREEN_OPT
)

type Game struct{
    assetloader AssetLoader
    gamestate GameState
    audioContext *audio.Context
    startScreen int
    volume128 int
}

func (g *Game)Init() {
    g.assetloader.Init()
    g.audioContext = audio.NewContext(sampleRate)
    g.volume128 = 10//for the sake of our ears

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

    sound_dir := "assets/sounds/"
    sounds := [5]string{"player_hit", "player_shoot", "enemy_die", "enemy_shoot", "enemy_move"}
    for _, sound := range sounds {
        sound_path := fmt.Sprintf("%s%s.wav", sound_dir, sound)
        g.assetloader.LoadSound(sound_path, sound)
    }

    g.gamestate.Init()
}

func (g *Game)ImageToWall() WallBody {
    maskColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
    wallImg, _ := g.assetloader.getSprite("wall", 0)

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
    if g.gamestate.IsGameRunning() {
        if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
            g.gamestate.PlayerMoveLeft()
        }

        if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
            g.gamestate.PlayerMoveRight()
        }

        if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
            g.gamestate.PlayerShoot()
        }
    } else {
        switch g.startScreen {
        case SCREEN_MAIN:
            if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
                g.gamestate.StartGame()
            }
            if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
                g.startScreen = SCREEN_OPT
            }
        case SCREEN_OPT:
            if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
                g.startScreen = SCREEN_MAIN
            }
            if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
                g.volume128--;
            }

            if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
                g.volume128++;
                if g.volume128 > 128 {
                    g.volume128 = 128
                }
                if g.volume128 < 0 {
                    g.volume128 = 0
                }
            }
        }
    }
}

func (g *Game) Update() error {
    g.HandleInputs()

    if !g.gamestate.IsGameRunning() {
        return nil
    }

    if !g.gamestate.wallsBodySet {
        body := g.ImageToWall()
        g.gamestate.SetWallsBody(body)
    }

    g.gamestate.GameLoop()
    g.PlaySounds()

    return nil
}

func (g *Game)PlaySounds() {
    soundsQueue := g.gamestate.GetSoundQueue()
    g.gamestate.ClearSoundQueue()

    volume := float64(g.volume128) / 128

    for _, sound_name := range soundsQueue {
        soundBytes := g.assetloader.getSound(sound_name)

        sound, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewBuffer(soundBytes))
        if err != nil {
            return
        }

        player, err := g.audioContext.NewPlayer(sound)
        if err != nil {
            return
        }

        player.SetVolume(volume)
        player.Play()
    } 
}

func (g *Game) Draw(screen *ebiten.Image) {
    if g.gamestate.IsGameRunning() {
        g.DrawGameObjects(screen)
    }

    g.DrawUI(screen)
}

func (g *Game)DrawGameObjects(screen *ebiten.Image) {
    walls := g.gamestate.GetWalls()
    for _, w := range walls {
        wallImg := WallToImage(w)
        op := &ebiten.DrawImageOptions{}
        op.GeoM.Translate(w.position.x - w.hitbox.x/2 + H_MARGIN, w.position.y - w.hitbox.y/2)
        screen.DrawImage(wallImg, op)
    }

    objects := g.gamestate.GetObjectsToDraw()
    for _, entity := range objects {
        entitySprite, err := g.assetloader.getSprite(entity.getId(), entity.getCurrentFrame())
        if err != nil {
            log.Fatal(err)
            return
        }

        op := &ebiten.DrawImageOptions{}
        drawPosition := getDrawPosition(entity)
        op.GeoM.Translate(drawPosition.x + H_MARGIN, drawPosition.y)

        screen.DrawImage(entitySprite, op)
    }
}

func (g *Game)DrawUI(screen *ebiten.Image) {
    switch g.gamestate.pauseState {
    case GAME_STARTING:
        g.DrawUIStarting(screen)
    case GAME_RUNNING:
        g.DrawUIRunning(screen)
    case GAME_OVER:
        g.DrawUIOver(screen)
    case GAME_WIN:
        g.DrawUIWin(screen)
    }
}

func (g *Game)DrawUIStarting(screen *ebiten.Image) {
    title := "Space Invaders"
    OptionsStr := "Press Space for controls"
    OptionsStrcd := "and options"
    startGameStr := "Press Enter to start"

    optionsTitle := "Options"
    optionsControlsStr := "Controls: "
    optionsLeft := "Move left: Arrow Left"
    optionsRight := "Move right: Arrow Right"
    optionsShoot := "Shoot: Arrow Up"
    optionsVolStr := "Adjust volume with"
    optionsVolStrcd := "left and right arrows"
    optionsReturn := "Press space to go back"

    switch g.startScreen {
    case SCREEN_MAIN:
        DrawTextImage(screen, title, SCREEN_WIDTH/ 2 - 40, 20, 1, 1)
        DrawTextImage(screen, OptionsStr, 10, 50, 1, 1)
        DrawTextImage(screen, OptionsStrcd, 10, 60, 1, 1)
        DrawTextImage(screen, startGameStr, SCREEN_WIDTH/ 2 - 75, SCREEN_HEIGHT - 10, 1.2, 1.2)
    case SCREEN_OPT:
        DrawTextImage(screen, optionsTitle, SCREEN_WIDTH/ 2 - 30, 20, 1, 1)
        DrawTextImage(screen, optionsControlsStr, 10, 30, 1, 1)
        DrawTextImage(screen, optionsLeft, 15, 43, 1, 1)
        DrawTextImage(screen, optionsRight, 15, 56, 1, 1)
        DrawTextImage(screen, optionsShoot, 15, 69, 1, 1)
        DrawTextImage(screen, optionsVolStr, 10, 100, 1, 1)
        DrawTextImage(screen, optionsVolStrcd, 10, 113, 1, 1)
        DrawVolumeBar(screen, g.volume128, 10, 120, SCREEN_WIDTH- 20, 10)
        DrawTextImage(screen, optionsReturn, SCREEN_WIDTH / 2 - 75, SCREEN_HEIGHT - 10, 1, 1)
    }
}

func DrawVolumeBar(screen *ebiten.Image, volume128, x, y, w, h int) {
    COLOR_FILL := color.RGBA{R: 0, G: 255, B: 0, A: 255}
    COLOR_EMPTY := color.RGBA{R: 255, G: 255, B: 255, A: 255}
    fillPerc := float64(volume128) / 128
    fillWidth := float32(fillPerc) * float32(w)

    vector.DrawFilledRect(screen, float32(x), float32(y), fillWidth, float32(h), COLOR_FILL, true)
    vector.DrawFilledRect(screen, float32(x) + fillWidth, float32(y), float32(w) - fillWidth, float32(h), COLOR_EMPTY, true)
}

func (g *Game)DrawUIRunning(screen *ebiten.Image) {
    livesStr := g.gamestate.GetPlayerLivesStr()
    DrawTextImage(screen, livesStr, 0, GAME_HEIGHT + 20, 1, 0.9)


    scoreStr := g.gamestate.GetScoreStr()
    DrawTextImage(screen, scoreStr, 100, GAME_HEIGHT + 20, 1, 0.9)
}

func (g *Game)DrawUIOver(screen *ebiten.Image) {
    DrawTextImage(screen, "Game Over! :(", GAME_WIDTH / 2 - 20, GAME_HEIGHT / 2, 1, 1)
    DrawTextImage(screen, "To try again reopen the app.", GAME_WIDTH / 2 - 60, GAME_HEIGHT / 2 + 20, 1, 1)
}

func (g *Game)DrawUIWin(screen *ebiten.Image) {
    scoreResultStr := g.gamestate.GetScoreResultStr()
    DrawTextImage(screen, "GG! Thanks for playing", GAME_WIDTH / 2 - 40, GAME_HEIGHT / 2, 1, 1)
    DrawTextImage(screen, scoreResultStr, GAME_WIDTH / 2 - 60, GAME_HEIGHT / 2 + 20, 1, 1)
    DrawTextImage(screen, "To try again reopen the app.", GAME_WIDTH / 2 - 60, GAME_HEIGHT / 2 + 40, 1, 1)
}

func DrawTextImage(screen *ebiten.Image, str string, posX, posY, scaleX, scaleY float64) {
    imgOp := &ebiten.DrawImageOptions{}

    imgOp.GeoM.Scale(scaleX, scaleY)
    imgOp.GeoM.Translate(posX, posY)

    text.DrawWithOptions(screen, str, bitmapfont.Face, imgOp)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
    return SCREEN_WIDTH, SCREEN_HEIGHT
}

func main() {
    factor := 3
	ebiten.SetWindowSize(factor*SCREEN_WIDTH, factor*SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Space Invaders")
    game := Game{}
    game.Init()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
