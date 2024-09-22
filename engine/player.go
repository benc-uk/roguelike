package engine

import (
	"roguelike/core"
)

const (
	MAX_ITEMS = 10
)

type Player struct {
	core.Pos
	name string

	items []Item
}

func (p *Player) pickupItem(item *Item) bool {
	if len(p.items) >= MAX_ITEMS {
		return false
	}

	p.items = append(p.items, *item)
	return true
}
