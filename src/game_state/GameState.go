package game_state

import (
    "time"
    gauss "github.com/chobie/go-gaussian"
    "math"
    rand "math/rand"
    "fmt"

    "space_invaders/utils"
    "space_invaders/entity"
    entity_hitbox "space_invaders/entity/hitbox"
    entity_states "space_invaders/entity/states"
    entity_ids "space_invaders/entity/ids"
    "space_invaders/game_state/states"
)

const GAME_WIDTH = 150
const GAME_HEIGHT = 150

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
    bounds utils.Vec2[int]

    player entity.Player 

    enemies []*entity.Enemy
    enemyColProjectileCooldown [enemyCountColumn]bool
    deadEnemies []int
    enemySpeed float64

    projectiles []*entity.Projectile

    walls [wallCount]*entity.Wall

    score int
    pauseState int

    enemyMoveDone chan bool
    wallsBodySet bool

    soundQueue []string
}

func (g *GameState)Init() {
    g.bounds = utils.CreateVec(GAME_WIDTH, GAME_HEIGHT)
    g.pauseState = states.STARTING

    g.enemySpeed = 2.0

}

func (g *GameState)StartGame() {
    g.pauseState = states.RUNNING

    g.player.Init(g.bounds)

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
    if g.pauseState != states.RUNNING {
        return
    }

    g.removeDeadEnemies()
    g.RemoveDeadProjectiles()
    g.CheckEnemiesInBounds()
    g.MoveProjectiles()
    g.EnemyShoot()
    g.HandleCollisions()
    g.HandleIfWon()
}

func (g *GameState)IsGameRunning() bool {
    return g.pauseState == states.RUNNING
}

func (g *GameState)HandleIfWon() {
    if len(g.enemies) == 0 {
        g.pauseState = states.WIN
        g.enemyMoveDone <- true
    }
}

func (g *GameState)ClearSoundQueue() {
    g.soundQueue = []string{}
}

func (g *GameState)GetSoundQueue() []string {
    return g.soundQueue
}

func (g *GameState)SpawnWalls() {
    margin := g.bounds.X / wallCount - entity.WALL_SIZE_X
    wallPos := utils.CreateVec(float64(margin)/2.0, WALL_POS_Y)
    for i := range wallCount {
        g.walls[i] = entity.CreateWall(wallPos, i)

        wallPos.X += float64(entity.WALL_SIZE_X + margin)
    }
}

func (g GameState)GetWallsBodySet() bool {
    return g.wallsBodySet
}

func (g *GameState)SetWallsBody(body entity.WallBody) {
    g.wallsBodySet = true
    for i := range wallCount {
        g.walls[i].Body = body
    }
}

func (g *GameState)GetWalls() [wallCount]*entity.Wall {
    return g.walls
}

func (g *GameState)GetLastEnemyInCol() [enemyCountColumn]*entity.Enemy {
    result := [enemyCountColumn]*entity.Enemy{}
    for i, e := range g.enemies {
        col := e.GetRowData().X
        if result[col] == nil || e.GetRowData().Y > result[col].GetRowData().Y {
            result[col] = g.enemies[i]
        }
    }
    return result
}

func (g *GameState)GetChanceToShoot(gaussian *gauss.Gaussian, enemy *entity.Enemy) bool {
    dist := math.Abs(enemy.GetPosition().X - g.player.GetPosition().X)
    chance := 0.6 * gaussian.Pdf(dist) * math.Sqrt(2 * math.Pi)
    
    return rand.Float64() <= chance
}

func (g *GameState)CanEnemyShoot(enemy *entity.Enemy, col int) bool {
    return enemy != nil && enemy.DeathState == entity_states.ALIVE && !g.enemyColProjectileCooldown[col] 
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
                    g.soundQueue = append(g.soundQueue, "enemy_shoot")

                    go g.SetEnemyCooldownTimer(col, 750)
                }
            }
        }
    }
}

func (g *GameState)SpawnEnemyprojectile(enemy *entity.Enemy, col int) bool {
    shotOnCooldown := g.enemyColProjectileCooldown[col]

    if shotOnCooldown {
        return false
    }

    projectile := entity.CreateProjectile(
        enemy.GetProjectileId(),
        4,
        utils.CreateVec(enemy.GetPosition().X, enemy.GetPosition().Y + enemy.GetSpriteSize().Y/2.0),
        utils.CreateVec(3, 7).ToFloat64(),
        enemyProjectileSpeed,
        entity_hitbox.PLAYER,
        entity_hitbox.ENEMY,
    )

    g.projectiles = append(g.projectiles, &projectile)

    return true
}

func (g *GameState)GetObjectsToDraw() []utils.EntityDraw {
    enemiesLen := len(g.enemies)
    projectilesLen := len(g.projectiles)

    objects := make([]utils.EntityDraw, projectilesLen + enemiesLen + 1)

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
    g.player.Position.X -= playerSpeed
    if g.player.GetPosition().X - g.player.GetSpriteSize().X/2.0 < 0 {
        g.player.Position.X = g.player.GetSpriteSize().X/2.0
    }
}

func (g *GameState)PlayerMoveRight() {
    g.player.Position.X += playerSpeed
    if g.player.GetPosition().X + g.player.GetSpriteSize().X/2.0 > float64(g.bounds.X) {
        g.player.Position.X = float64(g.bounds.X) - g.player.GetSpriteSize().X/2.0
    }
}

func (g *GameState)SpawnEnemiesRow(row int, enemySuffix int, enemyPoints int, count int) {
    position := utils.CreateVec(6.0, (enemyHeight + 6.0) * float64(row) + 6.0)

    enemiesNew := make([]*entity.Enemy, count)

    enemyId := fmt.Sprintf("enemy%d", enemySuffix)
    enemyProjectileId := fmt.Sprintf("enemy_projectile_%d", enemySuffix)

    for i := range count {
        enemiesNew[i] = &entity.Enemy{}
        enemiesNew[i].Init(enemyId, enemyProjectileId, 3, position, utils.CreateVec(i, row), enemyPoints)

        position.X += 10.0//enemy x + margin
    }

    g.enemies = append(g.enemies, enemiesNew...)
}

func (g *GameState)removeDeadEnemies() {
    toRemove := make([]int, 0)
    for i, e := range g.enemies {
        if e.DeathState == entity_states.DEATH_END {
            toRemove = append(toRemove, i)
        }
    }

    g.enemies = utils.RemoveIndexesMany(g.enemies, toRemove)
}

func (g *GameState)SetCooldownTimer(t time.Duration) {
    timer := time.NewTimer(t * time.Millisecond)

    <-timer.C

    g.player.ShotOnCooldown = false
}

func (g *GameState)SetEnemyCooldownTimer(column int, t time.Duration) {
    timer := time.NewTimer(t * time.Millisecond)

    <-timer.C

    g.enemyColProjectileCooldown[column] = false
}

func (g *GameState)PlayerShoot() {
    if ok := g.SpawnPlayerProjectile(); ok {
        g.player.ShotOnCooldown = true
        g.soundQueue = append(g.soundQueue, "player_shoot")

        go g.SetCooldownTimer(750)
    }
}

func (g *GameState)SpawnPlayerProjectile() bool {
    shotOnCooldown := g.player.GetShotOnCooldown()

    if shotOnCooldown {
        return false
    }

    projectile := entity.CreateProjectile("player_projectile",
        1,
        utils.CreateVec(g.player.GetPosition().X, g.player.GetPosition().Y - g.player.GetSpriteSize().Y/2.0),
        utils.CreateVec(1, 6).ToFloat64(),
        playerProjectileSpeed,
        entity_hitbox.ENEMY,
        entity_hitbox.PLAYER,
    )

    g.projectiles = append(g.projectiles, &projectile)

    return true
}

func (g *GameState)MoveEnemies() {
    for i := range g.enemies {
        g.enemies[i].Move(g.enemySpeed)
    }

    g.soundQueue = append(g.soundQueue, "enemy_move")
}

func (g *GameState)MoveProjectiles() {
    for i := range g.projectiles {
        g.projectiles[i].Move()
    }
}


func (g *GameState)RemoveDeadProjectiles() {
    toRemove := make([]int, 0)
    for i, p := range g.projectiles {
        if p.DeathState == entity_states.DEATH_END {
            toRemove = append(toRemove, i)
        }
    }

    g.projectiles = utils.RemoveIndexesMany(g.projectiles, toRemove)
}

func (g *GameState)SetEnemyDeathTimer(enemy *entity.Enemy, t time.Duration) {
    timer := time.NewTimer(t * time.Millisecond)

    <-timer.C

    enemy.DeathState = entity_states.DEATH_END
}

func (g *GameState)SetProjectileDeathTimer(projectile *entity.Projectile, t time.Duration) {
    timer := time.NewTimer(t * time.Millisecond)

    <-timer.C
    projectile.DeathState = entity_states.DEATH_END
}

func (g *GameState)EnemiesShiftRow(row int, enemyIx int) {
    var shiftX float64

    e := g.enemies[enemyIx]
    if e.GetPosition().X + e.GetSpriteSize().X > float64(g.bounds.X) {
        shiftX = -(e.GetPosition().X + e.GetSpriteSize().X - float64(g.bounds.X))
    } else {
        shiftX = e.GetSpriteSize().X - e.GetPosition().X
    }

    for i := range g.enemies {
        if g.enemies[i].GetRowData().Y == row {
            g.enemies[i].Shift(utils.CreateVec(shiftX, 0))
        }
    }
}

func (g *GameState)EnemiesShiftDown() {
    for _, e := range g.enemies {
        if e.GetPosition().Y > 70 {
            return
        } 
    }

    for i := range g.enemies {
        g.enemies[i].Shift(utils.CreateVec(0, enemyHeight + 6.0))
    }
}

func (g *GameState)CheckEnemiesInBounds() {
    changeEnemiesDirection := false
    for i, e := range g.enemies {
        if utils.IsOutOfBounds(g.bounds, e) {
            changeEnemiesDirection = true
            g.EnemiesShiftRow(e.GetRowData().Y, i)
        } 
    }

    if changeEnemiesDirection {
        g.enemySpeed *= -1
        g.EnemiesShiftDown()
    }
}

func (g *GameState)HandleCollisions() {

    tree := utils.QTreeInitFromBounds(g.bounds)

    g.player.HandledCollisions = make(map[utils.EntityHit]bool)
    tree.Insert(&g.player)

    for i := range g.enemies {
        if g.enemies[i].DeathState == entity_states.ALIVE {
            g.enemies[i].GamestateIx = i
            g.enemies[i].HandledCollisions = make(map[utils.EntityHit]bool)
            tree.Insert(g.enemies[i])
        }
    }

    for i := range g.projectiles {
        if g.projectiles[i].DeathState == entity_states.ALIVE {
            g.projectiles[i].GamestateIx = i
            g.projectiles[i].HandledCollisions = make(map[utils.EntityHit]bool)
            tree.Insert(g.projectiles[i])
        }
    }

    for i := range g.walls {
        g.walls[i].HandledCollisions = make(map[utils.EntityHit]bool)
        tree.Insert(g.walls[i])
    }

    collisions := tree.GetAllIntersections()

    for i := range collisions {
        e := collisions[i].GetEntities()
        e1, e2 := e[0], e[1]

        collision1 := g.SetCollisionHandled(e1, e2)
        collision2 := g.SetCollisionHandled(e2, e1)
        if !collision1 || !collision2 {
            continue
        }

        if e1.GetEntityType() == entity_ids.WALL {
            wall := g.walls[e1.GetGamestateIx()]
            if pos, didHit := wall.GetHitPos(e2); didHit {
                if utils.HitboxReceive(e1, e2) {
                    g.HitWall(e1, e2, pos)
                }
                if utils.HitboxReceive(e2, e1) {
                    g.HitEntity(e2, e1)
                }
            }
        } else if e2.GetEntityType() == entity_ids.WALL {
            wall := g.walls[e2.GetGamestateIx()]
            if pos, didHit := wall.GetHitPos(e1); didHit {
                if utils.HitboxReceive(e1, e2) {
                    g.HitEntity(e1, e2)
                }
                if utils.HitboxReceive(e2, e1) {
                    g.HitWall(e2, e1, pos)
                }
            }
        } else {
            if utils.HitboxCollide(e1, e2) {
                if utils.HitboxReceive(e2, e1) {
                    g.HitEntity(e1, e2)
                }

                if utils.HitboxReceive(e1, e2) {
                    g.HitEntity(e2, e1)
                }
            }
        }
    }
}

func (g *GameState)SetCollisionHandled(e1 utils.EntityHit, e2 utils.EntityHit) bool {
    if e1.IsCollisionHandled(e2) {
        return false
    }
    
    eIx := e1.GetGamestateIx()

    switch e1.GetEntityType() {
    case entity_ids.ENEMY:
        enemy := g.enemies[eIx]
        enemy.HandledCollisions[e2] = true
    case entity_ids.PROJECTILE:
        projectile := g.projectiles[eIx]
        projectile.HandledCollisions[e2] = true
    case entity_ids.PLAYER:
        g.player.HandledCollisions[e2] = true
    case entity_ids.WALL:
        wall := g.walls[eIx]
        wall.HandledCollisions[e2] = true
    }

    return true
}

func (g *GameState)HitEntity(e utils.EntityHit, sender utils.EntityHit) {
    eIx := e.GetGamestateIx()

    switch e.GetEntityType() {
    case entity_ids.ENEMY:
        enemy := g.enemies[eIx]
        g.soundQueue = append(g.soundQueue, "enemy_die")
        enemy.StartDying()
        g.score += g.enemies[eIx].GetPoints()

        go g.SetEnemyDeathTimer(enemy, 500)

    case entity_ids.PROJECTILE:
        projectile := g.projectiles[eIx]
        projectile.StartDying()

        go g.SetProjectileDeathTimer(projectile, 500)

    case entity_ids.PLAYER:
        g.player.Hit()
        g.soundQueue = append(g.soundQueue, "player_hit")
        if g.player.GetLives() <= 0 {
            g.pauseState = states.OVER
            g.enemyMoveDone <- true
        }
    }
}

func (g *GameState)HitWall(e utils.EntityHit, sender utils.EntityHit, pos utils.Vec2[int]) {
    eIx := e.GetGamestateIx()

    wall := g.walls[eIx]
    wall.Hit(pos)
}

func (g *GameState)GetPlayerLivesStr() string {
    return fmt.Sprintf("lives %d", g.player.GetLives())
}

func (g *GameState)GetScoreStr() string {
    return fmt.Sprintf("score %d", g.score)
}

func (g *GameState)GetScoreResultStr() string {
    return fmt.Sprintf("Your score was: %d", g.score)
}

func (g *GameState)GetPauseState() int {
    return g.pauseState
}

