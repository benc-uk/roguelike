package engine

import (
	"fmt"
)

// ============================================================================
// Entities exist in the game world	- creatures, items, furniture etc
// This is the base entity type extended by other entity types
// ============================================================================

type entityType int

const (
	entityTypeCreature entityType = iota
	entityTypeItem
	entityTypeFurniture
)

type entityBase struct {
	*pos
	id         string
	instanceID string

	blocksMove bool
	blocksLOS  bool // nolint

	desc string
	name string

	graphicId string
	colour    string
}

type entity interface {
	Id() string
	InstanceID() string
	Graphic() string
	Description() string
	Type() entityType
	BlocksLOS() bool
	BlocksMove() bool
	setPos(*pos)
}

func (e entityBase) Id() string {
	return e.id
}

func (e entityBase) InstanceID() string {
	return e.instanceID
}

func (e entityBase) Description() string {
	return e.desc
}

func (e entityBase) Name() string {
	return e.name
}

func (e entityBase) Appearance() Appearance {
	return Appearance{
		Graphic: e.graphicId,
		Colour:  e.colour,
	}
}

func (e entityBase) BlocksLOS() bool {
	return e.blocksLOS
}

func (e entityBase) BlocksMove() bool {
	return e.blocksMove
}

func (e entityBase) Graphic() string {
	return e.graphicId
}

func (e entityBase) String() string {
	return fmt.Sprintf("entity_%s_%s", e.id, e.instanceID)
}

func (e *entityBase) setPos(p *pos) {
	e.pos = p
}

// ===== Lists ========================================================================================================

type entityList []entity

func NewEntityList() entityList {
	return make(entityList, 0)
}

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

func (el entityList) Count() int {
	return len(el)
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

func (el entityList) Last() entity {
	if len(el) == 0 {
		return nil
	}

	return el[len(el)-1]
}

func (el entityList) First() entity {
	if len(el) == 0 {
		return nil
	}

	return el[0]
}

func (el entityList) IsEmpty() bool {
	return len(el) == 0
}

func (el *entityList) Remove(e entity) {
	for i, ent := range *el {
		if ent.InstanceID() == e.InstanceID() {
			*el = append((*el)[:i], (*el)[i+1:]...)
			return
		}
	}
}

func (el *entityList) Add(e ...entity) {
	*el = append(*el, e...)
}

func (el entityList) Contains(thing entity) bool {
	for _, e := range el {
		if e.InstanceID() == thing.InstanceID() {
			return true
		}
	}

	return false
}
