package entity

import (
    "space_invaders/utils"
    "space_invaders/entity/ids"
    "space_invaders/entity/hitbox"
)

const WALL_SIZE_X = 24
const WALL_SIZE_Y = 24
const WALL_HIT_RADIUS = 2.0

type WallBody = [WALL_SIZE_X][WALL_SIZE_Y]bool

type Wall struct {
    position utils.Vec2[float64]
    Body WallBody

    hitbox utils.Vec2[float64]
    hitboxSendMask uint8
    hitboxReceiveMask uint8

    gamestateIx int

    HandledCollisions map[utils.EntityHit]bool
}

func CreateWall(wallPos utils.Vec2[float64], i int) *Wall {
    return &Wall{position: wallPos,
       hitbox: utils.CreateVec(WALL_SIZE_X, WALL_SIZE_Y).ToFloat64(),
       hitboxReceiveMask: hitbox.PROJECTILE,
       hitboxSendMask: hitbox.WALL,
       gamestateIx: i,
   }
}

func (w Wall)GetPosition() utils.Vec2[float64] {
    return w.position
}

func (w Wall)GetHitbox() utils.Vec2[float64] {
    return w.hitbox
}

func (w Wall)GetHitboxSendMask() uint8 {
    return w.hitboxSendMask
}

func (w Wall)GetHiboxReceiveMask() uint8 {
    return w.hitboxReceiveMask
}

func (w Wall)GetEntityType() int {
    return ids.WALL
}

func (w Wall)GetGamestateIx() int {
    return w.gamestateIx
}

func IsInWallRect(w Wall, pos utils.Vec2[float64]) bool {
    wMin, wMax := utils.GetHitboxBounds(w)

    return pos.X >= wMin.X &&
           pos.X <= wMax.X &&
           pos.Y >= wMin.Y && 
           pos.Y <= wMax.Y
} 

func (w Wall)IsCollisionHandled(e utils.EntityHit) bool {
    _, exists := w.HandledCollisions[e]
    return exists
}

func (w Wall)GetHitPos(e utils.EntityHit) (utils.Vec2[int], bool) {
    eMin, eMax := utils.GetHitboxBounds(e)
    wMin, _ := utils.GetHitboxBounds(w)

    for x := int(eMin.X); x <= int(eMax.X); x++ {
        for y := int(eMin.Y); y <= int(eMax.Y); y++ {
            pos := utils.CreateVec(x, y).ToFloat64()
            
            if IsInWallRect(w, pos) {
                wallHitboxPos := pos.Subtract(wMin).ToInt()
                if wallHitboxPos.X >= 0 &&
                   wallHitboxPos.X < WALL_SIZE_X &&
                   wallHitboxPos.Y >= 0 &&
                   wallHitboxPos.Y < WALL_SIZE_Y &&
                   w.Body[wallHitboxPos.X][wallHitboxPos.Y] {
                    return wallHitboxPos, true
                }
            }
        }
    }

    return utils.Vec2[int]{}, false
}

func (w *Wall)Hit(pos utils.Vec2[int]) {
    for x := range WALL_SIZE_X {
        for y := range WALL_SIZE_Y {
            pos := pos.Subtract(utils.CreateVec(x, y)).ToFloat64()

            if utils.VecLen(pos) <= WALL_HIT_RADIUS {
                w.Body[x][y] = false
            }
        }
    } 
}

func (w Wall)GetBody(i, j int) bool {
    return w.Body[i][j]
}

