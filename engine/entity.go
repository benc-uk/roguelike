package engine

import (
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
	blocksMove bool
	blocksLOS  bool // nolint

	desc string

	graphicId string
	colour    string
}

type entity interface {
	Id() string
	Description() string
	Type() entityType
	BlocksLOS() bool
	BlocksMove() bool
}

func (e *entityBase) Id() string {
	return e.id
}

func (e *entityBase) Description() string {
	return e.desc
}

func (e *entityBase) Appearance() Appearance {
	return Appearance{
		Graphic: e.graphicId,
		Colour:  e.colour,
	}
}

func (e *entityBase) BlocksLOS() bool {
	return e.blocksLOS
}

func (e *entityBase) BlocksMove() bool {
	return e.blocksMove
}

// ===== Items ========================================================================================================

type Item struct {
	entityBase
	consumable bool
	equippable bool //nolint
}

func (i *Item) Type() entityType {
	return entityTypeItem
}

type itemFactoryDB map[string](func() *Item)

var itemFactory itemFactoryDB

func (factory itemFactoryDB) CreateItem(id string) *Item {
	factoryFunc, ok := factory[id]
	if !ok {
		return nil
	}

	// Create the item
	return factoryFunc()
}

// ===== Furniture ========================================================================================================

type Furniture struct {
	entityBase
}

func (f *Furniture) Type() entityType {
	return entityTypeFurniture
}

func (f *Furniture) BlocksLOS() bool {
	return true
}

func (f *Furniture) BlocksMove() bool {
	return true
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
