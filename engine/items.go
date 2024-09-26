package engine

import (
	"roguelike/core"
	"slices"

	"gopkg.in/yaml.v3"
)

// ============================================================================
// Item entities are things like weapons, armour, potions etc
// ============================================================================

type Item struct {
	entityBase
	usable     bool
	equippable bool // nolint
}

func (i *Item) Type() entityType {
	return entityTypeItem
}

func (i *Item) String() string {
	return "item_" + i.id + "_" + i.instanceID
}

// ===== Item Generator =================================================================================================

type itemGenerator struct {
	genFunctions map[string](func() *Item)
	keys         []string
}

type yamlItem struct {
	Description string `yaml:"description"`
	Name        string `yaml:"name"`
	Graphic     string `yaml:"graphic"`
	Colour      string `yaml:"colour"`
	Usable      bool   `yaml:"usable"`
}

type yamlItemsFile struct {
	Items map[string]yamlItem `yaml:"items"`
}

func newItemGenerator(dataFile string) (*itemGenerator, error) {
	data, err := core.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}

	var itemsFile yamlItemsFile
	err = yaml.Unmarshal(data, &itemsFile)
	if err != nil {
		return nil, err
	}

	gen := itemGenerator{
		genFunctions: make(map[string](func() *Item)),
		keys:         make([]string, 0),
	}

	for id, item := range itemsFile.Items {
		gen.genFunctions[id] = func() *Item {
			return &Item{
				entityBase: entityBase{
					id:         id,
					instanceID: core.RandId(6),
					desc:       item.Description,
					name:       item.Name,
					graphicId:  item.Graphic,
					colour:     item.Colour,
				},
				usable: item.Usable,
			}
		}

		gen.keys = append(gen.keys, id)
	}

	// Sort the keys as the map iteration order above is random
	slices.Sort(gen.keys)

	return &gen, nil
}

// nolint
func (gen itemGenerator) createItem(id string) *Item {
	itemFunc, ok := gen.genFunctions[id]
	if !ok {
		return nil
	}

	// Create the item by invoking the generation function
	return itemFunc()
}

func (gen itemGenerator) createRandomItem() *Item {
	if len(gen.genFunctions) == 0 {
		return nil
	}

	// Get a random item, we have to use the keys slice as the map iteration order is random
	id := gen.keys[rng.IntN(len(gen.keys))]
	return gen.createItem(id)
}
