package engine

import (
	"roguelike/core"
)

type GameState int // TileType is an integer representing the type of tile

const (
	GameStateTitle GameState = iota
	GameStatePlaying
	GameStateGameOver
	GameStatePlayerGen
)

type Game struct {
	State   GameState
	player  *Player
	gameMap *GameMap
}

func (g *Game) Map() *GameMap {
	return g.gameMap
}

func (g *Game) Player() *Player {
	return g.player
}

func (g *Game) UpdateFOV(p Player) {
	// Remove all previous FOV
	for _, t := range g.gameMap.fovList {
		t.inFOV = false
	}
	g.gameMap.fovList = nil

	// cast rays from center of player to edges of a 7x7 square centered on the player
	// if the ray hits a wall stop and otherwise the tile is inFov and seen
	for x := p.X - 3; x < p.X+4; x++ {
		for y := p.Y - 3; y < p.Y+4; y++ {
			if x < 0 || x >= g.gameMap.width || y < 0 || y >= g.gameMap.height {
				continue
			}
			if castRay(g.gameMap, p.X, p.Y, x, y) {
				g.gameMap.tiles[x][y].inFOV = true
				g.gameMap.tiles[x][y].seen = true
				g.gameMap.fovList = append(g.gameMap.fovList, &g.gameMap.tiles[x][y])
			}
		}
	}
}

func castRay(gameMap *GameMap, x0, y0, x1, y1 int) bool {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := -1
	sy := -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	err := dx - dy

	leave := 1
	for {
		if leave == 0 {
			return false
		}
		if gameMap.tiles[x0][y0].blocksLOS {
			leave = 0
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// TODO: All a massive placeholder for now
func NewGame() *Game {
	g := &Game{
		State: GameStateTitle,
		player: &Player{
			Pos:  core.Pos{X: 2, Y: 2},
			name: "Wizard Bob",
		},
		gameMap: NewMap(40, 40),
	}

	sword := itemFactory.CreateItem("sword")
	potion := itemFactory.CreateItem("potion")
	door := itemFactory.CreateItem("door")
	rat := itemFactory.CreateItem("rat")
	poison := itemFactory.CreateItem("potion_poison")

	g.gameMap.tiles[3][4].placeItem(potion)
	g.gameMap.tiles[12][5].placeItem(poison)
	g.gameMap.tiles[6][6].placeItem(sword)
	g.gameMap.tiles[7][5].placeItem(door)
	g.gameMap.tiles[14][7].placeItem(rat)

	g.gameMap.makeRectRoom(2, 2, 5, 5)
	g.gameMap.makeRectRoom(10, 4, 6, 5)
	g.gameMap.makeRectRoom(7, 5, 3, 1)

	return g
}
