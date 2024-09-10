package engine

import (
	"log"
	"roguelike/core"
)

type entityType int

const (
	entityTypeCreature entityType = iota
	entityTypeItem
	entityTypeFurniture
)

type entityBase struct {
	id string
	*core.Pos
	desc          string
	displayString string
	blocksMove    bool
	blocksLOS     bool // nolint
	hints         []string
}

type entity interface {
	Id() string
	Description() string
	Type() entityType
}

func (i *entityBase) Id() string {
	return i.id
}

func (i *entityBase) Description() string {
	return i.desc
}

func (i *entityBase) Appearance() Appearance {
	return Appearance{
		Details: i.displayString,
		Hints:   i.hints,
	}
}

// ===== Items ========================================================================================================

type Item struct {
	entityBase
	consumable bool
}

func (i *Item) Type() entityType {
	return entityTypeItem
}

type itemFactoryType map[string](func() *Item)

var itemFactory itemFactoryType

// TODO: This is placeholder code
func LoadItemFactory() {
	itemFactory = make(map[string](func() *Item))

	itemFactory["sword"] = func() *Item {
		return &Item{
			entityBase: entityBase{
				id:            "sword",
				desc:          "A rusty sword",
				displayString: "sword",
			},
			consumable: false,
		}
	}

	itemFactory["potion"] = func() *Item {
		return &Item{
			entityBase: entityBase{
				id:            "potion",
				desc:          "A refreshing looking blue potion",
				displayString: "potion",
			},
			consumable: true,
		}
	}

	itemFactory["potion_poison"] = func() *Item {
		return &Item{
			entityBase: entityBase{
				id:            "potion_poison",
				desc:          "A bubbling pink potion",
				displayString: "potion",
				hints:         []string{"colour::11"},
			},
			consumable: true,
		}
	}

	itemFactory["door"] = func() *Item {
		return &Item{
			entityBase: entityBase{
				id:            "door",
				desc:          "A sturdy wooden door",
				displayString: "door",
				blocksMove:    true,
			},
			consumable: false,
		}
	}

	itemFactory["rat"] = func() *Item {
		return &Item{
			entityBase: entityBase{
				id:            "rat",
				desc:          "A small, scurrying rat",
				displayString: "rat",
				blocksMove:    true,
			},
			consumable: false,
		}
	}
}

func (factory itemFactoryType) CreateItem(id string) *Item {
	factoryFunc, ok := factory[id]
	if !ok {
		log.Printf("ItemFactory: Item with id %s not found", id)
		return nil
	}

	// Create the item
	return factoryFunc()
}

// ===== Lists ========================================================================================================

type entityList []entity

func (el entityList) AllItems() []*Item {
	items := make([]*Item, 0)

	for _, e := range el {
		if e.Type() == entityTypeItem {
			i, ok := e.(*Item)
			if !ok {
				continue
			}
			items = append(items, i)
		}
	}

	return items
}

func (el entityList) AllCreatures() []*creature {
	creatures := make([]*creature, 0)

	for _, e := range el {
		if e.Type() == entityTypeCreature {
			c, ok := e.(*creature)
			if !ok {
				continue
			}
			creatures = append(creatures, c)
		}
	}

	return creatures
}

func (el entityList) Last() *entity {
	if len(el) == 0 {
		return nil
	}
	return &el[len(el)-1]
}

func (el entityList) First() *entity {
	if len(el) == 0 {
		return nil
	}

	return &el[0]
}

func (el entityList) IsEmpty() bool {
	return len(el) == 0
}

// ===== Creatures ======================================================================================================

type creature struct {
	entityBase
}

func (m *creature) Type() entityType {
	return entityTypeCreature
}
