package main

import (
	"fmt"
	"image/color"
	"roguelike/core"
	"roguelike/engine"
	"roguelike/game/controls"
	"roguelike/game/graphics"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *EbitenGame) UpdatePlaying(heldKeys []ebiten.Key, tappedKeys []ebiten.Key) {
	var move *engine.MoveAction

	// Loop over taps (there should only be one for reasons)
	for _, tap := range g.taps {
		if tap.X < g.scrWidth/4 {
			move = engine.NewMoveAction(core.DirWest)
		} else if tap.X > g.scrWidth/4*3 {
			move = engine.NewMoveAction(core.DirEast)
		}

		if tap.Y < g.scrHeight/4 {
			move = engine.NewMoveAction(core.DirNorth)
		} else if tap.Y > g.scrHeight/4*3 {
			move = engine.NewMoveAction(core.DirSouth)
		}
	}

	// Held keys require a delay before moving the player
	for _, key := range heldKeys {
		if controls.Up.IsKey(key) && inpututil.KeyPressDuration(key) > 20 {
			move = engine.NewMoveAction(core.DirNorth)
		}
		if controls.Down.IsKey(key) && inpututil.KeyPressDuration(key) > 20 {
			move = engine.NewMoveAction(core.DirSouth)
		}
		if controls.Left.IsKey(key) && inpututil.KeyPressDuration(key) > 20 {
			move = engine.NewMoveAction(core.DirWest)
			g.playerLeft = true
		}
		if controls.Right.IsKey(key) && inpututil.KeyPressDuration(key) > 20 {
			move = engine.NewMoveAction(core.DirEast)
			g.playerLeft = false
		}
	}

	// Tapped keys (just pressed) reset the delayFrames and move the player immediately
	for _, key := range tappedKeys {
		if controls.Up.IsKey(key) {
			move = engine.NewMoveAction(core.DirNorth)
			g.delayFrames = 0
		}
		if controls.Down.IsKey(key) {
			move = engine.NewMoveAction(core.DirSouth)
			g.delayFrames = 0
		}
		if controls.Left.IsKey(key) {
			move = engine.NewMoveAction(core.DirWest)
			g.playerLeft = true
			g.delayFrames = 0
		}
		if controls.Right.IsKey(key) {
			move = engine.NewMoveAction(core.DirEast)
			g.playerLeft = false
			g.delayFrames = 0
		}

		if controls.Inventory.IsKey(key) {
			if g.state == GameStateInventory {
				g.state = GameStatePlaying
			} else {
				g.SwitchStateInventory()
			}
		}

		if controls.Escape.IsKey(key) {
			if g.state == GameStateInventory {
				g.state = GameStatePlaying
			}
		}
	}

	// This stops the whole game from running too fast
	if g.delayFrames > 0 {
		g.delayFrames--
		return
	}

	if move != nil {
		result := move.Execute(*g.game)
		if !result.Success {
			return
		}

		// We translate the energy spent into frames to delay the game
		if result.EnergySpent > 0 {
			g.delayFrames = result.EnergySpent
		}

		g.effect.Play("walk")

		g.viewPort = g.game.GetViewPort(VP_COLS, VP_ROWS)
		g.game.UpdateFOV(g.viewDist)

		// Handle events and age them
		for _, e := range g.events {
			e.Age++
		}

		// Remove old events
		for i := 0; i < len(g.events); i++ {
			e := g.events[i]
			if e.Age > MAX_EVENT_AGE {
				g.events = append(g.events[:i], g.events[i+1:]...)
			}
		}
	}
}

func (g *EbitenGame) DrawPlaying(screen *ebiten.Image) {
	gameMap := g.game.Map()
	p := g.game.Player()

	offsetX := g.viewPort.X * g.spSize
	offsetY := g.viewPort.Y * g.spSize

	// Draw the map
	for x := g.viewPort.X; x < g.viewPort.Width+g.viewPort.X; x++ {
		for y := g.viewPort.Y; y < g.viewPort.Height+g.viewPort.Y; y++ {
			tile := gameMap.Tile(x, y)
			appear := tile.GetAppearance(gameMap)
			drawX := x*g.spSize - offsetX
			drawY := y*g.spSize - offsetY

			// Unseen areas are blank/not drawn
			if appear == nil {
				continue
			}

			palIndex := 0
			if appear.Colour != "" {
				if i, err := strconv.Atoi(appear.Colour); err == nil {
					palIndex = i
				}
			}
			colour := g.palette[palIndex]

			// Walls
			if appear.Graphic == "wall" {
				g.bank.Sprite("wall").Draw(screen, drawX, drawY, g.palette[PAL_INDEX_WALL], appear.InFOV, false, false)
				continue
			}

			// Draw the player
			if x == p.Pos().X && y == p.Pos().Y {
				g.bank.Sprite("player").Draw(screen, drawX, drawY, g.palette[PAL_INDEX_PLAYER], appear.InFOV, g.playerLeft, false)
				continue
			}

			if appear.Graphic == "floor" {
				g.bank.Sprite("floor").Draw(screen, drawX, drawY, g.palette[PAL_INDEX_FLOOR], appear.InFOV, false, false)
				continue
			}

			// Then items/monsters/stuff that might have a sprite
			itemSprite := g.bank.Sprite(appear.Graphic)
			if itemSprite != nil {
				itemSprite.Draw(screen, drawX, drawY, colour, appear.InFOV, false, false)
				continue
			}
		}
	}

	// Draw the status bar, it was at row VP_ROWS-1 but we added a row for the status bar
	statusText := fmt.Sprintf("%s    ♥%d/%d   ⌘%d   ▼%d", p.Name(), p.CurrentHP(), p.MaxHP(), p.Exp(), p.Level())
	graphics.DrawTextRow(screen, statusText, VP_ROWS, color.RGBA{0x10, 0x50, 0x10, 0xff})

	for i, e := range g.events {
		graphics.DrawTextRow(screen, e.Text, i, color.RGBA{0x00, 0x00, 0x30, 0x30})
	}
}
