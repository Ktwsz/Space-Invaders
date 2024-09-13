package utils

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
    boundsF := Vec2[int]{X: bounds.X, Y: bounds.Y}.ToFloat64()
    T.center = boundsF.scale(0.5)
    T.radius = max(boundsF.X, boundsF.Y)

    return T
}

func (T *QTree)divide() {
    T.isDivided = true

    newRadius := T.radius/2
    x, y := T.center.X, T.center.Y

    T.children[topleft] = &QTree{center: Vec2[float64]{X: x - newRadius, Y: y - newRadius}, radius: newRadius}
    T.children[topright] = &QTree{center: Vec2[float64]{X: x + newRadius, Y: y - newRadius}, radius: newRadius}
    T.children[bottomleft] = &QTree{center: Vec2[float64]{X: x - newRadius, Y: y + newRadius}, radius: newRadius}
    T.children[bottomright] = &QTree{center: Vec2[float64]{X: x + newRadius, Y: y + newRadius}, radius: newRadius}

    for _, e := range T.entities {
        for i := range childrenCount {
            T.children[i].Insert(e)
        }
    }
}

func (T *QTree)Insert(entity EntityHit) {
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
        T.children[i].Insert(entity)
    }
}

func (T *QTree)GetAllIntersections() []*QTree {
    if !T.isDivided {
        if countEntites(T) > 1 {
            return []*QTree{T}
        } else {
            return []*QTree{}
        }
    }

    result := make([]*QTree, 0)
    for i := range childrenCount {
        if childrenResult := T.children[i].GetAllIntersections(); len(childrenResult) > 0 {
            result = append(result, childrenResult...)
        }
    }

    return result
}

func (T *QTree)GetEntities() [maxEntities]EntityHit {
    return T.entities
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
    ePos := entity.GetPosition()
    eHitbox := entity.GetHitbox().scale(0.5)
    eMin, eMax := ePos.Subtract(eHitbox), ePos.Add(eHitbox)

    TPos := T.center
    TRad := Vec2[float64]{X: T.radius, Y: T.radius}
    TMin, TMax := TPos.Subtract(TRad), TPos.Add(TRad)

    return RectIntersects(eMin, eMax, TMin, TMax)
}
