package engine

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	_ "image/png"
	"math/rand"
	"os"
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

// nolint
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
	if t == nil {
		return nil
	}

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
	core.Size
	tiles [][]tile // 2D array of tiles

	fovList     []*tile // List of all tiles in the FOV
	depth       int     // Depth of the map
	description string
}

func (m *GameMap) Tile(x, y int) *tile {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return nil
	} else {
		return &m.tiles[x][y]
	}
}

func (m *GameMap) TileAt(pos core.Pos) *tile {
	return &m.tiles[pos.X][pos.Y]
}

func (m *GameMap) Rect() core.Rect {
	return core.NewRect(0, 0, m.Width, m.Height)
}

// NewMap creates a new map with the given width and height
// Note that the map is initially filled with walls, and you should
// call the generation functions to create structure in the map
func NewMap(width, height, depth int) *GameMap {
	m := &GameMap{
		Size:        core.Size{Width: width, Height: height},
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
		x := rand.Intn(m.Width)
		y := rand.Intn(m.Height)

		if m.Tile(x, y).Type == tileTypeFloor {
			if noItems && !m.Tile(x, y).entities.IsEmpty() {
				return *m.Tile(x, y)
			} else {
				return *m.Tile(x, y)
			}
		}
	}
}

func (m *GameMap) dumpPNG() {
	tilesize := 16
	// Create a new image
	img := image.NewRGBA(image.Rect(0, 0, m.Width*tilesize, m.Height*tilesize))

	// Draw the map
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			tile := m.Tile(x, y)
			if tile.Type == tileTypeWall {
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
