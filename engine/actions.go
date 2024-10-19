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

type PickupAction struct {
	item *Item
}

type DropAction struct {
	item *Item
}

type UseAction struct {
	item *Item
}

type EquipAction struct {
	item *Item
}

func NewMoveAction(d core.Direction) *MoveAction {
	return &MoveAction{d}
}

func NewAttackAction(target *creature) *AttackAction {
	return &AttackAction{target}
}

func NewPickupAction(item *Item) *PickupAction {
	return &PickupAction{item}
}

func NewDropAction(item *Item) *DropAction {
	return &DropAction{item}
}

func NewUseAction(item *Item) *UseAction {
	return &UseAction{item}
}

func NewEquipAction(item *Item) *EquipAction {
	return &EquipAction{item}
}

func (a *MoveAction) Execute(g Game) ActionResult {
	p := g.Player()
	m := g.Map()

	destTile := m.AdjacentTile(p.currentTile, a.direction)

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
			return ActionResult{true, energy}
		}

		if item.dropped {
			events.new(EventItemSkipped, item, fmt.Sprintf("You see a %s you previously dropped", item.Name()))
			return ActionResult{true, 40}
		}

		pickupAction := NewPickupAction(item)
		return pickupAction.Execute(g)
	} else if len(items) > 1 {
		events.new(EventItemMultiple, nil, fmt.Sprintf("You stand over a pile of %d items", len(items)))
	}

	g.updateFOV()
	return ActionResult{true, energy}
}

func (a *AttackAction) Execute(g Game) ActionResult {
	p := g.Player()

	// Check adjacent
	// TODO: Maybe remove this for ranged attacks
	if !p.pos.IsNeighbour(*a.target.pos) {
		return ActionResult{false, 0}
	}

	// Attack the target
	// TODO: Add a combat system here :)
	hitRoll := d100.Roll()
	events.new(EventMiscMessage, a.target, fmt.Sprintf("Hit roll: %d need < %d", hitRoll, p.attackChance))
	if hitRoll > p.attackChance {
		events.new(EventCombatMissed, a.target, "You missed!")
		return ActionResult{false, 0}
	}

	rollResult := p.attackRoll.Roll() + p.attackDamage
	events.new(EventCombatHit, nil, fmt.Sprintf("Attack: %s+%d → %d", p.attackRoll, p.attackDamage, rollResult))

	a.target.currentTile.creature = nil
	a.target.currentTile = nil
	message := fmt.Sprintf("You %s a %s",
		randString("killed", "defeated", "felled", "vanquished", "slayed", "destroyed", "murdered"),
		a.target.Name())
	events.new(EventCreatureKilled, a.target, message)
	p.exp += a.target.xp

	return ActionResult{true, 60}
}

func (a *PickupAction) Execute(g Game) ActionResult {
	p := g.Player()

	if p.PickupItem(a.item) {
		events.new(EventItemPickup, a.item, "Picked up "+a.item.Name())
		return ActionResult{true, 40}
	}

	events.new(EventPackFull, a.item, "Carrying too much!")
	return ActionResult{false, 0}
}

func (a *DropAction) Execute(g Game) ActionResult {
	p := g.Player()

	if a.item.IsEquipped() {
		events.new(EventItemDropped, a.item, "Can't drop an equipped item")
		return ActionResult{false, 0}
	}

	if p.DropItem(a.item) {
		events.new(EventItemDropped, a.item, fmt.Sprintf("Dropped the %s", a.item.Name()))
		return ActionResult{true, 40}
	}

	events.new(EventItemDropped, a.item, fmt.Sprintf("Can't drop the %s here", a.item.Name()))
	return ActionResult{false, 0}
}

func (a *UseAction) Execute(g Game) ActionResult {
	if a.item.use(g) {
		return ActionResult{true, 40}
	}

	return ActionResult{false, 0}
}

func (a *EquipAction) Execute(g Game) ActionResult {
	p := g.Player()

	if !a.item.IsEquipment() {
		return ActionResult{false, 0}
	}

	slot := a.item.EquipLocation()

	if p.IsEquipped(a.item) {
		p.UnequipItem(slot)

		msg := fmt.Sprintf("The %s was taken off", a.item.Name())
		if slot == equipLocationWeapon || slot == equipLocationShield || slot == equipLocationMissile {
			msg = fmt.Sprintf("No longer wielding a %s", a.item.Name())
		}

		events.new(EventItemUnequipped, a.item, msg)
	} else {
		// Is there an existing item in the slot?
		if i, ok := p.equipSlots[slot]; ok {
			p.UnequipItem(i.EquipLocation())
		}

		p.EquipItem(a.item, slot)

		msg := fmt.Sprintf("Now wearing a %s", a.item.Name())
		if slot == equipLocationWeapon || slot == equipLocationShield || slot == equipLocationMissile {
			msg = fmt.Sprintf("Wielding a %s", a.item.Name())
		}

		events.new(EventItemEquipped, a.item, msg)
	}

	return ActionResult{true, 40}
}
