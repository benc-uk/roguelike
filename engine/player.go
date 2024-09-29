package engine

import (
	"fmt"
	"roguelike/core"
	"strings"

	fn "github.com/s0rg/fantasyname"
)

// ============================================================================
// Player data structure holds the player's state
// ============================================================================

const playerMaxItems = 10

type Player struct {
	pos
	currentTile *tile
	name        string

	currentHP int
	maxHP     int

	exp   int
	level int

	backpack entityList

	// Inspired by Angband https://angband.readthedocs.io/en/latest/command.html#inventory-commands
	equipSlots map[equipLocation]*Item
}

func NewPlayer(tile *tile, items ...Item) *Player {
	name := "Jimmy No Name"
	gen, err := fn.Compile("sd", fn.Collapse(true), fn.RandFn(rng.IntN))
	if err == nil {
		name = gen.String()
		// Capitalize the first letter
		name = strings.ToUpper(name[:1]) + name[1:]
	}

	p := &Player{
		pos:         tile.pos,
		currentTile: tile,
		name:        name,
		currentHP:   10,
		maxHP:       50,
		exp:         0,
		level:       1,
		backpack:    NewEntityList(),
		equipSlots:  make(map[equipLocation]*Item),
	}

	return p
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) CurrentHP() int {
	return p.currentHP
}

func (p *Player) MaxHP() int {
	return p.maxHP
}

func (p *Player) Exp() int {
	return p.exp
}

func (p *Player) Level() int {
	return p.level
}

func (p *Player) Pos() core.Pos {
	return p.pos
}

func (p *Player) Inventory() []*Item {
	return p.backpack.AllItems()
}

func (p *Player) MaxItems() int {
	return playerMaxItems
}

func (p *Player) DropItem(item *Item) {
	if !p.backpack.Contains(item) {
		return
	}

	if placedOK := p.currentTile.addEntity(item); placedOK {
		p.backpack.Remove(item)
		events.new(EventItemDropped, item, fmt.Sprintf("You dropped the %s", item.Name()))
		item.dropped = true
	}
}

func (p *Player) CurrentTile() tile {
	return *p.currentTile
}

func (p *Player) moveToTile(t *tile) {
	p.pos = t.pos
	p.currentTile = t
}

// Pickup an item from the ground, returning true if the item was picked up
func (p *Player) PickupItem(item *Item) bool {
	if len(p.backpack) >= playerMaxItems {
		return false
	}

	p.backpack.Add(item)
	p.currentTile.entities.Remove(item)
	item.pos = nil

	return true
}
