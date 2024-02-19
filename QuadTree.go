package main

const childrenCount = 4
const (
    topleft = iota
    topright
    bottomleft
    bottomright
)
const maxEntities = 2

type QTree struct {
    center Vec2[float64]
    radius float64
    children [childrenCount]*QTree
    entities [maxEntities]EntityHit
    isDivided bool
}

func QTreeInitFromBounds(bounds Vec2[int]) QTree {
    T := QTree{}
    boundsF := Vec2[float64]{x: float64(bounds.x), y: float64(bounds.y)}
    T.center = boundsF.scale(0.5)
    T.radius = max(boundsF.x, boundsF.y)

    return T
}

func (T *QTree)divide() {
    T.isDivided = true

    newRadius := T.radius/2
    x, y := T.center.x, T.center.y

    T.children[topleft] = &QTree{center: Vec2[float64]{x: x - newRadius, y: y - newRadius}, radius: newRadius}
    T.children[topright] = &QTree{center: Vec2[float64]{x: x + newRadius, y: y - newRadius}, radius: newRadius}
    T.children[bottomleft] = &QTree{center: Vec2[float64]{x: x - newRadius, y: y + newRadius}, radius: newRadius}
    T.children[bottomright] = &QTree{center: Vec2[float64]{x: x + newRadius, y: y + newRadius}, radius: newRadius}

    for _, e := range T.entities {
        for i := range childrenCount {
            T.children[i].insert(e)
        }
    }
}

func (T *QTree)insert(entity EntityHit) {
    if !T.IsInBounds(entity) {
        return 
    }

    if eCount := countEntites(T); !T.isDivided && eCount < maxEntities {
        T.entities[eCount] = entity 
        return 
    }

    if !T.isDivided {
        T.divide()
    }

    for i := range childrenCount {
        T.children[i].insert(entity)
    }
}

func (T *QTree)getAllIntersections() []*QTree {
    if !T.isDivided {
        if countEntites(T) > 1 {
            return []*QTree{T}
        } else {
            return []*QTree{}
        }
    }

    result := make([]*QTree, 0)
    for i := range childrenCount {
        if childrenResult := T.children[i].getAllIntersections(); len(childrenResult) > 0 {
            result = append(result, childrenResult...)
        }
    }

    return result
}

func countEntites(T *QTree) int {
    for i := range maxEntities {
        if T.entities[i] == nil {
            return i
        }
    }    

    return maxEntities
}

func (T *QTree)IsInBounds(entity EntityHit) bool {
    ePos := entity.getPosition()
    eHitbox := entity.getHitbox().scale(0.5)
    eMin, eMax := ePos.subtract(eHitbox), ePos.add(eHitbox)

    TPos := T.center
    TRad := Vec2[float64]{x: T.radius, y: T.radius}
    TMin, TMax := TPos.subtract(TRad), TPos.add(TRad)

    return RectIntersects[float64](eMin, eMax, TMin, TMax)
}
