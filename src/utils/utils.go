package utils

import (
    "golang.org/x/exp/constraints"
    "sort"
    "math"
)

type Vec2[T constraints.Float | constraints.Integer] struct {
    X, Y T
}

func CreateVec[T constraints.Float | constraints.Integer](x, y T) Vec2[T] {
    return Vec2[T] {x, y}
}

func (v Vec2[T])Add(other Vec2[T]) Vec2[T] {
    return Vec2[T]{X: v.X + other.X, Y: v.Y + other.Y}
}

func (v Vec2[T])Subtract(other Vec2[T]) Vec2[T] {
    return Vec2[T]{X: v.X - other.X, Y: v.Y - other.Y}
}

func (v Vec2[T])scale(scalar T) Vec2[T] {
    return Vec2[T]{X: v.X * scalar, Y: v.Y * scalar}
}

func (v Vec2[float64])ToInt() Vec2[int] {
    return Vec2[int]{X: int(v.X), Y: int(v.Y)}
}

func (v Vec2[int])ToFloat64() Vec2[float64] {
    return Vec2[float64]{X: float64(v.X), Y: float64(v.Y)}
}

func VecLen(v Vec2[float64]) float64 {
    return math.Sqrt(v.X * v.X + v.Y * v.Y)
}


func RectIntersects[T constraints.Float | constraints.Integer](Min1 Vec2[T], Max1 Vec2[T], Min2 Vec2[T], Max2 Vec2[T]) bool {
    return Min1.X < Max2.X &&
           Max1.X > Min2.X &&
           Min1.Y < Max2.Y &&
           Max1.Y > Min2.Y
}

type EntityDraw interface {
    GetId() string
    GetCurrentFrame() int
    GetPosition() Vec2[float64]
    GetSpriteSize() Vec2[float64]
    GetEntityType() int
}


func GetDrawPosition(e EntityDraw) Vec2[float64] {
    pos := e.GetPosition()
    sprite := e.GetSpriteSize()
    return pos.Subtract(sprite.scale(0.5))
}

func IsOutOfBounds(bounds Vec2[int], e EntityDraw) bool {
    pos := e.GetPosition()
    sprite := e.GetSpriteSize()
    minEdge := pos.Subtract(sprite.scale(0.5))
    maxEdge := pos.Add(sprite.scale(0.5))

    return minEdge.X < 0 || minEdge.Y < 0 || maxEdge.X > float64(bounds.X) || maxEdge.Y > float64(bounds.Y)
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
