package main

const WALL_SIZE_X = 24
const WALL_SIZE_Y = 24

type Wall struct {
    position Vec2[float64]
    body WallBody

    hitbox Vec2[float64]
    hitboxSendMask uint8
    hitboxReceiveMask uint8
}

func (w Wall)getPosition() Vec2[float64] {
    return w.position
}

func (w Wall)getHitbox() Vec2[float64] {
    return w.hitbox
}

func (w Wall)getHitboxSendMask() uint8 {
    return w.hitboxSendMask
}

func (w Wall)getHiboxReceiveMask() uint8 {
    return w.hitboxReceiveMask
}

func (w Wall)getEntityType() int {
    return ENTITY_WALL
}

func (w Wall)getGamestateIx() int {
    return 0
}

func (w Wall)didCollideWith(e EntityHit) bool {
    return false//TODO
}
