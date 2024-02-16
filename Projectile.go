package main

type Projectile struct {
    id string
    frame int
    frameCount int

    position Vec2[float64]

    hitbox Vec2[float64]
    hitboxSendMask uint8
    hitboxReceiveMask uint8

    spriteSize Vec2[float64]
    speed float64

    deathState int
    gamestateIx int
}

func (p Projectile)getId() string {
    return p.id
}

func (p Projectile)getCurrentFrame() int {
    return p.frame
}

func (p Projectile)getPosition() Vec2[float64] {
    return p.position
}

func (p Projectile)getSpriteSize() Vec2[float64] {
    return p.spriteSize
}

func (p *Projectile)Move() {
    if p.deathState != STATE_ALIVE {
        return 
    }
    p.position.y += p.speed
    p.frame = (p.frame + 1) % p.frameCount
}

func (p Projectile)getHitbox() Vec2[float64] {
    return p.hitbox
}

func (p *Projectile)StartDying() {
    p.deathState = STATE_DEATH_START
    p.frame = p.frameCount
}

func (p Projectile)getEntityType() int {
    return ENTITY_PROJECTILE
}
    
func (p Projectile)getGamestateIx() int {
    return p.gamestateIx
}

func (p Projectile)getHitboxSendMask() uint8 {
    return p.hitboxSendMask
}

func (p Projectile)getHiboxReceiveMask() uint8 {
    return p.hitboxReceiveMask
}
