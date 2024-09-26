package engine

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

func NewMoveAction(d core.Direction) *MoveAction {
	return &MoveAction{d}
}

func (a *MoveAction) Execute(g Game) ActionResult {
	p := g.Player()
	m := g.Map()

	destPos := p.pos.Add(a.Pos())
	destTile := m.TileAt(destPos)

	creatures := destTile.entities.AllCreatures()

	// TODO: Not implemented
	energy := 5

	// Check if the player can move in the direction
	if !destPos.InBounds(m.Width, m.Height) || destTile.BlocksMove() {
		// Creature blocking the way
		if len(creatures) > 0 {
			creature := creatures[0]
			destTile.entities.Remove(creature)
			message := fmt.Sprintf("You %s a %s",
				randString("killed", "defeated", "felled", "vanquished", "slew", "destroyed", "murdered"),
				creature.Name())
			events.new("creature_killed", creature, message)
			energy = 40
			return ActionResult{true, energy}
		}

		return ActionResult{false, energy}
	}

	p.pos.X += a.Pos().X
	p.pos.Y += a.Pos().Y

	// Check for items
	items := destTile.entities.AllItems()
	if len(items) == 1 {
		item := items[0]

		if p.pickupItem(item) {
			destTile.removeItem(item)
			events.new("item_pickup", item, "Picked up "+item.Name())
			energy = 40
		} else {
			events.new("item_pickup_fail", item, "Inventory full")
		}
	} else if len(items) > 1 {
		events.new("item_pickup_multiple", nil, "You stand over a pile of items")
	}

	return ActionResult{true, energy}
}
