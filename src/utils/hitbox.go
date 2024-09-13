package utils

type EntityHit interface {
    GetPosition() Vec2[float64]
    GetHitbox() Vec2[float64]
    GetHitboxSendMask() uint8
    GetHiboxReceiveMask() uint8
    GetEntityType() int
    GetGamestateIx() int
    IsCollisionHandled(e EntityHit) bool
}

func GetHitboxBounds(e EntityHit) (Vec2[float64], Vec2[float64]) {
    ePos := e.GetPosition()
    eHitbox := e.GetHitbox().scale(0.5)

    return ePos.Subtract(eHitbox), ePos.Add(eHitbox)
}

func HitboxCollide(e1 EntityHit, e2 EntityHit) bool {
    e1Min, e1Max := GetHitboxBounds(e1)
    e2Min, e2Max := GetHitboxBounds(e2)

    return RectIntersects(e1Min, e1Max, e2Min, e2Max)
} 

func HitboxReceive(sender EntityHit, receiver EntityHit) bool {
    sMask := sender.GetHitboxSendMask()
    rMask := receiver.GetHiboxReceiveMask()

    return sMask & rMask > 0
}

