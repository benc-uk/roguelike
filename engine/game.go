package engine

import "roguelike/core"

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
