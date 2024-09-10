package core

import "math/rand"

type Size struct {
	Width  int
	Height int
}

type Pos struct {
	X int
	Y int
}

type Direction int

const (
	North Direction = iota
	South
	East
	West
)

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

func (p Pos) Neighbours() []Pos {
	return []Pos{
		{p.X - 1, p.Y},
		{p.X + 1, p.Y},
		{p.X, p.Y - 1},
		{p.X, p.Y + 1},
	}
}
