package entity

import (
    "space_invaders/utils"
    "space_invaders/entity/states"
    "space_invaders/entity/hitbox"
    "space_invaders/entity/ids"
)

type Enemy struct {
    id string
    projectileId string

    frame int
    frameCount int

    rowData utils.Vec2[int]
    position utils.Vec2[float64]

    hitbox utils.Vec2[float64]
    hitboxSendMask uint8
    hitboxReceiveMask uint8

    spriteSize utils.Vec2[float64]

    DeathState int
    GamestateIx int

    points int

    HandledCollisions map[utils.EntityHit]bool
}

func (e Enemy)GetId() string {
    return e.id
}

func (e Enemy)GetCurrentFrame() int {
    return e.frame
}

func (e Enemy)GetPosition() utils.Vec2[float64] {
    return e.position
}

func (e Enemy)GetSpriteSize() utils.Vec2[float64] {
    return e.spriteSize
}

func (e *Enemy)Init(id string, projectileId string, frameCount int, position utils.Vec2[float64], rowData utils.Vec2[int], points int) {
    e.id = id
    e.projectileId = projectileId

    e.frameCount = frameCount

    e.hitboxReceiveMask = hitbox.PLAYER
    e.hitboxSendMask = hitbox.ENEMY

    e.rowData = rowData
    e.position = position
    e.spriteSize = utils.CreateVec(8.0, 8.0)
    e.hitbox = utils.CreateVec(8.0, 8.0)

    e.points = points
}

func (e *Enemy)Move(speed float64) {
    if e.DeathState != states.ALIVE {
        return
    }
    e.position.X += speed
    e.frame = (e.frame + 1) % e.frameCount
}

func (e *Enemy)Shift(shift utils.Vec2[float64]) {
    if e.DeathState != states.ALIVE {
        return
    }
    e.position = e.position.Add(shift)
}

func (e Enemy)GetHitbox() utils.Vec2[float64] {
    return e.hitbox
}

func (e *Enemy)StartDying() {
    e.DeathState = states.DEATH_START
    e.frame = e.frameCount
}

func (e Enemy)GetEntityType() int {
    return ids.ENEMY    
}
    
func (e Enemy)GetGamestateIx() int {
    return e.GamestateIx
}

func (e Enemy)GetHitboxSendMask() uint8 {
    return e.hitboxSendMask
}

func (e Enemy)GetHiboxReceiveMask() uint8 {
    return e.hitboxReceiveMask
}

func (e Enemy)IsCollisionHandled(ent utils.EntityHit) bool {
    _, exists := e.HandledCollisions[ent]
    return exists
}

func(e Enemy)GetRowData() utils.Vec2[int] {
    return e.rowData
}

func (e Enemy)GetProjectileId() string {
    return e.projectileId
}

func (e Enemy)GetPoints() int {
    return e.points
}

