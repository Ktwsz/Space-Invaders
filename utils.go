package main

import "golang.org/x/exp/constraints"

type Vec2[T constraints.Float | constraints.Integer] struct {
    x, y T
}

type Entity interface {
    getId() string
    getCurrentFrame() int
    getPosition() Vec2[float64]
    getDrawPosition() Vec2[float64]
}
