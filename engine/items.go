package engine

import (
	"roguelike/core"

	"gopkg.in/yaml.v3"
)

// ===== Items ========================================================================================================

type Item struct {
	entityBase
	usable     bool
	equippable bool //nolint:unused
}

func (i *Item) Type() entityType {
	return entityTypeItem
}

func (i *Item) String() string {
	return "item_" + i.id + "_" + i.instanceID
}

// ===== Item Factory =================================================================================================

type itemFactoryDB map[string](func() *Item)

type yamlItem struct {
	Description string `yaml:"description"`
	Short       string `yaml:"short"`
	Graphic     string `yaml:"graphic"`
	Colour      string `yaml:"colour"`
	Usable      bool   `yaml:"usable"`
}

type yamlItemsFile struct {
	Items map[string]yamlItem `yaml:"items"`
}

func newItemFactory(dataFile string) (itemFactoryDB, error) {
	// Load items from YAML data file
	data, err := core.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}

	var itemsFile yamlItemsFile
	err = yaml.Unmarshal(data, &itemsFile)
	if err != nil {
		return nil, err
	}

	itemFactory := make(map[string](func() *Item))
	for id, item := range itemsFile.Items {
		itemFactory[id] = func() *Item {
			return &Item{
				entityBase: entityBase{
					id:         id,
					instanceID: core.RandId(6),
					desc:       item.Description,
					shortDesc:  item.Short,
					graphicId:  item.Graphic,
					colour:     item.Colour,
				},
				usable: item.Usable,
			}
		}
	}

	return itemFactory, nil
}

// nolint
func (factory itemFactoryDB) createItem(id string) *Item {
	itemFunc, ok := factory[id]
	if !ok {
		return nil
	}

	// Create the item
	return itemFunc()
}

func (factory itemFactoryDB) createRandomItem() *Item {
	if len(factory) == 0 {
		return nil
	}

	// Iteration in go is random, so we can just grab the first item
	for _, itemFunc := range factory {
		return itemFunc()
	}

	return nil
}
