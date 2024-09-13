package entity

import (
    "space_invaders/utils"
    "space_invaders/entity/hitbox"
    "space_invaders/entity/ids"
)

type Player struct {
    lives int

    Position utils.Vec2[float64]

    hitbox utils.Vec2[float64]
    hitboxSendMask uint8
    hitboxReceiveMask uint8

    spriteSize utils.Vec2[float64]

    ShotOnCooldown bool
    HandledCollisions map[utils.EntityHit]bool
}

func (p Player)GetId() string {
    return "player"
}

func (p Player)GetCurrentFrame() int {
    return 0
}

func (p Player)GetPosition() utils.Vec2[float64] {
    return p.Position
}

func (p Player)GetSpriteSize() utils.Vec2[float64] {
    return p.spriteSize
}

func (p Player)GetHitbox() utils.Vec2[float64] {
    return p.hitbox
}

func (p *Player)Init(bounds utils.Vec2[int]) {
    p.lives = 3

    p.spriteSize = utils.CreateVec(11.0, 8.0)
    p.hitbox = p.spriteSize
    p.Position = utils.CreateVec(float64(bounds.X)/2.0, float64(bounds.Y) - p.spriteSize.Y)

    p.hitboxReceiveMask = hitbox.ENEMY
    p.hitboxSendMask = hitbox.PLAYER
}

func (p Player)GetEntityType() int {
    return ids.PLAYER
}

func (p Player)GetGamestateIx() int {
    return 0
}

func (p Player)GetHitboxSendMask() uint8 {
    return p.hitboxSendMask
}

func (p Player)GetHiboxReceiveMask() uint8 {
    return p.hitboxReceiveMask
}

func (p *Player)Hit() {
    p.lives -= 1
}

func (p Player)IsCollisionHandled(ent utils.EntityHit) bool {
    _, exists := p.HandledCollisions[ent]
    return exists
}

func (p Player)GetShotOnCooldown() bool {
    return p.ShotOnCooldown
}

func (p Player)GetLives() int {
    return p.lives
}

