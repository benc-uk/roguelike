package engine

import (
	"roguelike/core"
)

// ============================================================================
// Generator for the game map using a simple BSP tree
// ============================================================================

// Representation of a BSP tree node, used in the generation of the map
type bspNode struct {
	core.Rect
	center core.Pos
	depth  int

	// Child nodes
	Left, Right *bspNode
}

// BSP generator
type bspGenerator struct {
	// Maximum depth of the BSP tree, majorly affects the number of rooms
	maxDepth int

	// Holds the game map to generate into
	gameMap *GameMap
}

func newBSPGenerator(maxDepth int, gameMap *GameMap) *bspGenerator {
	return &bspGenerator{maxDepth, gameMap}
}

func (b bspNode) isLeaf() bool {
	return b.Left == nil && b.Right == nil
}

func (gen *bspGenerator) generate() {
	// Create a BSP tree
	root := gen.buildBSP(rect{pos{0, 0}, gen.gameMap.Size()}, 0) //nolint

	// Traverse the tree which creates rooms at leaf nodes
	root.traverseBSP(gen.gameMap)

	// Separate pass to create corridors
	root.createCorridors(gen.gameMap)

	gen.gameMap.description += " dungeon"
	gen.gameMap.generationMethod = "bsp"
}

func (gen bspGenerator) buildBSP(r core.Rect, depth int) *bspNode {
	// If we've reached the maximum depth, return a leaf node and end the recursion
	if depth >= gen.maxDepth {
		return &bspNode{r, r.Center(), depth, nil, nil}
	}

	aspect := float64(r.Width / r.Height)
	var horiz bool
	if aspect >= 1.25 {
		horiz = false
	} else if aspect <= 0.75 {
		horiz = true
	} else {
		horiz = rng.IntN(100) < 50
	}

	if horiz {
		// Split horizontally
		split := r.Height / 2

		// Move the split point by random factor
		factor := 0.4
		splitRand := int(float64(split) * factor * 2)
		if splitRand <= 0 {
			splitRand = 1
		}
		split += rng.IntN(splitRand) - int(float64(split)*factor)

		left := core.NewRect(r.X, r.Y, r.Width, split)
		right := core.NewRect(r.X, r.Y+split, r.Width, r.Height-split)
		return &bspNode{r, r.Center(), depth, gen.buildBSP(left, depth+1), gen.buildBSP(right, depth+1)}
	} else {
		// Split vertically
		split := r.Width / 2

		// Move the split point by random factor
		factor := 0.4
		splitRand := int(float64(split) * factor * 2)
		if splitRand <= 0 {
			splitRand = 1
		}
		split += rng.IntN(splitRand) - int(float64(split)*factor)

		left := core.NewRect(r.X, r.Y, split, r.Height)
		right := core.NewRect(r.X+split, r.Y, r.Width-split, r.Height)
		return &bspNode{r, r.Center(), depth, gen.buildBSP(left, depth+1), gen.buildBSP(right, depth+1)}
	}
}

func (node *bspNode) traverseBSP(gm *GameMap) {
	if node.Left != nil {
		node.Left.traverseBSP(gm)
	}
	if node.Right != nil {
		node.Right.traverseBSP(gm)
	}

	// Percentage chance to create a room
	if rng.IntN(100) < 70 {
		// Create a room at the center of leaf nodes
		if node.Left == nil && node.Right == nil {
			// Room size is randomly % of the node size
			width := node.Width * (40 + rng.IntN(50)) / 100
			height := node.Height * (40 + rng.IntN(50)) / 100
			room := core.NewRect(node.center.X-width/2, node.center.Y-height/2, width, height)

			// Carve the room area
			gm.setAreaRect(false, room)
		}
	}
}

func (node *bspNode) createCorridors(gm *GameMap) {
	if !node.isLeaf() {
		var corridor core.Rect
		a := node.Left
		b := node.Right

		if a.center.X == b.center.X {
			// Vertical corridor
			corridor = core.NewRect(a.center.X-1, a.center.Y, 1, b.center.Y-a.center.Y)
		} else {
			// Horizontal corridor
			corridor = core.NewRect(a.center.X, a.center.Y-1, b.center.X-a.center.X, 1)
		}

		// Carve the corridor
		gm.setAreaRect(false, corridor)

		// Recurse
		node.Left.createCorridors(gm)
		node.Right.createCorridors(gm)
	}
}
