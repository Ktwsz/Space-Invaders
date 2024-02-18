package main

type Player struct {
    lives int

    position Vec2[float64]

    hitbox Vec2[float64]
    hitboxSendMask uint8
    hitboxReceiveMask uint8

    spriteSize Vec2[float64]

    shotOnCooldown bool
    collideMap map[Entity]bool
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
    p.lives = 3

    p.spriteSize = Vec2[float64]{x: 11.0, y: 8.0}
    p.hitbox = p.spriteSize
    p.position = Vec2[float64]{x: float64(bounds.x)/2.0, y: float64(bounds.y) - p.spriteSize.y}

    p.hitboxReceiveMask = HITBOX_ENEMY
    p.hitboxSendMask = HITBOX_PLAYER
}

func (p Player)getEntityType() int {
    return ENTITY_PLAYER
}

func (p Player)getGamestateIx() int {
    return 0
}

func (p Player)getHitboxSendMask() uint8 {
    return p.hitboxSendMask
}

func (p Player)getHiboxReceiveMask() uint8 {
    return p.hitboxReceiveMask
}

func (p *Player)Hit() {
    p.lives -= 1
}

func (p Player)didCollideWith(ent Entity) bool {
    _, exists := p.collideMap[ent]
    return exists
}
