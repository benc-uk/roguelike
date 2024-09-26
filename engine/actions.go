package engine

import (
	"roguelike/core"
)

type Action interface {
	Execute(g Game) bool
}

type MoveAction struct {
	direction
}

func NewMoveAction(d core.Direction) *MoveAction {
	return &MoveAction{d}
}

func (a *MoveAction) Execute(g Game) bool {
	p := g.Player()
	m := g.Map()

	destPos := p.pos.Add(a.Pos())
	destTile := m.TileAt(destPos)

	creatures := destTile.entities.AllCreatures()

	// Check if the player can move in the direction
	if !destPos.InBounds(m.Width, m.Height) || destTile.BlocksMove() {
		// Creature blocking the way
		if len(creatures) > 0 {
			creature := creatures[0]
			destTile.entities.Remove(creature)
			events.new(EventCreatureKilled, creature, "You killed a "+creature.name)
		}

		return false
	}

	p.pos.X += a.Pos().X
	p.pos.Y += a.Pos().Y

	// Check for items
	items := destTile.entities.AllItems()
	if len(items) == 1 {
		item := items[0]

		if p.pickupItem(item) {
			destTile.removeItem(item)
			events.new("item_pickup", item, "Picked up "+item.ShortDesc())
		} else {
			events.new("item_pickup_fail", item, "Inventory full")
		}
	} else if len(items) > 1 {
		events.new("item_pickup_multiple", nil, "You stand over a pile of items")
	}

	return true
}
