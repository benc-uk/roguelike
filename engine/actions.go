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
	target *creature //nolint
}

func NewMoveAction(d core.Direction) *MoveAction {
	return &MoveAction{d}
}

func (a *MoveAction) Execute(g Game) ActionResult {
	p := g.Player()
	m := g.Map()

	destPos := p.pos.Add(a.Pos())
	destTile := m.TileAt(destPos)
	if destTile == nil {
		return ActionResult{false, 0}
	}

	creatures := destTile.entities.AllCreatures()

	// TODO: Energy not really implemented
	energy := 3

	// Check if the player can move in the direction
	if !destPos.InBounds(m.Width, m.Height) || destTile.BlocksMove() {
		// Creature blocking the way
		if len(creatures) > 0 {
			creature := creatures[0]
			destTile.entities.Remove(creature)
			message := fmt.Sprintf("You %s a %s",
				randString("killed", "defeated", "felled", "vanquished", "slayed", "destroyed", "murdered"),
				creature.Name())
			events.new("creature_killed", creature, message)
			p.exp += creature.xp
			energy = 60
			return ActionResult{true, energy}
		}

		return ActionResult{false, energy}
	}

	p.moveToTile(destTile)

	// Check for items and auto pick them up
	items := destTile.entities.AllItems()
	if len(items) == 1 {
		item := items[0]
		if item.dropped {
			return ActionResult{false, 0}
		}

		if p.PickupItem(item) {
			events.new(EventItemPickup, item, "Picked up "+item.Name())
			energy = 40
		} else {
			events.new("item_pickup_fail", item, "You are carrying too much!")
		}
	} else if len(items) > 1 {
		events.new(EventItemMultiple, nil, "You stand over a pile of items")
	}

	return ActionResult{true, energy}
}
