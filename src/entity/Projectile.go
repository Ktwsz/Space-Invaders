package entity

import (
	"space_invaders/entity/hitbox"
	"space_invaders/entity/ids"
	"space_invaders/entity/states"
	"space_invaders/utils"
)

type Projectile struct {
    id string
    frame int
    frameCount int

    position utils.Vec2[float64]

    hitbox utils.Vec2[float64]
    hitboxSendMask uint8
    hitboxReceiveMask uint8

    spriteSize utils.Vec2[float64]
    speed float64

    DeathState int
    GamestateIx int

    HandledCollisions map[utils.EntityHit]bool
}

func CreateProjectile(id string, frameCount int, position utils.Vec2[float64], size utils.Vec2[float64], speed float64, hitboxReceive, hitboxSend uint8) Projectile {
    return Projectile {
        id: id,
        frameCount: frameCount,
        position: position,
        speed: speed,
        hitbox: size,
        spriteSize: size,
        hitboxReceiveMask: hitbox.PROJECTILE | hitbox.WALL | hitboxReceive,
        hitboxSendMask: hitbox.PROJECTILE | hitboxSend,
    }
}


func (p Projectile)GetId() string {
    return p.id
}

func (p Projectile)GetCurrentFrame() int {
    return p.frame
}

func (p Projectile)GetPosition() utils.Vec2[float64] {
    return p.position
}

func (p Projectile)GetSpriteSize() utils.Vec2[float64] {
    return p.spriteSize
}

func (p *Projectile)Move() {
    if p.DeathState != states.ALIVE {
        return 
    }
    p.position.Y += p.speed
    p.frame = (p.frame + 1) % p.frameCount
}

func (p Projectile)GetHitbox() utils.Vec2[float64] {
    return p.hitbox
}

func (p *Projectile)StartDying() {
    p.DeathState = states.DEATH_START
    p.frame = p.frameCount
}

func (p Projectile)GetEntityType() int {
    return ids.PROJECTILE
}
    
func (p Projectile)GetGamestateIx() int {
    return p.GamestateIx
}

func (p Projectile)GetHitboxSendMask() uint8 {
    return p.hitboxSendMask
}

func (p Projectile)GetHiboxReceiveMask() uint8 {
    return p.hitboxReceiveMask
}

func (p Projectile)IsCollisionHandled(ent utils.EntityHit) bool {
    _, exists := p.HandledCollisions[ent]
    return exists
}
