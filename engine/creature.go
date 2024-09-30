package engine

// ============================================================================
// Creature entities are mostly monsters but can be also non-hostile
// ============================================================================

import (
	"fmt"
	"roguelike/core"
	"slices"

	"gopkg.in/yaml.v3"
)

type creature struct {
	entityBase

	currentTile *tile
	hp          int //nolint
	xp          int //nolint
}

func (c creature) String() string {
	return fmt.Sprintf("creature_%v_%v", c.id, c.instanceID)
}

func (c creature) Type() entityType {
	return entityTypeCreature
}

func (c creature) BlocksLOS() bool {
	return false
}

func (c creature) BlocksMove() bool {
	return true
}

// ===== Creature Generator =================================================================================================

type creatureGenerator struct {
	genFunctions map[string](func() *creature)
	keys         []string
}

type yamlCreature struct {
	// Entity fields
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Graphic     string `yaml:"graphic"`
	Colour      string `yaml:"colour"`

	// Creature specific fields
	Hostile bool `yaml:"hostile"`
	Hp      int  `yaml:"hp"`
	Xp      int  `yaml:"xp"`
}

type yamlCreaturesFile struct {
	Creatures map[string]yamlCreature `yaml:"creatures"`
}

func newCreatureGenerator(dataFile string) (*creatureGenerator, error) {
	data, err := core.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}

	var file yamlCreaturesFile
	err = yaml.Unmarshal(data, &file)
	if err != nil {
		return nil, err
	}

	gen := creatureGenerator{
		genFunctions: make(map[string](func() *creature)),
		keys:         make([]string, 0),
	}

	for id, creat := range file.Creatures {
		gen.genFunctions[id] = func() *creature {
			return &creature{
				entityBase: entityBase{
					id:         id,
					instanceID: core.RandId(6),
					desc:       creat.Description,
					name:       creat.Name,
					graphicId:  creat.Graphic,
					colour:     creat.Colour,
				},
				hp: creat.Hp,
				xp: creat.Xp,
			}
		}

		gen.keys = append(gen.keys, id)
	}

	// Sort the keys as the map iteration order above is random
	slices.Sort(gen.keys)
	return &gen, nil
}

func (gen creatureGenerator) createCreature(id string) *creature {
	genFunc, ok := gen.genFunctions[id]
	if !ok {
		return nil
	}

	// Create the creature by invoking the generation function
	return genFunc()
}

func (gen creatureGenerator) createRandomCreature() *creature {
	if len(gen.genFunctions) == 0 {
		return nil
	}

	// Get a random creature, we have to use the keys slice as the map iteration order is random
	id := gen.keys[rng.IntN(len(gen.keys))]
	return gen.createCreature(id)
}
