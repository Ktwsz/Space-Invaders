package main

import "golang.org/x/exp/constraints"

type Vec2[T constraints.Float | constraints.Integer] struct {
    x, y T
}
