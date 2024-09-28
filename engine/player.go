package engine

import (
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

	items []Item
}

func NewPlayer(pos core.Pos) *Player {
	name := "Jimmy No Name"
	gen, err := fn.Compile("sd", fn.Collapse(true), fn.RandFn(rng.IntN))
	if err == nil {
		name = gen.String()
		// Capitalize the first letter
		name = strings.ToUpper(name[:1]) + name[1:]
	}

	return &Player{
		pos:       pos,
		name:      name,
		currentHP: 10,
		maxHP:     50,
		exp:       0,
		level:     1,
	}
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

func (p *Player) Inventory() []Item {
	return p.items
}

func (p *Player) MaxItems() int {
	return playerMaxItems
}

func (p *Player) DropItem(index int) {
	if index < 0 || index >= len(p.items) {
		return
	}

	item := p.items[index]
	p.items = append(p.items[:index], p.items[index+1:]...)
	item.pos = &p.pos
	p.currentTile.placeItem(&item)
}

func (p *Player) moveToTile(t *tile) {
	p.pos = t.pos
	p.currentTile = t
}

// Pickup an item from the ground, returning true if the item was picked up
func (p *Player) pickupItem(item *Item) bool {
	if len(p.items) >= playerMaxItems {
		return false
	}

	p.items = append(p.items, *item)
	item.pos = nil

	p.currentTile.removeItem(item)
	return true
}
