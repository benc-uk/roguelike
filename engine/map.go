package engine

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	_ "image/png"
	"os"
	"roguelike/core"
)

// ============================================================================
// GameMap represents the structure of the game world as a 2D grid of tiles
// ============================================================================

// ===== Tiles ================================================================

type tileType int // TileType is an integer representing the type of tile

const (
	tileTypeFloor tileType = iota
	tileTypeWall
)

type tile struct {
	pos
	tileType   tileType
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
	return tile{pos: pos{X: x, Y: y}, tileType: tileTypeWall, blocksMove: true, blocksLOS: true}
}

// nolint
func newFloor(pos core.Pos) tile {
	return tile{pos: pos, tileType: tileTypeFloor, blocksMove: false, blocksLOS: false}
}

func (t *tile) makeFloor() {
	t.tileType = tileTypeFloor
	t.blocksMove = false
	t.blocksLOS = false
}

// nolint
func (t *tile) makeWall() {
	t.tileType = tileTypeWall
	t.blocksMove = true
	t.blocksLOS = true
}

func (t *tile) placeItem(item *Item) {
	if item == nil {
		return
	}

	t.entities = append(t.entities, item)
	item.pos = &t.pos
}

func (t *tile) removeItem(item *Item) {
	t.entities.Remove(item)
}

// GetAppearance returns the appearance of the tile as a string
// to be used by the renderer and UI to display this tile
func (t *tile) GetAppearance(gameMap *GameMap) *Appearance {
	if t == nil {
		return nil
	}

	if !t.seen {
		return nil
	}

	if t.tileType == tileTypeWall {
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
	if t == nil {
		return true
	}

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

// ===== GameMap ==============================================================

type GameMap struct {
	size
	tiles [][]tile // 2D array of tiles

	fovList     []*tile // List of all tiles in the FOV
	depth       int     // Depth of the map
	description string
}

// NewMap creates a new map with the given width and height
// Note that the map is initially filled with walls, and you should
// call the generation functions to create structure in the map
func NewMap(width, height, depth int) *GameMap {
	m := &GameMap{
		size:        size{Width: width, Height: height},
		fovList:     make([]*tile, 0),
		depth:       depth,
		description: "Empty map",
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

func (m *GameMap) Tile(x, y int) *tile {
	p := pos{X: x, Y: y}
	return m.TileAt(p)
}

func (m *GameMap) TileAt(pos core.Pos) *tile {
	if !pos.InBounds(m.Width, m.Height) {
		return nil
	}

	return &m.tiles[pos.X][pos.Y]
}

func (m *GameMap) Size() core.Size {
	return m.size
}

func (m *GameMap) Description() string {
	return m.description
}

func (m *GameMap) Depth() int {
	return m.depth
}

func (m *GameMap) floorArea(x, y, w, h int) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			m.tiles[i][j].makeFloor()
		}
	}
}

func (m *GameMap) floorAreaRect(r core.Rect) {
	m.floorArea(r.X, r.Y, r.Width, r.Height)
}

// nolint
func (m *GameMap) revealMap() {
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			m.tiles[x][y].seen = true
		}
	}
}

func (m *GameMap) randomFloorTile(noItems bool) tile {
	for {
		x := rng.IntN(m.Width)
		y := rng.IntN(m.Height)

		if m.Tile(x, y).tileType == tileTypeFloor {
			if noItems && !m.Tile(x, y).entities.IsEmpty() {
				return *m.Tile(x, y)
			} else {
				return *m.Tile(x, y)
			}
		}
	}
}

// nolint
func (m *GameMap) dumpPNG() {
	tilesize := 16
	// Create a new image
	img := image.NewRGBA(image.Rect(0, 0, m.Width*tilesize, m.Height*tilesize))

	// Draw the map
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			tile := m.Tile(x, y)
			if tile.tileType == tileTypeWall {
				draw.Draw(img, image.Rect(x*tilesize, y*tilesize, x*tilesize+tilesize, y*tilesize+tilesize), &image.Uniform{color.Black}, image.Point{}, draw.Src)
			} else {
				draw.Draw(img, image.Rect(x*tilesize, y*tilesize, x*tilesize+tilesize, y*tilesize+tilesize), &image.Uniform{color.White}, image.Point{}, draw.Src)
			}
		}
	}

	// Encode as PNG file
	file := "map.png"
	f, _ := os.Create(file)
	_ = png.Encode(f, img)
}
