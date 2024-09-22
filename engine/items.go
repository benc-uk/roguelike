package engine

import (
	"roguelike/core"

	"math/rand"

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

func (factory itemFactoryDB) CreateItem(id string) *Item {
	factoryFunc, ok := factory[id]
	if !ok {
		return nil
	}

	// Create the item
	return factoryFunc()
}

func (i *Item) String() string {
	return "item_" + i.id + "_" + i.instanceID
}

// ===== Item Factory =================================================================================================

type itemFactoryDB map[string](func() *Item)

type yamlItem struct {
	Description string `yaml:"description"`
	Graphic     string `yaml:"graphic"`
	Colour      string `yaml:"colour"`
	Usable      bool   `yaml:"usable"`
}

type yamlItemsFile struct {
	Items map[string]yamlItem `yaml:"items"`
}

func NewItemFactory(dataFile string) (itemFactoryDB, error) {
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
					instanceID: randId(),
					desc:       item.Description,
					graphicId:  item.Graphic,
					colour:     item.Colour,
				},
				usable: item.Usable,
			}
		}
	}

	return itemFactory, nil
}

// Simple ID generator 8 characters long
func randId() string {
	// generate a random string
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, 8)
	for i := range id {
		id[i] = chars[rand.Intn(len(chars))]
	}
	return string(id)
}
