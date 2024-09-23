package engine

import (
	"roguelike/core"
)

// =====================================================================================================================
// Tile stuff
// =====================================================================================================================

type tileType int // TileType is an integer representing the type of tile

const (
	tileTypeFloor tileType = iota
	tileTypeWall
)

type tile struct {
	core.Pos
	Type       tileType
	seen       bool
	inFOV      bool
	blocksMove bool
	blocksLOS  bool
	entities   entityList
}

type Appearance struct {
	Graphic string
	Colour  string
	InFOV   bool
}

func newWall(x, y int) tile {
	return tile{Pos: core.Pos{X: x, Y: y}, Type: tileTypeWall, blocksMove: true, blocksLOS: true}
}

// nolint
func newFloor(pos core.Pos) tile {
	return tile{Pos: pos, Type: tileTypeFloor, blocksMove: false, blocksLOS: false}
}

func (t *tile) makeFloor() {
	t.Type = tileTypeFloor
	t.blocksMove = false
	t.blocksLOS = false
}

func (t *tile) makeWall() {
	t.Type = tileTypeWall
	t.blocksMove = true
	t.blocksLOS = true
}

func (t *tile) Entities() entityList {
	return t.entities
}

func (t *tile) placeItem(item *Item) {
	if item == nil {
		return
	}

	t.entities = append(t.entities, item)
	item.Pos = &t.Pos
}

func (t *tile) removeItem(item *Item) {
	t.entities.Remove(item)
}

// GetAppearance returns the appearance of the tile as a string
// to be used by the renderer and UI to display this tile
func (t *tile) GetAppearance(gameMap *GameMap) *Appearance {
	if !t.seen {
		return nil
	}

	if t.Type == tileTypeWall {
		return &Appearance{Graphic: "wall", InFOV: t.inFOV}
	}

	// If there are entities on this tile, return the appearance of the last one
	if !t.entities.IsEmpty() {
		// Creatures take precedence over items, and will be displayed first
		creatures := t.entities.AllCreatures()
		if len(creatures) > 0 {
			a := creatures[len(creatures)-1].Appearance()
			a.InFOV = t.inFOV

			return &a
		}

		items := t.entities.AllItems()
		if len(items) > 0 {
			a := items[len(items)-1].Appearance()
			a.InFOV = t.inFOV

			return &a
		}
	}

	return &Appearance{Graphic: "floor", InFOV: t.inFOV}
}

func (t *tile) BlocksMove() bool {
	for _, e := range t.entities {
		if e.BlocksMove() {
			return true
		}
	}
	return t.blocksMove
}

func (t *tile) BlocksLOS() bool {
	for _, e := range t.entities {
		if e.BlocksLOS() {
			return true
		}
	}
	return t.blocksLOS
}

// =====================================================================================================================
// Map stuff
// =====================================================================================================================

type GameMap struct {
	tiles   [][]tile // 2D array of tiles
	width   int      // Width of the map
	height  int      // Height of the map
	fovList []*tile  // List of all tiles in the FOV
}

func (m *GameMap) Tile(x, y int) *tile {
	return &m.tiles[x][y]
}

func (m *GameMap) TileAt(pos core.Pos) *tile {
	return &m.tiles[pos.X][pos.Y]
}

func (m *GameMap) Size() core.Size {
	return core.Size{Width: m.width, Height: m.height}
}

func (m *GameMap) Rect() core.Rect {
	return core.NewRect(0, 0, m.width, m.height)
}

func NewMap(width, height int) *GameMap {
	events.new("map_created", nil, "Map created")
	m := &GameMap{
		width:   width,
		height:  height,
		fovList: make([]*tile, 0),
	}

	m.tiles = make([][]tile, width)
	for x := range m.tiles {
		m.tiles[x] = make([]tile, height)
		for y := range m.tiles[x] {
			m.tiles[x][y] = newWall(x, y)
		}
	}

	return m
}

func (m *GameMap) makeRectRoom(x, y, w, h int) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			m.tiles[i][j].makeFloor()
		}
	}
}

// nolint
func (m *GameMap) seeAll() {
	for x := 0; x < m.width; x++ {
		for y := 0; y < m.height; y++ {
			m.tiles[x][y].seen = true
		}
	}
}
