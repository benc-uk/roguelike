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

// TileType represents the type of tile
type tileType int

const (
	tileTypeFloor tileType = iota
	tileTypeWall
)

const maxTileItems = 10

// Tile is a single square on the map which can contain stuff
type tile struct {
	pos
	tileType   tileType
	seen       bool
	inFOV      bool
	blocksMove bool
	blocksLOS  bool
	items      entityList
	creature   *creature
	furniture  *furniture //nolint
}

// Appearance is a struct that holds the appearance of a tile
type Appearance struct {
	Graphic string
	Colour  string
	InFOV   bool
}

func newWall(x, y int) tile {
	return tile{
		pos:        pos{X: x, Y: y},
		tileType:   tileTypeWall,
		blocksMove: true,
		blocksLOS:  true,
		items:      nil,
	}
}

func (t *tile) makeFloor() {
	if t == nil {
		return
	}

	t.tileType = tileTypeFloor
	t.blocksMove = false
	t.blocksLOS = false
}

// nolint
func (t *tile) makeWall() {
	if t == nil {
		return
	}

	t.tileType = tileTypeWall
	t.blocksMove = true
	t.blocksLOS = true
}

// Add an item to the tile's item stack, up to the max allowed
func (t *tile) addItem(i *Item) bool {
	if t == nil {
		return false
	}

	// Only allow a certain number of items on a tile
	if t.items.Count() >= maxTileItems {
		return false
	}

	if i == nil {
		return false
	}

	t.items.Add(i)
	i.setPos(&t.pos)

	return true
}

// Place a creature on this tile
func (t *tile) placeCreature(c *creature) bool {
	// Only one creature per tile
	if t.creature != nil {
		return false
	}

	if c == nil {
		return false
	}

	t.creature = c
	c.setPos(&t.pos)
	c.currentTile = t

	return true
}

// Appearance returns the appearance of the tile as a string
// Used by the renderer and UI to display this tile
func (t *tile) Appearance() *Appearance {
	if t == nil || !t.seen {
		return nil
	}

	if t.tileType == tileTypeWall {
		return &Appearance{Graphic: "wall", InFOV: t.inFOV}
	}

	// Creatures take precedence over items, and will be displayed first
	if t.creature != nil {
		a := t.creature.Appearance()
		a.InFOV = t.inFOV

		return &a
	}

	// If there are items on this tile, return the appearance of the last one
	if last := t.items.Last(); last != nil {
		item := last.(*Item)
		appear := item.Appearance()
		appear.InFOV = t.inFOV

		return &appear
	}

	return &Appearance{Graphic: "floor", InFOV: t.inFOV}
}

// BlocksMove returns true if this tile currently blocks movement
func (t *tile) BlocksMove() bool {
	if t == nil {
		return true
	}

	if t.creature != nil {
		return t.creature.BlocksMove()
	}

	// Items don't block move, but lets check them anyway
	for _, i := range t.items {
		if i.BlocksMove() {
			return true
		}
	}

	return t.blocksMove
}

// BlocksLOS returns true if this tile currently blocks line of sight
func (t *tile) BlocksLOS() bool {
	if t == nil {
		return true
	}

	if t.creature != nil {
		return t.creature.BlocksLOS()
	}

	for _, e := range t.items {
		if e.BlocksLOS() {
			return true
		}
	}

	return t.blocksLOS
}

// ListItems returns a list of items on this tile
func (t *tile) ListItems() []Item {
	itemsOut := make([]Item, 0)
	for _, item := range t.items.AllItems() {
		itemsOut = append(itemsOut, *item)
	}

	return itemsOut
}

// Returns the creature on this tile, if any
func (t *tile) Creature() *creature {
	if t == nil {
		return nil
	}

	return t.creature
}

// ===== GameMap ==============================================================

// A game level is a 2D grid of tiles
type GameMap struct {
	size
	tiles [][]tile // 2D array of tiles, this holds the world

	fovList     []*tile // List of all tiles in the FOV
	depth       int     // Depth of the map
	description string  // Some human-readable description of the map
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

// Get a tile from x, y coordinates
func (m *GameMap) Tile(x, y int) *tile {
	return m.TileAt(pos{X: x, Y: y})
}

// Get a tile from a position
func (m *GameMap) TileAt(pos core.Pos) *tile {
	if !pos.InBounds(m.Width, m.Height) {
		return nil
	}

	return &m.tiles[pos.X][pos.Y]
}

// The size (width and height) of the map
func (m *GameMap) Size() core.Size {
	return m.size
}

// The description of the map
func (m *GameMap) Description() string {
	return m.description
}

// The depth/level of the map
func (m *GameMap) Depth() int {
	return m.depth
}

// Returns an adjacent tile in the given direction
func (m *GameMap) AdjacentTile(t *tile, dir core.Direction) *tile {
	destPos := t.pos.Add(dir.Pos())
	destTile := m.TileAt(destPos)

	if destTile == nil {
		return nil
	}

	return destTile
}

// Set an area of the map to be floor or wall
func (m *GameMap) setArea(wall bool, x, y, w, h int) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			if wall {
				m.tiles[i][j].makeWall()
			} else {
				m.tiles[i][j].makeFloor()
			}
		}
	}
}

// Set an area of the map to be floor or wall
func (m *GameMap) setAreaRect(wall bool, r core.Rect) {
	m.setArea(wall, r.X, r.Y, r.Width, r.Height)
}

// Reveal the entire map
// nolint
func (m *GameMap) revealMap() {
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			m.tiles[x][y].seen = true
		}
	}
}

// Find a random floor tile on the map
func (m *GameMap) randomFloorTile(noItems bool) *tile {
	for {
		x := rng.IntN(m.Width)
		y := rng.IntN(m.Height)
		t := m.Tile(x, y)

		if t.tileType == tileTypeFloor {
			if t.BlocksMove() {
				continue
			}

			if noItems && !t.items.IsEmpty() {
				continue
			}

			return t
		}
	}
}

// Dump the whole map to a PNG file called map.png
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
