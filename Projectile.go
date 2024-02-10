package main

type Projectile struct {
    id string
    frame int
    frameCount int

    position Vec2[float64]
    hitbox Vec2[float64]
    spriteSize Vec2[float64]
    speed float64
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
    p.position.y += p.speed
    p.frame = (p.frame + 1) % p.frameCount
}
