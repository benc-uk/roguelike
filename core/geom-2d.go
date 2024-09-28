package core

// ====================================================================================================================
// General 2D geometry and direction utilities, such as positions, rectangles, and directions
// ====================================================================================================================

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// Size represents a width and height
type Size struct {
	Width  int
	Height int
}

func (s Size) String() string {
	return fmt.Sprintf("%dx%d", s.Width, s.Height)
}

func (s Size) Area() int {
	return s.Width * s.Height
}

// ===== Rectangles =====

// Rect represents a rectangle with a position (x & y) and size (width & height)
type Rect struct {
	Pos
	Size
}

func (r Rect) String() string {
	return fmt.Sprintf("%s %s", r.Pos, r.Size)
}

func (r Rect) Contains(p Pos) bool {
	return p.X >= r.X && p.X < r.X+r.Width && p.Y >= r.Y && p.Y < r.Y+r.Height
}

func (r Rect) IntersectingRect(other Rect) Rect {
	x1 := MaxInt(r.X, other.X)
	y1 := MaxInt(r.Y, other.Y)
	x2 := MinInt(r.X+r.Width, other.X+other.Width)
	y2 := MinInt(r.Y+r.Height, other.Y+other.Height)
	if x2 > x1 && y2 > y1 {
		return Rect{Pos{x1, y1}, Size{x2 - x1, y2 - y1}}
	}
	return Rect{}
}

func (r Rect) Center() Pos {
	if r.Width == 0 || r.Height == 0 {
		return Pos{}
	}

	return Pos{r.X + r.Width/2, r.Y + r.Height/2}
}

func (r Rect) RandomPos(rng *rand.Rand) Pos {
	if r.Width <= 0 || r.Height <= 0 {
		return Pos{}
	}

	return Pos{rng.IntN(r.Width) + r.X, rng.IntN(r.Height) + r.Y}
}

func (r Rect) Area() int {
	return r.Width * r.Height
}

func (r Rect) ContainsPos(p Pos) bool {
	return p.X >= r.X && p.X < r.X+r.Width && p.Y >= r.Y && p.Y < r.Y+r.Height
}

func NewRect(x, y, width, height int) Rect {
	return Rect{Pos{x, y}, Size{width, height}}
}

// ===== Positions =====

// Pos represents a 2D position with X and Y coordinates
type Pos struct {
	X int `yaml:"x"`
	Y int `yaml:"y"`
}

func (p Pos) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func RandomPos(maxX, maxY int) Pos {
	return Pos{rand.IntN(maxX), rand.IntN(maxY)}
}

func (p Pos) Add(p2 Pos) Pos {
	return Pos{p.X + p2.X, p.Y + p2.Y}
}

func (p Pos) Sub(p2 Pos) Pos {
	return Pos{p.X - p2.X, p.Y - p2.Y}
}

func (p Pos) Distance(p2 Pos) float64 {
	return math.Sqrt(float64((p.X-p2.X)*(p.X-p2.X) + (p.Y-p2.Y)*(p.Y-p2.Y)))
}

func (p Pos) InBounds(width, height int) bool {
	return p.X >= 0 && p.X < width && p.Y >= 0 && p.Y < height
}

func (p Pos) InRect(r Rect) bool {
	return p.X >= r.X && p.X < r.X+r.Width && p.Y >= r.Y && p.Y < r.Y+r.Height
}

func (p Pos) IsNeighbour(p2 Pos) bool {
	return p.Distance(p2) == 1
}

func (p Pos) NeighboursCardinal() []Pos {
	return []Pos{
		{p.X - 1, p.Y},
		{p.X + 1, p.Y},
		{p.X, p.Y - 1},
		{p.X, p.Y + 1},
	}
}

func (p Pos) NeighboursAll() []Pos {
	return []Pos{
		{p.X - 1, p.Y},
		{p.X + 1, p.Y},
		{p.X, p.Y - 1},
		{p.X, p.Y + 1},
		{p.X - 1, p.Y - 1},
		{p.X + 1, p.Y - 1},
		{p.X - 1, p.Y + 1},
		{p.X + 1, p.Y + 1},
	}
}

func (p Pos) RayCastTo(p2 Pos, maxDist float64) []Pos {
	// Bresenham's line algorithm
	x0, y0 := p.X, p.Y
	x1, y1 := p2.X, p2.Y
	dx := AbsInt(x1 - x0)
	dy := AbsInt(y1 - y0)
	sx := 1
	sy := 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}

	err := dx - dy
	var points []Pos
	dist := 0.0

	for {
		points = append(points, Pos{x0, y0})
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}

		dist = p.Distance(Pos{x0, y0})
		if dist > maxDist {
			break
		}
	}

	return points
}

// ===== Directions =====

// Direction represents a cardinal direction, such as North, South, East, or West
type Direction int

const (
	DirNorth Direction = iota
	DirSouth
	DirEast
	DirWest
)

var Directions = []Direction{DirNorth, DirSouth, DirEast, DirWest}

func (d Direction) Pos() Pos {
	switch d {
	case DirNorth:
		return Pos{0, -1}
	case DirSouth:
		return Pos{0, 1}
	case DirEast:
		return Pos{1, 0}
	case DirWest:
		return Pos{-1, 0}
	}
	return Pos{}
}
