package engine

import (
	"fmt"
	"strconv"
)

// ============================================================================
// Effects of equipment and other things
// ============================================================================

type effectType int

const (
	effectTypeNone effectType = iota
	effectTypeDefence
	effectTypeAttackChance
	effectTypeAttackDamage
	effectTypeAttackRoll
)

type effect struct {
	effectType effectType
	value      int
	roll       *DiceRoll
}

func newEffect(effectName string, effectValue string) *effect {
	val, _ := strconv.ParseInt(effectValue, 10, 64)
	roll, _ := ParseDiceRoll(effectValue)

	switch effectName {
	case "defence":
		return &effect{effectTypeDefence, int(val), nil}
	case "attack":
		return &effect{effectTypeAttackRoll, 0, &roll}
	case "toHit":
		return &effect{effectTypeAttackChance, int(val), nil}
	case "damage":
		return &effect{effectTypeAttackDamage, int(val), nil}
	}

	return &effect{effectTypeNone, 0, nil}
}

func (e effect) apply(p *Player) {
	switch e.effectType {
	case effectTypeDefence:
		p.defence += e.value
	case effectTypeAttackChance:
		p.attackChance += e.value
	case effectTypeAttackDamage:
		p.attackDamage += e.value
	case effectTypeAttackRoll:
		p.attackRoll = *e.roll
	}
}

func (e effect) remove(p *Player) {
	switch e.effectType {
	case effectTypeDefence:
		p.defence -= e.value
	case effectTypeAttackChance:
		p.attackChance -= e.value
	case effectTypeAttackDamage:
		p.attackDamage -= e.value
	case effectTypeAttackRoll:
		p.attackRoll = DiceRoll{0, 0, 0}
	}
}

func (e effect) description() string {
	mod := fmt.Sprintf("+%d", e.value)
	if e.value < 0 {
		mod = fmt.Sprintf("%d", e.value)
	}

	switch e.effectType {
	case effectTypeDefence:
		return fmt.Sprintf("%sdef", mod)
	case effectTypeAttackChance:
		return fmt.Sprintf("%shit", mod)
	case effectTypeAttackDamage:
		return fmt.Sprintf("%sdam", mod)
	case effectTypeAttackRoll:
		return e.roll.String()
	}

	return "Unknown"
}
