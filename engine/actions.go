package engine

import (
	"roguelike/core"
)

type Action interface {
	Execute() bool
}

type MoveAction struct {
	Direction core.Direction
}

func NewMoveAction(d core.Direction) *MoveAction {
	return &MoveAction{Direction: d}
}

func (a *MoveAction) Execute(p *Player, m *GameMap) bool {
	// heck if the player can move in the direction
	newPos := p.Pos.Add(a.Direction.Pos())
	if !newPos.InBounds(m.Width, m.Height) {
		return false
	}

	tile := m.TileAt(newPos)
	if tile.BlocksMove() {
		return false
	}

	p.X += a.Direction.Pos().X
	p.Y += a.Direction.Pos().Y

	// Check for items
	items := tile.entities.AllItems()
	if len(items) == 1 {
		item := items[0]

		if p.pickupItem(item) {
			tile.removeItem(item)
			events.new("item_pickup", item, "Picked up "+item.ShortDesc())
		} else {
			events.new("item_pickup_fail", item, "Inventory full")
		}
	} else if len(items) > 1 {
		events.new("item_pickup_multiple", nil, "You stand over a pile of items")
	}

	return true
}
