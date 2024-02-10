package main

import (
    "golang.org/x/exp/constraints"
)

type Vec2[T constraints.Float | constraints.Integer] struct {
    x, y T
}

func (v Vec2[T])add(other Vec2[T]) Vec2[T] {
    return Vec2[T]{x: v.x + other.x, y: v.y + other.y}
}

type Entity interface {
    getId() string
    getCurrentFrame() int
    getPosition() Vec2[float64]
    getSpriteSize() Vec2[float64]
}

func getDrawPosition(e Entity) Vec2[float64] {
    pos := e.getPosition()
    sprite := e.getSpriteSize()
    return Vec2[float64]{x: pos.x - sprite.x/2.0, y: pos.y - sprite.y/2.0}
}

func IsOutOfBounds(bounds Vec2[int], e Entity) bool {
    minEdge := Vec2[float64]{x: e.getPosition().x - e.getSpriteSize().x/2.0, y: e.getPosition().y - e.getSpriteSize().y/2.0}
    maxEdge := Vec2[float64]{x: e.getPosition().x + e.getSpriteSize().x/2.0, y: e.getPosition().y + e.getSpriteSize().y/2.0}

    return minEdge.x < 0 || minEdge.y < 0 || maxEdge.x > float64(bounds.x) || maxEdge.y > float64(bounds.y)
}

func RemoveIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}

func RemoveIndexesMany[T any](s []T, indexes []int) []T {
    for count, ix := range indexes {
        s = RemoveIndex(s, ix - count)
    }
    return s
}
