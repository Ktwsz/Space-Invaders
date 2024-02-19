package main

const WALL_SIZE_X = 24
const WALL_SIZE_Y = 24
const WALL_HIT_RADIUS = 2.0

type WallBody = [WALL_SIZE_X][WALL_SIZE_Y]bool

type Wall struct {
    position Vec2[float64]
    body WallBody

    hitbox Vec2[float64]
    hitboxSendMask uint8
    hitboxReceiveMask uint8

    gamestateIx int

    handledCollisions map[EntityHit]bool
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
    return w.gamestateIx
}

func IsInWallRect(w Wall, pos Vec2[float64]) bool {
    wMin, wMax := getHitboxBounds(w)

    return pos.x >= wMin.x &&
           pos.x <= wMax.x &&
           pos.y >= wMin.y && 
           pos.y <= wMax.y
} 

func (w Wall)IsCollisionHandled(e EntityHit) bool {
    _, exists := w.handledCollisions[e]
    return exists
}

func (w Wall)getHitPos(e EntityHit) (Vec2[int], bool) {
    eMin, eMax := getHitboxBounds(e)
    wMin, _ := getHitboxBounds(w)

    for x := int(eMin.x); x <= int(eMax.x); x++ {
        for y := int(eMin.y); y <= int(eMax.y); y++ {
            pos := Vec2[int]{x: x, y: y}.toFloat64()
            
            if IsInWallRect(w, pos) {
                wallHitboxPos := pos.subtract(wMin).toInt()
                if wallHitboxPos.x >= 0 &&
                   wallHitboxPos.x < WALL_SIZE_X &&
                   wallHitboxPos.y >= 0 &&
                   wallHitboxPos.y < WALL_SIZE_Y &&
                   w.body[wallHitboxPos.x][wallHitboxPos.y] {
                    return wallHitboxPos, true
                }
            }
        }
    }

    return Vec2[int]{}, false
}

func (w *Wall)Hit(pos Vec2[int]) {
    for x := range WALL_SIZE_X {
        for y := range WALL_SIZE_Y {
            pos := pos.subtract(Vec2[int]{x: x, y: y}).toFloat64()

            if vecLen(pos) <= WALL_HIT_RADIUS {
                w.body[x][y] = false
            }
        }
    } 
}
