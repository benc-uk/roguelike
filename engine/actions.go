package engine

// ============================================================================
// Actions & action execution operate on the game state
// ============================================================================

import (
	"fmt"
	"roguelike/core"
)

type Action interface {
	Execute(g Game) ActionResult
}

type ActionResult struct {
	Success     bool
	EnergySpent int
}

type MoveAction struct {
	direction
}

type AttackAction struct {
	target *creature
}

func NewMoveAction(d core.Direction) *MoveAction {
	return &MoveAction{d}
}

func NewAttackAction(target *creature) *AttackAction {
	return &AttackAction{target}
}

func (a *MoveAction) Execute(g Game) ActionResult {
	p := g.Player()
	m := g.Map()

	destTile := p.currentTile.AdjacentTileDir(a.direction, m)

	if destTile == nil || destTile.BlocksMove() {
		return ActionResult{false, 0}
	}

	energy := 4
	p.moveToTile(destTile)

	// Check for items and auto pick them up
	items := destTile.items
	if len(items) == 1 {
		item, isItem := items[0].(*Item)
		if !isItem {
			return ActionResult{false, 0}
		}

		if item.dropped {
			events.new("item_pickup_dropped", item, fmt.Sprintf("You see a %s you previously dropped", item.Name()))
			return ActionResult{false, 0}
		}

		if p.PickupItem(item) {
			events.new(EventItemPickup, item, "Picked up "+item.Name())
			energy = 40
		} else {
			events.new("item_pickup_fail", item, "You are carrying too much!")
		}
	} else if len(items) > 1 {
		events.new(EventItemMultiple, nil, fmt.Sprintf("You stand over a pile of %d items", len(items)))
	}

	return ActionResult{true, energy}
}

func (a *AttackAction) Execute(g Game) ActionResult {
	p := g.Player()

	// Check adjacent
	if !p.pos.IsNeighbour(*a.target.pos) {
		return ActionResult{false, 0}
	}

	// Attack the target
	// TODO: Add a combat system here :)
	a.target.currentTile.creature = nil
	a.target.currentTile = nil
	message := fmt.Sprintf("You %s a %s",
		randString("killed", "defeated", "felled", "vanquished", "slayed", "destroyed", "murdered"),
		a.target.Name())
	events.new("creature_killed", a.target, message)
	p.exp += a.target.xp

	return ActionResult{true, 60}
}
