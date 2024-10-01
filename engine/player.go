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

	// Health
	currentHP int
	maxHP     int

	// Base attributes, can be modified by equipment
	defence      int //nolint Defence protects against damage
	attackDamage int //nolint Base damage added to weapon damage
	attackChance int //nolint Base chance to hit

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

func (p *Player) HP() int {
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

func (p *Player) BackpackSize() int {
	return playerMaxItems
}

func (p *Player) DropItem(item *Item) bool {
	if !p.backpack.Contains(item) {
		return false
	}

	if placedOK := p.currentTile.addItem(item); placedOK {
		p.backpack.Remove(item)
		item.dropped = true
		return true
	}

	return false
}

func (p *Player) Tile() *tile {
	return p.currentTile
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
	p.currentTile.items.Remove(item)
	item.pos = nil

	return true
}

func (p *Player) SetHP(hp int) {
	p.currentHP = hp
}

func (p *Player) SetMaxHP(hp int) {
	p.maxHP = hp
}
