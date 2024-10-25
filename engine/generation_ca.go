package engine

import (
	"math/rand/v2"
)

// ============================================================================
// Generator for the game map using celluar automata, to make caves
// ============================================================================

// https://code.tutsplus.com/generate-random-cave-levels-using-cellular-automata--gamedev-9664t

const startChance = 0.43
const deathLimit = 3
const birthLimit = 4
const iterations = 15

// Celluar automata generator
type caGenerator struct {
	gameMap *GameMap
	cells   [][]bool
	width   int
	height  int
}

func newCaGenerator(gameMap *GameMap) *caGenerator {
	cells := make([][]bool, gameMap.Width)
	for i := 0; i < gameMap.Width; i++ {
		cells[i] = make([]bool, gameMap.Height)
	}

	return &caGenerator{
		gameMap: gameMap,
		cells:   cells,
		width:   gameMap.Width,
		height:  gameMap.Height,
	}
}

func (gen *caGenerator) generate() {
	// Start with a random map
	for x := 0; x < gen.width; x++ {
		for y := 0; y < gen.height; y++ {
			if rand.Float64() < startChance {
				gen.cells[x][y] = true
			} else {
				gen.cells[x][y] = false
			}
		}
	}

	// Run the simulation
	for i := 0; i < iterations; i++ {
		gen.doSimulationStep()
	}

	// Convert cells to tiles in the game map
	for x := 0; x < gen.width; x++ {
		for y := 0; y < gen.height; y++ {
			if gen.cells[x][y] {
				gen.gameMap.Tile(x, y).makeWall()
			} else {
				gen.gameMap.Tile(x, y).makeFloor()
			}
		}
	}

	gen.gameMap.description += " cave"
	gen.gameMap.generationMethod = "cellular automata"
}

func (gen *caGenerator) doSimulationStep() {
	newCells := make([][]bool, gen.width)
	for i := 0; i < gen.height; i++ {
		newCells[i] = make([]bool, gen.height)
	}

	for x := 0; x < gen.width; x++ {
		for y := 0; y < gen.height; y++ {
			aliveNeighbours := gen.countAliveNeighbours(x, y)
			if gen.cells[x][y] {
				if aliveNeighbours < deathLimit {
					newCells[x][y] = false
				} else {
					newCells[x][y] = true
				}
			} else {
				if aliveNeighbours > birthLimit {
					newCells[x][y] = true
				} else {
					newCells[x][y] = false
				}
			}
		}
	}

	gen.cells = newCells
}

func (gen *caGenerator) countAliveNeighbours(x, y int) int {
	count := 0

	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			neighbourX := x + i
			neighbourY := y + j

			if i == 0 && j == 0 {
				continue
			} else if neighbourX < 0 || neighbourY < 0 || neighbourX >= gen.width || neighbourY >= gen.height {
				count++
			} else if gen.cells[neighbourX][neighbourY] {
				count++
			}
		}
	}

	return count
}
