package main

type Enemy struct {
    id string
    frame int

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

func (e Enemy)getDrawPosition() Vec2[float64] {
    return Vec2[float64]{x: e.position.x - e.spriteSize.x/2.0, y: e.position.y - e.spriteSize.y/2.0}
}

func (e *Enemy)Init(id string, position Vec2[float64]) {
    e.id = id
    e.frame = 0

    e.position = position
    e.spriteSize = Vec2[float64]{x: 8.0, y: 8.0}
    e.hitbox = Vec2[float64]{x: 8.0, y: 8.0}
}
