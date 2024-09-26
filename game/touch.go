package main

// ====================================================================================================================
// Touch handling for mobile devices (web)
// Lifted from https://github.com/hajimehoshi/ebiten/blob/main/examples/touch/main.go
// ====================================================================================================================

import (
	"roguelike/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type touch struct {
	originX, originY int
	currX, currY     int
	duration         int
	wasPinch, isPan  bool
}

type tap struct {
	X, Y int
}

func handleTaps(taps []tap, touches map[ebiten.TouchID]*touch) []tap {
	taps = taps[:0]
	for id, t := range touches {
		if inpututil.IsTouchJustReleased(id) {

			// If this one has not been touched long (30 frames can be assumed to be 500ms), or moved far, then it's a tap.
			dist := core.Pos{t.originX, t.originY}.Distance(core.Pos{t.currX, t.currY}) //nolint
			if !t.wasPinch && !t.isPan && (t.duration <= 30 || dist < 2) {
				taps = append(taps, tap{
					X: t.currX,
					Y: t.currY,
				})
			}

			delete(touches, id)
		}
	}

	return taps
}

func handleTouches(touchIDs []ebiten.TouchID, touches map[ebiten.TouchID]*touch) {
	// What touches are new in this frame?
	touchIDs = inpututil.AppendJustPressedTouchIDs(touchIDs[:0])
	for _, id := range touchIDs {
		x, y := ebiten.TouchPosition(id)
		touches[id] = &touch{
			originX: x, originY: y,
			currX: x, currY: y,
		}
	}
}
