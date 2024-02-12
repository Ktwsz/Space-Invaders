package main

type Enemy struct {
    id string
    frame int
    frameCount int

    rowData Vec2[int]
    position Vec2[float64]
    hitbox Vec2[float64]
    spriteSize Vec2[float64]

    deathState int
    gamestateIx int
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

func (e *Enemy)Init(id string, frameCount int, position Vec2[float64], rowData Vec2[int]) {
    e.id = id
    e.frameCount = frameCount

    e.rowData = rowData
    e.position = position
    e.spriteSize = Vec2[float64]{x: 8.0, y: 8.0}
    e.hitbox = Vec2[float64]{x: 8.0, y: 8.0}
}

func (e *Enemy)Move(speed float64) {
    if e.deathState != STATE_ALIVE {
        return
    }
    e.position.x += speed
    e.frame = (e.frame + 1) % e.frameCount
}

func (e *Enemy)Shift(shift Vec2[float64]) {
    if e.deathState != STATE_ALIVE {
        return
    }
    e.position = e.position.add(shift)
}

func (e Enemy)getHitbox() Vec2[float64] {
    return e.hitbox
}

func (e *Enemy)StartDying() {
    e.deathState = STATE_DEATH_START
    e.frame = e.frameCount
}

func (e Enemy)getEntityType() int {
    return ENTITY_ENEMY    
}
    
func (e Enemy)getGamestateIx() int {
    return e.gamestateIx
}
