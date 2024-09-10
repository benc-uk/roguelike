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
	blocksMove bool
	entities   entityList
}

type Appearance struct {
	Details string
	Seen    bool
	Hints   []string
}

func newWall(x, y int) tile {
	return tile{Pos: core.Pos{X: x, Y: y}, Type: tileTypeWall, blocksMove: false}
}

// nolint
func newFloor(pos core.Pos) tile {
	return tile{Pos: pos, Type: tileTypeFloor, blocksMove: true}
}

func (t *tile) makeFloor() {
	t.Type = tileTypeFloor
	t.blocksMove = true
}

// nolint
func (t *tile) makeWall() {
	t.Type = tileTypeWall
	t.blocksMove = false
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

// GetAppearance returns the appearance of the tile as a string
// to be used by the renderer and UI to display this tile
func (t *tile) GetAppearance() Appearance {
	if t.Type == tileTypeWall {
		return Appearance{Details: tileDispWall, Seen: t.seen, Hints: nil}
	}

	allItems := t.entities.AllItems()
	if len(allItems) > 0 {
		last := allItems[len(allItems)-1]
		a := last.Appearance()
		a.Seen = t.seen
		return a
	}

	return Appearance{Details: tileDispFloor, Seen: t.seen, Hints: nil}
}

// =====================================================================================================================
// Map stuff
// =====================================================================================================================

type GameMap struct {
	tiles  [][]tile
	width  int
	height int
}

func (m *GameMap) Tile(x, y int) *tile {
	return &m.tiles[x][y]
}

// TODO: Placeholder for now
func NewMap(width, height int) *GameMap {
	m := &GameMap{
		width:  width,
		height: height,
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
