package engine

import (
	"encoding/json"
	"fmt"
	"roguelike/core"
)

// ============================================================================
// This is core of the game engine
// The struct `Game` holds all the game state and all logic hangs off it
// ============================================================================

type Game struct {
	player  *Player
	gameMap *GameMap

	// Entity generators
	itemGen     *itemGenerator
	creatureGen *creatureGenerator
}

// Create a new game instance, it all starts here
func NewGame(dataFileDir string, seed uint64, viewDist int, listeners ...EventListener) *Game {
	seedRNG(seed)

	g := &Game{}

	// Reset the global event manager
	events = eventManager{}

	// Register event listeners
	events.addEventListeners(listeners...)

	generateMap(g, dataFileDir)

	g.player = NewPlayer(g.gameMap.randomFloorTile(true))
	g.player.fovDistance = viewDist

	levelText := fmt.Sprintf("You are on level %d of %s", g.Map().Depth(), g.Map().Description())
	events.new(EventMiscMessage, nil, "Welcome adventurer "+g.Player().Name())
	events.new(EventMiscMessage, nil, levelText)

	g.gameMap.dumpPNG()

	// DEBUG
	// g.player.backpack.Add(g.itemGen.createItem("meat"))

	g.updateFOV()
	return g
}

func (g *Game) Map() *GameMap {
	return g.gameMap
}

func (g *Game) Player() *Player {
	return g.player
}

// Update what the player can see, called after every action
func (g *Game) updateFOV() {
	p := g.player
	fovRange := p.fovDistance

	// Remove all previous FOV
	for _, t := range g.gameMap.fovList {
		t.inFOV = false
	}
	g.gameMap.fovList = nil

	// TODO: Implement a better FOV algorithm!

	// Cast rays from center of player to edges of a square centred on the player
	// - if the ray hits a wall stop and otherwise the tile is inFov and seen
	for x := p.X - fovRange; x < p.X+fovRange+1; x++ {
		for y := p.Y - fovRange; y < p.Y+fovRange+1; y++ {
			if x < 0 || x >= g.gameMap.Width || y < 0 || y >= g.gameMap.Height {
				continue
			}

			// Step along the ray from the player to the edge of the square, stopping if we hit a wall
			rayCoords := p.RayCastTo(pos{X: x, Y: y}, float64(fovRange)) // nolint
			for _, checkCoord := range rayCoords {
				tile := &g.gameMap.tiles[checkCoord.X][checkCoord.Y]

				// Check if title contains any entities that block LOS
				for _, entity := range tile.items {
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

// ViewPort is the area of the map that is centered on the player
// Implementations of the game should use this to render the game world
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

// MarshalJSON is a custom JSON marshaller for the Game struct
func (g *Game) MarshalJSON() ([]byte, error) {
	// TODO: This is not implemented yet, we need custom marshalling
	return json.Marshal(struct {
		Player  *Player  `json:"player"`
		GameMap *GameMap `json:"gameMap"`
	}{
		Player:  g.player,
		GameMap: g.gameMap,
	})
}
