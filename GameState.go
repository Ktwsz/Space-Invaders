package main

import (
    "time"
    gauss "github.com/chobie/go-gaussian"
    "math"
    rand "math/rand"
    "fmt"
)

const playerSpeed float64 = 1.0
const enemyHeight float64 = 8.0
const enemyCountColumn = 11
const playerProjectileSpeed = -3.5
const playerProjectileCooldown = 750
const enemyProjectileSpeed = 1.0
const enemyProjectileCooldown = 1250
const wallCount = 4
const WALL_POS_Y = 102

type GameState struct {
    bounds Vec2[int]

    player Player 

    enemies []*Enemy
    enemyColProjectileCooldown [enemyCountColumn]bool
    deadEnemies []int
    enemySpeed float64

    projectiles []*Projectile

    walls [wallCount]*Wall

    score int
    pauseState int

    enemyMoveDone chan bool
    wallsBodySet bool

}

func (g *GameState)Init() {
    g.bounds = Vec2[int]{x: GAME_WIDTH, y: GAME_HEIGHT}
    g.pauseState = GAME_RUNNING//TODO: remove later
    g.player.Init(g.bounds)

    g.enemySpeed = 2.0

    g.SpawnEnemiesRow(0, 3, 50, enemyCountColumn)
    g.SpawnEnemiesRow(1, 2, 20, enemyCountColumn)
    g.SpawnEnemiesRow(2, 1, 10, enemyCountColumn)

    g.SpawnWalls()

    enemyMoveTicker := time.NewTicker(enemyProjectileCooldown * time.Millisecond)
    g.enemyMoveDone = make(chan bool)
    go func() {
        for {
            select{
                case <-g.enemyMoveDone:
                    return
                case <-enemyMoveTicker.C:
                    g.MoveEnemies()
            }
        }
    }()
}

func (g *GameState)GameLoop() {
    if g.pauseState != GAME_RUNNING {
        return
    }

    g.removeDeadEnemies()
    g.RemoveDeadProjectiles()
    g.CheckForMissedProjectiles()
    g.CheckEnemiesInBounds()
    g.MoveProjectiles()
    g.EnemyShoot()
    g.HandleCollisions()
    g.HandleIfWon()
}

func (g *GameState)IsGameRunning() bool {
    return g.pauseState == GAME_RUNNING
}

func (g *GameState)HandleIfWon() {
    if len(g.enemies) == 0 {
        g.pauseState = GAME_WIN
        g.enemyMoveDone <- true
    }
}

func (g *GameState)SpawnWalls() {
    margin := g.bounds.x / wallCount - WALL_SIZE_X
    wallPos := Vec2[float64]{x: float64(margin)/2.0, y: WALL_POS_Y}
    for i := range wallCount {
        g.walls[i] = &Wall{position: wallPos,
                           hitbox: Vec2[float64]{x: WALL_SIZE_X, y: WALL_SIZE_Y},
                           hitboxReceiveMask: HITBOX_PROJECTILE,
                           hitboxSendMask: HITBOX_WALL,
                           gamestateIx: i}

        wallPos.x += float64(WALL_SIZE_X + margin)
    }
}

func (g *GameState)SetWallsBody(body WallBody) {
    g.wallsBodySet = true
    for i := range wallCount {
        g.walls[i].body = body
    }
}

func (g *GameState)GetWalls() [wallCount]*Wall {
    return g.walls
}

func (g *GameState)GetLastEnemyInCol() [enemyCountColumn]*Enemy {
    result := [enemyCountColumn]*Enemy{}
    for i, e := range g.enemies {
        col := e.rowData.x
        if result[col] == nil || e.rowData.y > result[col].rowData.y {
            result[col] = g.enemies[i]
        }
    }
    return result
}

func (g *GameState)GetChanceToShoot(gaussian *gauss.Gaussian, enemy *Enemy) bool {
    dist := math.Abs(enemy.position.x - g.player.position.x)
    chance := 0.6 * gaussian.Pdf(dist) * math.Sqrt(2 * math.Pi)
    
    return rand.Float64() <= chance
}

func (g *GameState)CanEnemyShoot(enemy *Enemy, col int) bool {
    return enemy != nil && enemy.deathState == STATE_ALIVE && !g.enemyColProjectileCooldown[col] 
}

func (g *GameState)EnemyShoot() {
    gaussDist := gauss.NewGaussian(0, 1)

    enemiesLast := g.GetLastEnemyInCol()

    for col := range enemyCountColumn {
        enemy := enemiesLast[col]
        if g.CanEnemyShoot(enemy, col) {
            if chance := g.GetChanceToShoot(gaussDist, enemy); chance {
                didShoot := g.SpawnEnemyprojectile(enemy, col)

                if didShoot {
                    g.enemyColProjectileCooldown[col] = true

                    go g.SetEnemyCooldownTimer(col, 750)
                }
            }
        }
    }
}

func (g *GameState)SpawnEnemyprojectile(enemy *Enemy, col int) bool {
    shotOnCooldown := g.enemyColProjectileCooldown[col]

    if shotOnCooldown {
        return false
    }
    projectile := Projectile{id: enemy.projectileId, 
                             frameCount: 4,
                             position: Vec2[float64]{x: enemy.position.x, y: enemy.position.y + enemy.spriteSize.y/2.0},
                             hitbox: Vec2[float64]{x: 3, y: 7},
                             spriteSize: Vec2[float64]{x: 3, y: 7},
                             speed: enemyProjectileSpeed,
                             hitboxReceiveMask: HITBOX_PROJECTILE | HITBOX_PLAYER | HITBOX_WALL,
                             hitboxSendMask: HITBOX_PROJECTILE | HITBOX_ENEMY,
                         }


    g.projectiles = append(g.projectiles, &projectile)

    return true
}

func (g *GameState)GetObjectsToDraw() []EntityDraw {
    enemiesLen := len(g.enemies)
    projectilesLen := len(g.projectiles)

    objects := make([]EntityDraw, projectilesLen + enemiesLen + 1)

    objects[0] = &g.player

    for i := range g.enemies {
        objects[i + 1] = g.enemies[i]
    }

    for i := range g.projectiles {
        objects[enemiesLen+i+1] = g.projectiles[i]
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

func (g *GameState)SpawnEnemiesRow(row int, enemySuffix int, enemyPoints int, count int) {
    position := Vec2[float64]{x: 6.0, y: (enemyHeight + 6.0) * float64(row) + 6.0}

    enemiesNew := make([]*Enemy, count)

    enemyId := fmt.Sprintf("enemy%d", enemySuffix)
    enemyProjectileId := fmt.Sprintf("enemy_projectile_%d", enemySuffix)

    for i := range count {
        enemiesNew[i] = &Enemy{}
        enemiesNew[i].Init(enemyId, enemyProjectileId, 3, position, Vec2[int]{x: i, y: row}, enemyPoints)

        position.x += 10.0//enemy x + margin
    }

    g.enemies = append(g.enemies, enemiesNew...)
}

func (g *GameState)removeDeadEnemies() {
    toRemove := make([]int, 0)
    for i, e := range g.enemies {
        if e.deathState == STATE_DEATH_END {
            toRemove = append(toRemove, i)
        }
    }

    g.enemies = RemoveIndexesMany(g.enemies, toRemove)
}

func (g *GameState)SetCooldownTimer(t time.Duration) {
    timer := time.NewTimer(t * time.Millisecond)

    <-timer.C

    g.player.shotOnCooldown = false
}

func (g *GameState)SetEnemyCooldownTimer(column int, t time.Duration) {
    timer := time.NewTimer(t * time.Millisecond)

    <-timer.C

    g.enemyColProjectileCooldown[column] = false
}

func (g *GameState)PlayerShoot() {
    if ok := g.SpawnPlayerProjectile(); ok {
        g.player.shotOnCooldown = true

        go g.SetCooldownTimer(750)
    }
}

func (g *GameState)SpawnPlayerProjectile() bool {
    shotOnCooldown := g.player.shotOnCooldown

    if shotOnCooldown {
        return false
    }
    projectile := Projectile{id: "player_projectile", 
                             frameCount: 1,
                             position: Vec2[float64]{x: g.player.position.x, y: g.player.position.y - g.player.spriteSize.y/2.0},
                             hitbox: Vec2[float64]{x: 1, y: 6},
                             spriteSize: Vec2[float64]{x: 1, y: 6},
                             speed: playerProjectileSpeed,
                             hitboxReceiveMask: HITBOX_PROJECTILE | HITBOX_ENEMY | HITBOX_WALL,
                             hitboxSendMask: HITBOX_PROJECTILE | HITBOX_PLAYER,
                         }


    g.projectiles = append(g.projectiles, &projectile)

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


func (g *GameState)RemoveDeadProjectiles() {
    toRemove := make([]int, 0)
    for i, p := range g.projectiles {
        if p.deathState == STATE_DEATH_END {
            toRemove = append(toRemove, i)
        }
    }

    g.projectiles = RemoveIndexesMany(g.projectiles, toRemove)
}

func (g *GameState)SetEnemyDeathTimer(enemy *Enemy, t time.Duration) {
    timer := time.NewTimer(t * time.Millisecond)

    <-timer.C

    enemy.deathState = STATE_DEATH_END
}

func (g *GameState)SetProjectileDeathTimer(projectile *Projectile, t time.Duration) {
    timer := time.NewTimer(t * time.Millisecond)

    <-timer.C
    projectile.deathState = STATE_DEATH_END
}

func (g *GameState)CheckForMissedProjectiles() {
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
        if g.enemies[i].rowData.y == row {
            g.enemies[i].Shift(Vec2[float64]{x: shiftX, y: 0})
        }
    }
}

func (g *GameState)EnemiesShiftDown() {
    for i := range g.enemies {
        g.enemies[i].Shift(Vec2[float64]{x: 0, y: enemyHeight + 6.0})        
    }
}

func (g *GameState)CheckEnemiesInBounds() {
    changeEnemiesDirection := false
    for i, e := range g.enemies {
        if IsOutOfBounds(g.bounds, e) {
            changeEnemiesDirection = true
            g.EnemiesShiftRow(e.rowData.y, i)
        } else {
        }
    }

    if changeEnemiesDirection {
        g.enemySpeed *= -1
        g.EnemiesShiftDown()
    }
}

func (g *GameState)HandleCollisions() {

    tree := QTreeInitFromBounds(g.bounds)

    g.player.handledCollisions = make(map[EntityHit]bool)
    tree.insert(&g.player)

    for i := range g.enemies {
        if g.enemies[i].deathState == STATE_ALIVE {
            g.enemies[i].gamestateIx = i
            g.enemies[i].handledCollisions = make(map[EntityHit]bool)
            tree.insert(g.enemies[i])
        }
    }

    for i := range g.projectiles {
        if g.projectiles[i].deathState == STATE_ALIVE {
            g.projectiles[i].gamestateIx = i
            g.projectiles[i].handledCollisions = make(map[EntityHit]bool)
            tree.insert(g.projectiles[i])
        }
    }

    for i := range g.walls {
        g.walls[i].handledCollisions = make(map[EntityHit]bool)
        tree.insert(g.walls[i])
    }

    collisions := tree.getAllIntersections()

    for i := range collisions {
        e := collisions[i].entities 
        e1, e2 := e[0], e[1]

        collision1 := g.SetCollisionHandled(e1, e2)
        collision2 := g.SetCollisionHandled(e2, e1)
        if !collision1 || !collision2 {
            continue
        }

        if e1.getEntityType() == ENTITY_WALL {
            wall := g.walls[e1.getGamestateIx()]
            if pos, didHit := wall.getHitPos(e2); didHit {
                if HitboxReceive(e1, e2) {
                    g.HitWall(e1, e2, pos)
                }
                if HitboxReceive(e2, e1) {
                    g.HitEntity(e2, e1)
                }
            }
        } else if e2.getEntityType() == ENTITY_WALL {
            wall := g.walls[e2.getGamestateIx()]
            if pos, didHit := wall.getHitPos(e1); didHit {
                if HitboxReceive(e1, e2) {
                    g.HitEntity(e1, e2)
                }
                if HitboxReceive(e2, e1) {
                    g.HitWall(e2, e1, pos)
                }
            }
        } else {
            if HitboxCollide(e1, e2) {
                if HitboxReceive(e2, e1) {
                    g.HitEntity(e1, e2)
                }

                if HitboxReceive(e1, e2) {
                    g.HitEntity(e2, e1)
                }
            }
        }
    }
}

func (g *GameState)SetCollisionHandled(e1 EntityHit, e2 EntityHit) bool {
    if e1.IsCollisionHandled(e2) {
        return false
    }
    
    eIx := e1.getGamestateIx()

    switch e1.getEntityType() {
    case ENTITY_ENEMY:
        enemy := g.enemies[eIx]
        enemy.handledCollisions[e2] = true
    case ENTITY_PROJECTILE:
        projectile := g.projectiles[eIx]
        projectile.handledCollisions[e2] = true
    case ENTITY_PLAYER:
        g.player.handledCollisions[e2] = true
    case ENTITY_WALL:
        wall := g.walls[eIx]
        wall.handledCollisions[e2] = true
    }

    return true
}

func (g *GameState)HitEntity(e EntityHit, sender EntityHit) {
    eIx := e.getGamestateIx()

    switch e.getEntityType() {
    case ENTITY_ENEMY:
        enemy := g.enemies[eIx]
        enemy.StartDying()
        g.score += g.enemies[eIx].points

        go g.SetEnemyDeathTimer(enemy, 500)

    case ENTITY_PROJECTILE:
        projectile := g.projectiles[eIx]
        projectile.StartDying()

        go g.SetProjectileDeathTimer(projectile, 500)

    case ENTITY_PLAYER:
        g.player.Hit()
        if g.player.lives <= 0 {
            g.pauseState = GAME_OVER
            g.enemyMoveDone <- true
        }
    }
}

func (g *GameState)HitWall(e EntityHit, sender EntityHit, pos Vec2[int]) {
    eIx := e.getGamestateIx()

    wall := g.walls[eIx]
    wall.Hit(pos)
}

func (g *GameState)GetPlayerLivesStr() string {
    return fmt.Sprintf("lives %d", g.player.lives)
}

func (g *GameState)GetScoreStr() string {
    return fmt.Sprintf("score %d", g.score)
}
