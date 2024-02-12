package main

import (
    "time"
    "sync"
    gauss "github.com/chobie/go-gaussian"
    "math"
    rand "math/rand"
)

const playerSpeed float64 = 1.0
const enemyHeight float64 = 8.0
const enemyCountColumn = 11
const PlayerProjectileSpeed = -3.5
const enemyProjectileSpeed = 3.5

type GameState struct {
    bounds Vec2[int]

    player Player 

    enemies []*Enemy
    enemyColProjectileCooldown [enemyCountColumn]bool
    deadEnemies []int
    enemySpeed float64

    projectiles []*Projectile

    mutex sync.Mutex
}

func (g *GameState)Init() {
    g.bounds = Vec2[int]{x: 130, y: 120}
    g.player.Init(g.bounds)

    g.enemySpeed = 2.0

    g.SpawnEnemiesRow(0, "enemy1", enemyCountColumn)
    g.SpawnEnemiesRow(1, "enemy2", enemyCountColumn)
    g.SpawnEnemiesRow(2, "enemy3", enemyCountColumn)

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

func (g *GameState)GetLastEnemyInCol(column int) *Enemy {
    g.mutex.Lock()
    var result *Enemy = nil
    for i, e := range g.enemies {
        if e.rowData.x == column && (result == nil || e.rowData.y > result.rowData.y) {
            result = g.enemies[i]
        }
    }
    g.mutex.Unlock()
    return result
}

func (g *GameState)GetChanceToShoot(gaussian *gauss.Gaussian, enemy *Enemy) bool {
    dist := math.Abs(enemy.position.x - g.player.position.x)
    chance := gaussian.Pdf(dist) * math.Sqrt(2 * math.Pi)
    
    return rand.Float64() <= chance
}

func (g *GameState)EnemyShoot() {
    gaussDist := gauss.NewGaussian(0, 1)
    for col := range enemyCountColumn {
        if enemyLast := g.GetLastEnemyInCol(col); enemyLast != nil {
            if chance := g.GetChanceToShoot(gaussDist, enemyLast); chance {
                didShoot := g.SpawnEnemyprojectile(enemyLast, col)

                if didShoot {
                    g.mutex.Lock()
                    g.enemyColProjectileCooldown[col] = true
                    g.mutex.Unlock()

                    go g.SetEnemyCooldownTimer(col, 750)
                }
            }
        }
    }
}

func (g *GameState)SpawnEnemyprojectile(enemy *Enemy, col int) bool {
    g.mutex.Lock()
    shotOnCooldown := g.enemyColProjectileCooldown[col]
    g.mutex.Unlock()

    if shotOnCooldown {
        return false
    }
    projectile := Projectile{id: "enemy_projectile_1", 
                             frameCount: 4,
                             position: Vec2[float64]{x: enemy.position.x, y: enemy.position.y + enemy.spriteSize.y/2.0},
                             hitbox: Vec2[float64]{x: 3, y: 7},
                             spriteSize: Vec2[float64]{x: 3, y: 7},
                             speed: enemyProjectileSpeed}


    g.mutex.Lock()
    g.projectiles = append(g.projectiles, &projectile)
    g.mutex.Unlock()

    return true
}

func (g *GameState)GetObjectsToDraw() []Entity{
    enemiesLen := len(g.enemies)
    projectilesLen := len(g.projectiles)

    objects := make([]Entity, projectilesLen + enemiesLen + 1)

    g.mutex.Lock()
    objects[0] = &g.player

    for i := range g.enemies {
        objects[i + 1] = g.enemies[i]
    }

    for i := range g.projectiles {
        objects[enemiesLen+i+1] = g.projectiles[i]
    }
    g.mutex.Unlock()

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

    enemiesNew := make([]*Enemy, count)

    for i := range count {
        enemiesNew[i] = &Enemy{}
        enemiesNew[i].Init(enemyId, 3, position, Vec2[int]{x: i, y: row})

        position.x += 10.0//enemy x + margin
    }

    g.mutex.Lock()
    g.enemies = append(g.enemies, enemiesNew...)
    g.mutex.Unlock()
}

func (g *GameState)removeDeadEnemies() {
    toRemove := make([]int, 0)
    g.mutex.Lock()
    for i, e := range g.enemies {
        if e.deathState == STATE_DEATH_END {
            toRemove = append(toRemove, i)
        }
    }

    g.enemies = RemoveIndexesMany(g.enemies, toRemove)
    g.mutex.Unlock()
}

func (g *GameState)SetCooldownTimer(t time.Duration) {
    timer := time.NewTimer(t * time.Millisecond)

    <-timer.C

    g.mutex.Lock()
    g.player.shotOnCooldown = false
    g.mutex.Unlock()
}

func (g *GameState)SetEnemyCooldownTimer(column int, t time.Duration) {
    timer := time.NewTimer(t * time.Millisecond)

    <-timer.C

    g.mutex.Lock()
    g.enemyColProjectileCooldown[column] = false
    g.mutex.Unlock()
}

func (g *GameState)PlayerShoot() {
    if ok := g.SpawnPlayerProjectile(); ok {
        g.mutex.Lock()
        g.player.shotOnCooldown = true
        g.mutex.Unlock()

        go g.SetCooldownTimer(750)
    }
}

func (g *GameState)SpawnPlayerProjectile() bool {
    g.mutex.Lock()
    shotOnCooldown := g.player.shotOnCooldown
    g.mutex.Unlock()

    if shotOnCooldown {
        return false
    }
    projectile := Projectile{id: "player_projectile", 
                             frameCount: 1,
                             position: Vec2[float64]{x: g.player.position.x, y: g.player.position.y - g.player.spriteSize.y/2.0},
                             hitbox: Vec2[float64]{x: 1, y: 6},
                             spriteSize: Vec2[float64]{x: 1, y: 6},
                             speed: PlayerProjectileSpeed}


    g.mutex.Lock()
    g.projectiles = append(g.projectiles, &projectile)
    g.mutex.Unlock()

    return true
}

func (g *GameState)MoveEnemies() {
    g.mutex.Lock()
    for i := range g.enemies {
        g.enemies[i].Move(g.enemySpeed)
    }
    g.mutex.Unlock()
}

func (g *GameState)MoveProjectiles() {
    g.mutex.Lock()
    for i := range g.projectiles {
        g.projectiles[i].Move()
    }
    g.mutex.Unlock()
}


func (g *GameState)RemoveDeadProjectiles() {
    toRemove := make([]int, 0)
    g.mutex.Lock()
    for i, p := range g.projectiles {
        if p.deathState == STATE_DEATH_END {
            toRemove = append(toRemove, i)
        }
    }

    g.projectiles = RemoveIndexesMany(g.projectiles, toRemove)
    g.mutex.Unlock()
}

func (g *GameState)SetEnemyDeathTimer(enemy *Enemy, t time.Duration) {
    timer := time.NewTimer(t * time.Millisecond)

    <-timer.C

    g.mutex.Lock()
    enemy.deathState = STATE_DEATH_END
    g.mutex.Unlock()
}

func (g *GameState)SetProjectileDeathTimer(projectile *Projectile, t time.Duration) {
    timer := time.NewTimer(t * time.Millisecond)

    <-timer.C
    g.mutex.Lock()
    projectile.deathState = STATE_DEATH_END
    g.mutex.Unlock()
}

func (g *GameState)CheckForMissedProjectiles() {
    g.mutex.Lock()
    for i, p := range g.projectiles {
        if IsOutOfBounds(g.bounds, p) {
            g.projectiles[i].StartDying()

            go g.SetProjectileDeathTimer(g.projectiles[i], 500)
        }
    }
    g.mutex.Unlock()
}

func (g *GameState)EnemiesShiftRow(row int, enemyIx int) {
    var shiftX float64

    g.mutex.Lock()
    e := g.enemies[enemyIx]
    if e.position.x + e.spriteSize.x > float64(g.bounds.x) {
        shiftX = -(e.position.x + e.spriteSize.x - float64(g.bounds.x))
    } else {
        shiftX = e.spriteSize.x - e.position.x
    }

    for i := range g.enemies {
        if g.enemies[i].rowData.y == row {
            g.enemies[i].Shift(Vec2[float64]{x: shiftX, y: 0})
        }
    }
    g.mutex.Unlock()
}

func (g *GameState)EnemiesShiftDown() {
    g.mutex.Lock()
    for i := range g.enemies {
        g.enemies[i].Shift(Vec2[float64]{x: 0, y: enemyHeight + 6.0})        
    }
    g.mutex.Unlock()
}

func (g *GameState)CheckEnemiesInBounds() {
    changeEnemiesDirection := false
    for i, e := range g.enemies {
        g.mutex.Lock()
        if IsOutOfBounds(g.bounds, e) {
            changeEnemiesDirection = true
            g.mutex.Unlock()
            g.EnemiesShiftRow(e.rowData.y, i)
        } else {
            g.mutex.Unlock()
        }
    }

    if changeEnemiesDirection {
        g.enemySpeed *= -1
        g.EnemiesShiftDown()
    }
}

func (g *GameState)HandleCollisions() {
    g.mutex.Lock()

    tree := QTreeInitFromBounds(g.bounds)
    tree.insert(&g.player)

    for i := range g.enemies {
        if g.enemies[i].deathState == STATE_ALIVE {
            g.enemies[i].gamestateIx = i
            tree.insert(g.enemies[i])
        }
    }

    for i := range g.projectiles {
        if g.projectiles[i].deathState == STATE_ALIVE {
            g.projectiles[i].gamestateIx = i
            tree.insert(g.projectiles[i])
        }
    }
    g.mutex.Unlock()

    collisions := tree.getAllIntersections()

    for i := range collisions {
        e := collisions[i].entities 
        e1, e2 := e[0], e[1]

        if ok := EntitiesCollide(e1, e2); ok {
            g.KillEntity(e1)
            g.KillEntity(e2)
        }
    }
}

func (g *GameState)KillEntity(e Entity) {
    eIx := e.getGamestateIx()

    switch e.getEntityType() {
    case ENTITY_ENEMY:
        g.mutex.Lock()
        g.enemies[eIx].StartDying()
        enemy := g.enemies[eIx]
        g.mutex.Unlock()

        go g.SetEnemyDeathTimer(enemy, 500)
    case ENTITY_PROJECTILE:
        g.mutex.Lock()
        g.projectiles[eIx].StartDying()
        projectile := g.projectiles[eIx]
        g.mutex.Unlock()

        go g.SetProjectileDeathTimer(projectile, 500)
    }
}
