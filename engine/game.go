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

func (g *Game) UpdateFOV(p Player, dist int) {
	// Remove all previous FOV
	for _, t := range g.gameMap.fovList {
		t.inFOV = false
	}
	g.gameMap.fovList = nil

	// Cast rays from center of player to edges of a square centred on the player
	// - if the ray hits a wall stop and otherwise the tile is inFov and seen
	fovRange := dist
	for x := p.X - fovRange; x < p.X+fovRange+1; x++ {
		for y := p.Y - fovRange; y < p.Y+fovRange+1; y++ {
			if x < 0 || x >= g.gameMap.width || y < 0 || y >= g.gameMap.height {
				continue
			}

			// Step along the ray from the player to the edge of the square, stopping if we hit a wall
			//lint:ignore S1031
			points := p.Pos.RayCastTo(core.Pos{x, y}, fovRange)
			for _, t := range points {
				tile := &g.gameMap.tiles[t.X][t.Y]

				tile.inFOV = true
				tile.seen = true
				g.gameMap.fovList = append(g.gameMap.fovList, tile)

				// When tile blocks LOS, stop but after marking it as seen
				if tile.blocksLOS {
					break
				}
			}
		}
	}
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
	g.gameMap.makeRectRoom(10, 4, 7, 8)
	g.gameMap.makeRectRoom(7, 5, 3, 1)

	g.gameMap.tiles[13][8].makeWall()
	g.gameMap.tiles[13][9].makeWall()

	return g
}
