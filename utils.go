package main

import (
    "golang.org/x/exp/constraints"
    "sort"
    "math"
)

const (
    ENTITY_PLAYER = iota
    ENTITY_ENEMY
    ENTITY_PROJECTILE
    ENTITY_WALL
)

const (
    STATE_ALIVE = iota
    STATE_DEATH_START
    STATE_DEATH_END
)

const (
    HITBOX_PLAYER = 1 << iota
    HITBOX_ENEMY
    HITBOX_PROJECTILE
    HITBOX_WALL
)

const (
    GAME_STARTING = iota
    GAME_RUNNING
    GAME_OVER
    GAME_WIN
)

type Vec2[T constraints.Float | constraints.Integer] struct {
    x, y T
}

func (v Vec2[T])add(other Vec2[T]) Vec2[T] {
    return Vec2[T]{x: v.x + other.x, y: v.y + other.y}
}

func (v Vec2[T])subtract(other Vec2[T]) Vec2[T] {
    return Vec2[T]{x: v.x - other.x, y: v.y - other.y}
}

func (v Vec2[T])scale(scalar T) Vec2[T] {
    return Vec2[T]{x: v.x * scalar, y: v.y * scalar}
}

func (v Vec2[float64])toInt() Vec2[int] {
    return Vec2[int]{x: int(v.x), y: int(v.y)}
}

func (v Vec2[int])toFloat64() Vec2[float64] {
    return Vec2[float64]{x: float64(v.x), y: float64(v.y)}
}

func vecLen(v Vec2[float64]) float64 {
    return math.Sqrt(v.x * v.x + v.y * v.y)
}


func RectIntersects[T constraints.Float | constraints.Integer](Min1 Vec2[T], Max1 Vec2[T], Min2 Vec2[T], Max2 Vec2[T]) bool {
    return Min1.x < Max2.x &&
           Max1.x > Min2.x &&
           Min1.y < Max2.y &&
           Max1.y > Min2.y
}

type EntityDraw interface {
    getId() string
    getCurrentFrame() int
    getPosition() Vec2[float64]
    getSpriteSize() Vec2[float64]
    getEntityType() int
}

type EntityHit interface {
    getPosition() Vec2[float64]
    getHitbox() Vec2[float64]
    getHitboxSendMask() uint8
    getHiboxReceiveMask() uint8
    getEntityType() int
    getGamestateIx() int
    didCollideWith(e EntityHit) bool
}

func getHitboxBounds(e EntityHit) (Vec2[float64], Vec2[float64]) {
    ePos := e.getPosition()
    eHitbox := e.getHitbox().scale(0.5)

    return ePos.subtract(eHitbox), ePos.add(eHitbox)
}

func HitboxCollide(e1 EntityHit, e2 EntityHit) bool {
    e1Min, e1Max := getHitboxBounds(e1)
    e2Min, e2Max := getHitboxBounds(e2)

    return RectIntersects(e1Min, e1Max, e2Min, e2Max)
} 

func HitboxReceive(sender EntityHit, receiver EntityHit) bool {
    sMask := sender.getHitboxSendMask()
    rMask := receiver.getHiboxReceiveMask()

    return sMask & rMask > 0
}

func getDrawPosition(e EntityDraw) Vec2[float64] {
    pos := e.getPosition()
    sprite := e.getSpriteSize()
    return pos.subtract(sprite.scale(0.5))
}

func IsOutOfBounds(bounds Vec2[int], e EntityDraw) bool {
    pos := e.getPosition()
    sprite := e.getSpriteSize()
    minEdge := pos.subtract(sprite.scale(0.5))
    maxEdge := pos.add(sprite.scale(0.5))

    return minEdge.x < 0 || minEdge.y < 0 || maxEdge.x > float64(bounds.x) || maxEdge.y > float64(bounds.y)
}

func RemoveIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}

func RemoveIndexesMany[T any](s []T, indexes []int) []T {
    if indexes == nil {
        return s
    }

    sort.Ints(indexes)

    for count, ix := range indexes {
        s = RemoveIndex(s, ix - count)
    }
    return s
}
