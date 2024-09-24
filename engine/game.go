package engine

import (
	"roguelike/core"
)

type Game struct {
	player  *Player
	gameMap *GameMap

	// Entity factories
	itemFactory itemFactoryDB
}

func (g *Game) Map() *GameMap {
	return g.gameMap
}

func (g *Game) Player() *Player {
	return g.player
}

func (g *Game) UpdateFOV(fovRange int) {
	p := g.player
	// Remove all previous FOV
	for _, t := range g.gameMap.fovList {
		t.inFOV = false
	}
	g.gameMap.fovList = nil

	// Cast rays from center of player to edges of a square centred on the player
	// - if the ray hits a wall stop and otherwise the tile is inFov and seen
	for x := p.X - fovRange; x < p.X+fovRange+1; x++ {
		for y := p.Y - fovRange; y < p.Y+fovRange+1; y++ {
			if x < 0 || x >= g.gameMap.Width || y < 0 || y >= g.gameMap.Height {
				continue
			}

			// Step along the ray from the player to the edge of the square, stopping if we hit a wall
			rayCoords := p.Pos.RayCastTo(core.Pos{x, y}, float64(fovRange)) // nolint
			for _, checkCoord := range rayCoords {
				tile := &g.gameMap.tiles[checkCoord.X][checkCoord.Y]

				// Check if title contains any entities that block LOS
				for _, entity := range tile.entities {
					if entity.BlocksLOS() {
						tile.blocksLOS = true
						break
					}
				}

				tile.inFOV = true
				tile.seen = true
				g.gameMap.fovList = append(g.gameMap.fovList, tile)

				// When tile blocks LOS, stop but *after* marking it as seen
				if tile.blocksLOS {
					break
				}
			}
		}
	}
}

func (g *Game) AddEventListener(listener func(GameEvent)) {
	events.AddEventListener(listener)
}

// Create a new game instance
func NewGame(dataFileDir string) *Game {
	g := &Game{
		player: &Player{
			Pos:  core.Pos{X: 2, Y: 2},
			name: "Wizard Bob",
		},
	}

	// Reset the global event manager
	events = EventManager{}

	var err error
	g.itemFactory, err = newItemFactory(dataFileDir + "/items.yaml")
	if err != nil {
		panic(err)
	}

	// *******************************
	// HACK: PLACEHOLDER MAP SETUP
	// *******************************
	// Tiny
	//g.gameMap = NewMap(32, 32, 1)
	//g.gameMap.GenerateBSP(3) // 4 or 5 also works

	// Smaller
	g.gameMap = NewMap(40, 40, 1)
	g.gameMap.GenerateBSP(4, g.itemFactory) // 4 or 5 also works
	g.gameMap.description = "a small dungeon"

	// Small
	//g.gameMap = NewMap(48, 48, 1)
	//g.gameMap.GenerateBSP(6) // 5 also works

	// Medium
	// g.gameMap = NewMap(64, 64, 1)
	// g.gameMap.GenerateBSP(6)

	// HACK: Dump the map to a PNG file
	g.gameMap.dumpPNG()

	g.player.Pos = g.gameMap.randomFloorTile(true).Pos

	return g
}

// ViewPort is the area of the map that is centered on the player
func (g Game) GetViewPort(vpWidth, vpHeight int) core.Rect {
	vp := core.NewRect((g.player.X - vpWidth/2), (g.player.Y - vpHeight/2), vpWidth, vpHeight)

	if vp.X < 0 {
		vp.X = 0
	}
	if vp.Y < 0 {
		vp.Y = 0
	}
	if vp.X+vp.Width > g.gameMap.Width {
		vp.X = g.gameMap.Width - vp.Width
	}
	if vp.Y+vp.Height > g.gameMap.Height {
		vp.Y = g.gameMap.Height - vp.Height
	}

	return vp
}
