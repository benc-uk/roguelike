package engine

import (
	"fmt"
	"math/rand"
	"roguelike/core"
)

// Representation of a BSP tree node, used in the generation of the map
type bspNode struct {
	core.Rect
	center core.Pos
	depth  int

	// Child nodes
	Left, Right *bspNode
}

// BSP generator only has one method, to build the BSP tree
type bspGenerator struct {
	// Maximum depth of the BSP tree, majorly affects the number of rooms
	maxDepth int
}

func (b bspNode) isLeaf() bool {
	return b.Left == nil && b.Right == nil
}

func (b bspNode) String() string {
	return fmt.Sprintf("BSPNode{Rect: %v, Depth: %d, Center: %v, Left: %v, Right: %v}", b.Rect, b.depth, b.center, b.Left, b.Right)
}

func (gm *GameMap) GenerateBSP(maxGenDepth int, itemFactory itemFactoryDB) {
	// Create a BSP tree
	gen := bspGenerator{maxGenDepth}
	root := gen.buildBSP(gm.Rect(), 0)

	// Traverse the tree and create rooms
	root.traverseBSP(gm, itemFactory)

	// Separate pass to create corridors
	root.createCorridors(gm)
}

func (node *bspNode) traverseBSP(gm *GameMap, itemFactory itemFactoryDB) {
	if node.Left != nil {
		node.Left.traverseBSP(gm, itemFactory)
	}
	if node.Right != nil {
		node.Right.traverseBSP(gm, itemFactory)
	}

	// Percentage chance to create a room
	if rand.Intn(100) < 70 {
		// Create a room at the center of leaf nodes
		if node.Left == nil && node.Right == nil {
			// Room size is randomly 40-70% of the node size
			width := node.Width * (40 + rand.Intn(30)) / 100
			height := node.Height * (40 + rand.Intn(30)) / 100
			room := core.NewRect(node.center.X-width/2, node.center.Y-height/2, width, height)

			// Carve the room area
			gm.floorAreaRect(room)

			// Place 0-3 items in the room
			numItems := rand.Intn(4)
			for i := 0; i < numItems; i++ {
				pos := room.RandomPos()
				item := itemFactory.createRandomItem()
				gm.TileAt(pos).placeItem(item)
			}
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
		gm.floorAreaRect(corridor)

		// Recurse
		node.Left.createCorridors(gm)
		node.Right.createCorridors(gm)
	}
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
		horiz = rand.Intn(100) < 50
	}

	if horiz {
		// Split horizontally
		split := r.Height / 2
		// move the split point by a factor of 0.4
		factor := 0.4
		split += rand.Intn(int(float64(split)*factor*2)) - int(float64(split)*factor)

		left := core.NewRect(r.X, r.Y, r.Width, split)
		right := core.NewRect(r.X, r.Y+split, r.Width, r.Height-split)
		return &bspNode{r, r.Center(), depth, gen.buildBSP(left, depth+1), gen.buildBSP(right, depth+1)}
	} else {
		// Split vertically
		split := r.Width / 2
		// move the split point by a factor of 0.4
		factor := 0.4
		split += rand.Intn(int(float64(split)*factor*2)) - int(float64(split)*factor)

		left := core.NewRect(r.X, r.Y, split, r.Height)
		right := core.NewRect(r.X+split, r.Y, r.Width-split, r.Height)
		return &bspNode{r, r.Center(), depth, gen.buildBSP(left, depth+1), gen.buildBSP(right, depth+1)}
	}
}
