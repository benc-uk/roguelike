package controls

// ====================================================================================================================
// Touch handling for mobile devices (web) and mouse clicks too
// Lifted from https://github.com/hajimehoshi/ebiten/blob/main/examples/touch/main.go
// ====================================================================================================================

import (
	"roguelike/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Touch struct {
	core.Pos
	originX, originY int
	ID               ebiten.TouchID
}

type TouchData struct {
	touches  map[ebiten.TouchID]*Touch
	touchIDs []ebiten.TouchID
	taps     []core.Pos
}

func NewTouchData() TouchData {
	return TouchData{
		touches:  make(map[ebiten.TouchID]*Touch),
		touchIDs: make([]ebiten.TouchID, 0),
		taps:     make([]core.Pos, 0),
	}
}

func (td *TouchData) Update() {
	td.taps = td.taps[:0]
	for id, t := range td.touches {
		if inpututil.IsTouchJustReleased(id) {
			// If this one has not been touched long (30 frames can be assumed to be 500ms), or moved far, then it's a tap.
			dist := core.Pos{X: t.originX, Y: t.originY}.Distance(core.Pos{X: t.X, Y: t.Y})
			if dist < 3 {
				td.taps = append(td.taps, core.Pos{
					X: t.X,
					Y: t.Y,
				})
			}
			delete(td.touches, id)
		}

	}

	// We also handle mouse input as a touch
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		td.taps = append(td.taps, core.Pos{
			X: x, Y: y,
		})
	}

	// What touches are new in this frame?
	td.touchIDs = inpututil.AppendJustPressedTouchIDs(td.touchIDs[:0])
	for _, id := range td.touchIDs {
		x, y := ebiten.TouchPosition(id)
		td.touches[id] = &Touch{
			originX: x, originY: y,
			Pos: core.Pos{X: x, Y: y},
			ID:  id,
		}
	}
}

// Get the taps that happened in this frame
func (td *TouchData) Taps() []core.Pos {
	return td.taps
}

func (td *TouchData) Touches() []Touch {
	touches := make([]Touch, 0)
	for _, t := range td.touches {
		touches = append(touches, *t)
	}

	return touches
}

// Was the given rectangle tapped in in this frame
func (td *TouchData) DidTapIn(r core.Rect) bool {
	for _, tap := range td.taps {
		if r.ContainsPos(tap) {
			return true
		}
	}

	return false
}
