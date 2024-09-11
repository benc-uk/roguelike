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

	// TODO: Absolutely shit placeholder
	// Update all tiles in radius around the player as seen and in FOV
	for dx := -3; dx <= 3; dx++ {
		for dy := -3; dy <= 3; dy++ {
			x := p.Pos.X + dx
			y := p.Pos.Y + dy
			if x >= 0 && x < g.gameMap.width && y >= 0 && y < g.gameMap.height {
				g.gameMap.tiles[x][y].seen = true
				g.gameMap.tiles[x][y].inFOV = true
				g.gameMap.fovList = append(g.gameMap.fovList, &g.gameMap.tiles[x][y])
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
	g.gameMap.makeRectRoom(10, 4, 6, 5)
	g.gameMap.makeRectRoom(7, 5, 3, 1)

	return g
}
