package main

const playerSpeed float64 = 1.0
const enemyHeight float64 = 8.0
const enemyCountPerRow = 11

type GameState struct {
    bounds Vec2[int]
    player Player 
    enemies []Enemy
}

func (g *GameState)Init() {
    g.bounds = Vec2[int]{x: 160, y: 120}
    g.player.Init(g.bounds)

    g.SpawnEnemiesRow(0, "enemy1", enemyCountPerRow)
    g.SpawnEnemiesRow(1, "enemy2", enemyCountPerRow)
    g.SpawnEnemiesRow(2, "enemy3", enemyCountPerRow)
}

func (g *GameState)GetObjectsToDraw() []Entity{
    objects := make([]Entity, len(g.enemies) + 1)
    objects[0] = &g.player

    for i := range g.enemies {
        objects[i + 1] = &g.enemies[i]
    }

    return objects
}

func (g *GameState)MovePlayerLeft() {
    g.player.position.x -= playerSpeed
    if g.player.position.x - g.player.spriteSize.x/2.0 < 0 {
        g.player.position.x = g.player.spriteSize.x/2.0
    }
}

func (g *GameState)MovePlayerRight() {
    g.player.position.x += playerSpeed
    if g.player.position.x + g.player.spriteSize.x/2.0 > float64(g.bounds.x) {
        g.player.position.x = float64(g.bounds.x) - g.player.spriteSize.x/2.0
    }
}

func (g *GameState)SpawnEnemiesRow(row int, enemyId string, count int) {
    position := Vec2[float64]{x: 6.0, y: (enemyHeight + 6.0) * float64(row)}

    enemiesNew := make([]Enemy, count)
    for i  := range count {
        enemiesNew[i] = Enemy{}
        enemiesNew[i].Init(enemyId, position)

        position.x += 10.0//enemy x + margin
    }

    g.enemies = append(g.enemies, enemiesNew...)
}
