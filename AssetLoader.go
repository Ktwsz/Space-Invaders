package main

import (
    "errors"
	"github.com/hajimehoshi/ebiten/v2"
    img "image"
    png "image/png"
    "os"
)

type AssetLoader struct {
    spritesheet img.NRGBA
    assets map[string][]*ebiten.Image
}

func (assetloader *AssetLoader)Init() {
    assetloader.assets = map[string][]*ebiten.Image{}
}

func (assetloader *AssetLoader)LoadSpriteSheet(path string) error {
    reader, err := os.Open(path)
    if err != nil {
        return err
    }
    defer reader.Close()

    m, err := png.Decode(reader)
    if err != nil {
        return err
    }

    if spitesheetNRGBA, ok := m.(*img.NRGBA); ok {
        assetloader.spritesheet = *spitesheetNRGBA
    } else {
        return errors.New("failed to convert image to RGBA")
    }

    return nil
}

func (assetloader *AssetLoader)LoadSprite(name string, bounds img.Rectangle, count int) {
    sizeX, sizeY := bounds.Max.X, bounds.Max.Y
    pos := Vec2[int]{x: bounds.Min.X, y: bounds.Min.Y} 

    images := make([]*ebiten.Image, count)

    for i := range count {
        imageRect := img.Rectangle{Min: img.Point{X: pos.x, Y: pos.y}, Max: img.Point{X: pos.x + sizeX, Y: pos.y + sizeY}}
        subImage := assetloader.spritesheet.SubImage(imageRect) 
        images[i] = ebiten.NewImageFromImage(subImage)

        pos.x += sizeX + 1
    }

    assetloader.assets[name] = images
}

func (assetloader *AssetLoader)LoadSpriteWithDeath(name string, bounds img.Rectangle, count int, deathFrameBounds Vec2[int]) {
    assetloader.LoadSprite(name, bounds, count-1)

    sizeX := bounds.Max.X
    posX := bounds.Min.X + (count-1) * (sizeX+1)
    posY := bounds.Min.Y

    imageRect := img.Rectangle{Min: img.Point{X: posX, Y: posY}, Max: img.Point{X: posX + deathFrameBounds.x, Y: posY + deathFrameBounds.y}}

    subImage := assetloader.spritesheet.SubImage(imageRect)
    assetloader.assets[name] = append(assetloader.assets[name], ebiten.NewImageFromImage(subImage))
}

func (assetloader *AssetLoader)get(name string, frame int) (*ebiten.Image, error) {
    imgs := assetloader.assets[name]
    if imgs == nil {
        return nil, errors.New("no sprite found")
    }
    if frame >= len(imgs) {
        return nil, errors.New("frame out of sprite length")
    }
    return imgs[frame], nil
}
