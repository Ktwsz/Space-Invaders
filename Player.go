package main

type Player struct {
    position Vec2[float64]
    hitbox Vec2[float64]
    spriteSize Vec2[float64]

    shotOnCooldown bool
}

func (p Player)getId() string {
    return "player"
}

func (p Player)getCurrentFrame() int {
    return 0
}

func (p Player)getPosition() Vec2[float64] {
    return p.position
}

func (p Player)getSpriteSize() Vec2[float64] {
    return p.spriteSize
}

func (p Player)getHitbox() Vec2[float64] {
    return p.hitbox
}

func (p *Player)Init(bounds Vec2[int]) {
    p.spriteSize = Vec2[float64]{x: 11.0, y: 8.0}
    p.hitbox = p.spriteSize
    p.position = Vec2[float64]{x: float64(bounds.x)/2.0, y: float64(bounds.y) - p.spriteSize.y}
}

