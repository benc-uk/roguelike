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
	blocksLOS  bool // nolint
	entities   entityList
}

type Appearance struct {
	Details string
	Hints   []string
	InFOV   bool
}

func newWall(x, y int) tile {
	return tile{Pos: core.Pos{X: x, Y: y}, Type: tileTypeWall, blocksMove: true}
}

// nolint
func newFloor(pos core.Pos) tile {
	return tile{Pos: pos, Type: tileTypeFloor, blocksMove: false}
}

func (t *tile) makeFloor() {
	t.Type = tileTypeFloor
	t.blocksMove = false
}

// nolint
func (t *tile) makeWall() {
	t.Type = tileTypeWall
	t.blocksMove = true
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

var unseenTileAppearance = Appearance{Details: "blank", Hints: nil}

// GetAppearance returns the appearance of the tile as a string
// to be used by the renderer and UI to display this tile
func (t *tile) GetAppearance() Appearance {
	if !t.seen {
		return unseenTileAppearance
	}

	if t.Type == tileTypeWall {
		return Appearance{Details: "wall", Hints: nil, InFOV: t.inFOV}
	}

	// If there are entities on this tile, return the appearance of the last one
	if !t.entities.IsEmpty() {
		creatures := t.entities.AllCreatures()
		if len(creatures) > 0 {
			a := creatures[len(creatures)-1].Appearance()
			a.InFOV = t.inFOV
			return a
		}

		items := t.entities.AllItems()
		if len(items) > 0 {
			a := items[len(items)-1].Appearance()
			a.InFOV = t.inFOV
			return a
		}
	}

	return Appearance{Details: "floor", Hints: nil, InFOV: t.inFOV}
}

// =====================================================================================================================
// Map stuff
// =====================================================================================================================

type GameMap struct {
	tiles  [][]tile
	width  int
	height int
	inFOV  []*tile
}

func (m *GameMap) Tile(x, y int) *tile {
	return &m.tiles[x][y]
}

// TODO: Placeholder for now
func NewMap(width, height int) *GameMap {
	m := &GameMap{
		width:  width,
		height: height,
		inFOV:  make([]*tile, 0),
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
