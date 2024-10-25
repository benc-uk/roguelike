package engine

// ============================================================================
// Generation used to create the game world using several techniques
// ============================================================================

type Generator interface {
	generate()
}

func generateMap(g *Game, dataFileDir string) {
	var err error
	g.itemGen, err = newItemGenerator(dataFileDir + "/items.yaml")
	if err != nil {
		panic(err)
	}

	g.creatureGen, err = newCreatureGenerator(dataFileDir + "/creatures.yaml")
	if err != nil {
		panic(err)
	}

	size := rng.IntN(4)
	depth := 1
	genDepth := 1
	switch size {
	case 0:
		// Tiny
		g.gameMap = NewMap(32, 32, depth)
		genDepth = rng.IntN(3) + 3 // 3,4,5
		g.gameMap.description = "a tiny"

	case 1:
		// Small
		g.gameMap = NewMap(40, 40, depth)
		genDepth = rng.IntN(3) + 3 // 3,4,5
		g.gameMap.description = "a small"

	case 2:
		// Medium
		g.gameMap = NewMap(48, 48, depth)
		genDepth = rng.IntN(3) + 4 // 4,5,6
		g.gameMap.description = "a fair sized"

	case 3:
		// Large
		g.gameMap = NewMap(64, 64, depth)
		genDepth = rng.IntN(3) + 5 // 5,6,7
		g.gameMap.description = "a large"
	}

	var gen Generator
	gen = newBSPGenerator(genDepth, g.gameMap)

	// 25% chance of a cave map
	if rng.IntN(4) == 0 {
		gen = newCaGenerator(g.gameMap)
	}

	gen.generate()

	// Place items
	den := (size + 1) * 3
	numItems := rng.IntN(den) + den
	for i := 0; i < numItems; i++ {
		item := g.itemGen.createRandomItem(rarityCommon)
		g.gameMap.randomFloorTile(false).addItem(item)
	}

	// Place creatures
	numCreatures := rng.IntN(den) + den
	for i := 0; i < numCreatures; i++ {
		creature := g.creatureGen.createRandomCreature()
		g.gameMap.randomFloorTile(false).placeCreature(creature)
	}
}
