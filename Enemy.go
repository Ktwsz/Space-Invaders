package main

type Enemy struct {
    id string
    frame int
    frameCount int

    rowNumber int
    position Vec2[float64]
    hitbox Vec2[float64]
    spriteSize Vec2[float64]
}

func (e Enemy)getId() string {
    return e.id
}

func (e Enemy)getCurrentFrame() int {
    return e.frame
}

func (e Enemy)getPosition() Vec2[float64] {
    return e.position
}

func (e Enemy)getSpriteSize() Vec2[float64] {
    return e.spriteSize
}

func (e *Enemy)Init(id string, frameCount int, position Vec2[float64], rowNumber int) {
    e.id = id
    e.frameCount = frameCount

    e.rowNumber = rowNumber
    e.position = position
    e.spriteSize = Vec2[float64]{x: 8.0, y: 8.0}
    e.hitbox = Vec2[float64]{x: 8.0, y: 8.0}
}

func (e *Enemy)Move(speed float64) {
    e.position.y += speed
    e.frame = (e.frame + 1) % e.frameCount
}

func (e *Enemy)Shift(shift Vec2[float64]) {
    e.position = e.position.add(shift)
}
