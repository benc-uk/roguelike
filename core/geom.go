package core

// ====================================================================================================================
// General 2D geometry and direction utilities, such as positions, rectangles, and directions
// ====================================================================================================================

import (
	"fmt"
	"math/rand"
)

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

func NewRect(x, y, width, height int) Rect {
	return Rect{Pos{x, y}, Size{width, height}}
}

// ===== Positions =====

type Pos struct {
	X int
	Y int
}

func (p Pos) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func RandomPos(width, height int) Pos {
	x := rand.Intn(width)
	y := rand.Intn(height)
	return Pos{x, y}
}

func (p Pos) Add(p2 Pos) Pos {
	return Pos{p.X + p2.X, p.Y + p2.Y}
}

func (p Pos) Sub(p2 Pos) Pos {
	return Pos{p.X - p2.X, p.Y - p2.Y}
}

func (p Pos) Distance(p2 Pos) int {
	return AbsInt(p.X-p2.X) + AbsInt(p.Y-p2.Y)
}

func (p Pos) InBounds(width, height int) bool {
	return p.X >= 0 && p.X < width && p.Y >= 0 && p.Y < height
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

// ===== Directions =====

type Direction int

const (
	North Direction = iota
	South
	East
	West
)

var Directions = []Direction{North, South, East, West}

func (d Direction) Pos() Pos {
	switch d {
	case North:
		return Pos{0, -1}
	case South:
		return Pos{0, 1}
	case East:
		return Pos{1, 0}
	case West:
		return Pos{-1, 0}
	}
	return Pos{}
}
