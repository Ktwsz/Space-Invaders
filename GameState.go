package main

import (
    "time"
    "fmt"
)

const playerSpeed float64 = 1.0
const enemyHeight float64 = 8.0
const enemyCountPerRow = 11

type GameState struct {
    bounds Vec2[int]

    player Player 

    enemies []Enemy
    enemySpeed float64

    projectiles []Projectile
}

func (g *GameState)Init() {
    g.bounds = Vec2[int]{x: 130, y: 120}
    g.player.Init(g.bounds)

    g.enemySpeed = 2.0

    g.SpawnEnemiesRow(0, "enemy1", enemyCountPerRow)
    g.SpawnEnemiesRow(1, "enemy2", enemyCountPerRow)
    g.SpawnEnemiesRow(2, "enemy3", enemyCountPerRow)

    enemyMoveTicker := time.NewTicker(500 * time.Millisecond)
    tickerDone := make(chan bool)
    go func() {
        for {
            select{
                case <-tickerDone:
                    return
                case <-enemyMoveTicker.C:
                    g.MoveEnemies()
            }
        }
    }()
}

func (g *GameState)GetObjectsToDraw() []Entity{
    enemiesLen := len(g.enemies)
    projectilesLen := len(g.projectiles)

    objects := make([]Entity, projectilesLen + enemiesLen + 1)
    objects[0] = &g.player

    for i := range g.enemies {
        objects[i + 1] = &g.enemies[i]
    }

    for i := range g.projectiles {
        objects[enemiesLen+i+1] = &g.projectiles[i]
    }

    return objects
}

func (g *GameState)PlayerMoveLeft() {
    g.player.position.x -= playerSpeed
    if g.player.position.x - g.player.spriteSize.x/2.0 < 0 {
        g.player.position.x = g.player.spriteSize.x/2.0
    }
}

func (g *GameState)PlayerMoveRight() {
    g.player.position.x += playerSpeed
    if g.player.position.x + g.player.spriteSize.x/2.0 > float64(g.bounds.x) {
        g.player.position.x = float64(g.bounds.x) - g.player.spriteSize.x/2.0
    }
}

func (g *GameState)SpawnEnemiesRow(row int, enemyId string, count int) {
    position := Vec2[float64]{x: 6.0, y: (enemyHeight + 6.0) * float64(row) + 6.0}

    enemiesNew := make([]Enemy, count)
    for i  := range count {
        enemiesNew[i] = Enemy{}
        enemiesNew[i].Init(enemyId, 3, position, row)

        position.x += 10.0//enemy x + margin
    }

    g.enemies = append(g.enemies, enemiesNew...)
}

func (g *GameState)PlayerShoot() {
    if ok := g.SpawnPlayerProjectile(); ok {
        g.player.shotOnCooldown = true
        go func() {
            timer := time.NewTimer(750 * time.Millisecond)

            <-timer.C
            g.player.shotOnCooldown = false
        }()
    }
}

func (g *GameState)SpawnPlayerProjectile() bool {
    if g.player.shotOnCooldown {
        return false
    }
    projectile := Projectile{id: "player_projectile", 
                             frameCount: 1,
                             position: Vec2[float64]{x: g.player.position.x, y: g.player.position.y - g.player.spriteSize.y/2.0},
                             hitbox: Vec2[float64]{x: 1, y: 6},
                             spriteSize: Vec2[float64]{x: 1, y: 6},
                             speed: -3.5}


    g.projectiles = append(g.projectiles, projectile)
    return true
}

func (g *GameState)MoveEnemies() {
    for i := range g.enemies {
        g.enemies[i].Move(g.enemySpeed)
    }
}

func (g *GameState)MoveProjectiles() {
    for i := range g.projectiles {
        g.projectiles[i].Move()
    }
}


func (g *GameState)RemoveMissedProjectiles() {
    toRemove := make([]int, 0)
    for i, p := range g.projectiles {
        if IsOutOfBounds(g.bounds, p) {
            toRemove = append(toRemove, i)
        }
    }

    g.projectiles = RemoveIndexesMany(g.projectiles, toRemove)
}

func (g *GameState)EnemiesShiftRow(row int, enemyIx int) {
    var shiftX float64

    e := g.enemies[enemyIx]
    if e.position.x + e.spriteSize.x > float64(g.bounds.x) {
        shiftX = -(e.position.x + e.spriteSize.x - float64(g.bounds.x))
    } else {
        shiftX = e.spriteSize.x - e.position.x
    }

    for i := range g.enemies {
        if g.enemies[i].rowNumber == row {
            g.enemies[i].Shift(Vec2[float64]{x: shiftX, y: 0})
        }
    }
}

func (g *GameState)EnemiesShiftDown() {
    for i := range g.enemies {
        g.enemies[i].Shift(Vec2[float64]{x: 0, y: enemyHeight + 6.0})        
    }
}

func (g *GameState)HandleCollisions() {
    g.RemoveMissedProjectiles()

    changeEnemiesDirection := false
    for i, e := range g.enemies {
        if IsOutOfBounds(g.bounds, e) {
            changeEnemiesDirection = true
            g.EnemiesShiftRow(e.rowNumber, i)
        }
    }

    if changeEnemiesDirection {
        g.enemySpeed *= -1
        g.EnemiesShiftDown()
    }

    tree := QTreeInitFromBounds(g.bounds)
    tree.insert(&g.player)
    for i := range g.enemies {
        tree.insert(&g.enemies[i])
    }
    for i := range g.projectiles {
        tree.insert(&g.projectiles[i])
    }

    collisions := tree.getAllIntersections()

    for i := range collisions {
        e := collisions[i].entities 
        e1, e2 := e[0], e[1]

        if ok := EntitiesCollide(e1, e2); ok {
            fmt.Printf("collision between i%+v and %+v\n", e1, e2)
        }
    }
}
