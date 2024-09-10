package engine

import "roguelike/core"

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
	if !newPos.InBounds(m.width, m.height) {
		return false
	}

	tile := m.tiles[newPos.X][newPos.Y]
	if tile.blocksMove {
		return false
	}

	p.X += a.Direction.Pos().X
	p.Y += a.Direction.Pos().Y

	return true
}
