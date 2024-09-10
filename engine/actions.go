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

func (a *MoveAction) Execute(p *Player) bool {
	p.X += a.Direction.Pos().X
	p.Y += a.Direction.Pos().Y

	return true
}
