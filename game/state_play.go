package main

import (
	"fmt"
	"roguelike/core"
	"roguelike/engine"
	"roguelike/game/controls"
	"roguelike/game/graphics"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PlayingState struct {
	// Neatly encapsulate the state of the game
	*EbitenGame

	// Internal vars for this state
	pickUpItem  bool
	playerLeft  bool
	delayFrames int
}

func (s *PlayingState) Init() {
}

func (s *PlayingState) PassEvent(e engine.GameEvent) {
}

func (s *PlayingState) Update(heldKeys []ebiten.Key, tappedKeys []ebiten.Key) {
	var action engine.Action
	player := s.game.Player()
	currTile := player.Tile()
	gm := s.game.Map()

	// Loop over taps (there should only be one for reasons)
	for _, tap := range s.taps {
		if tap.X < s.scrWidth/4 {
			action = engine.NewMoveAction(core.DirWest)
		} else if tap.X > s.scrWidth/4*3 {
			action = engine.NewMoveAction(core.DirEast)
		}

		if tap.Y < s.scrHeight/4 {
			action = engine.NewMoveAction(core.DirNorth)
		} else if tap.Y > s.scrHeight/4*3 {
			action = engine.NewMoveAction(core.DirSouth)
		}
	}

	// Held keys require a delay before moving the player
	for _, key := range heldKeys {
		if controls.Up.IsKey(key) && inpututil.KeyPressDuration(key) > 20 {
			action = engine.NewMoveAction(core.DirNorth)
		}
		if controls.Down.IsKey(key) && inpututil.KeyPressDuration(key) > 20 {
			action = engine.NewMoveAction(core.DirSouth)
		}
		if controls.Left.IsKey(key) && inpututil.KeyPressDuration(key) > 20 {
			action = engine.NewMoveAction(core.DirWest)
			s.playerLeft = true
		}
		if controls.Right.IsKey(key) && inpututil.KeyPressDuration(key) > 20 {
			action = engine.NewMoveAction(core.DirEast)
			s.playerLeft = false
		}
	}

	// Tapped keys (just pressed) reset the delayFrames and move the player immediately
	for _, key := range tappedKeys {
		if s.pickUpItem {
			if controls.Escape.IsKey(key) {
				s.pickUpItem = false
			}

			if key >= ebiten.KeyDigit0 && key <= ebiten.KeyDigit9 {
				index := int(key-ebiten.KeyDigit0) - 1
				if key == ebiten.KeyDigit0 {
					index = 9
				}

				tile := s.game.Player().Tile()
				items := tile.ListItems()
				if index < len(items) {
					item := items[index]
					a := engine.NewPickupAction(&item)
					res := a.Execute(*s.game)

					// Last item picked up, switch out of pickup mode
					if res.Success && len(items) == 1 {
						s.pickUpItem = false
					}
				}
			}

			return
		}

		var tappedDir core.Direction = -1

		if controls.Up.IsKey(key) {
			tappedDir = core.DirNorth
		}
		if controls.Down.IsKey(key) {
			tappedDir = core.DirSouth
		}
		if controls.Left.IsKey(key) {
			tappedDir = core.DirWest
			s.playerLeft = true
		}
		if controls.Right.IsKey(key) {
			tappedDir = core.DirEast
			s.playerLeft = false
		}

		if tappedDir >= 0 {
			destTile := currTile.AdjacentTileDir(tappedDir, gm)
			if destTile.Creature() != nil {
				action = engine.NewAttackAction(destTile.Creature())
			} else {
				action = engine.NewMoveAction(tappedDir)
				s.delayFrames = 0
				s.sfxPlayer.Play("walk")
			}
		}

		if controls.Inventory.IsKey(key) {
			s.state = GameStateInventory
			s.handlers[s.state].Init()
		}

		if controls.Get.IsKey(key) {
			appear := currTile.Appearance()

			// If it's not floor, player must be on one or more items
			// as we can't stand on walls or creatures
			if appear.Graphic != "floor" {
				s.pickUpItem = true
			}
		}
	}

	// This stops the whole game from running too fast
	if s.delayFrames > 0 {
		s.delayFrames--
		return
	}

	if action != nil {
		result := action.Execute(*s.game)

		// TODO: This logic is probably not right long term
		if !result.Success {
			return
		}

		// We translate the energy spent into frames to delay the game
		if result.EnergySpent > 0 {
			s.delayFrames = result.EnergySpent
		}

		// Play sound effects for walking
		if _, ok := action.(*engine.MoveAction); ok {
			s.sfxPlayer.Play("walk")
		}

		s.viewPort = s.game.GetViewPort(VP_COLS, VP_ROWS)
		s.game.UpdateFOV(s.viewDist)

		// Handle events and age them
		for _, e := range s.events {
			e.Age++
		}

		// Remove old events
		for i := 0; i < len(s.events); i++ {
			e := s.events[i]
			if e.Age > MAX_EVENT_AGE {
				s.events = append(s.events[:i], s.events[i+1:]...)
			}
		}
	}
}

func (s *PlayingState) Draw(screen *ebiten.Image) {
	gameMap := s.game.Map()
	p := s.game.Player()

	offsetX := s.viewPort.X * s.spSize
	offsetY := s.viewPort.Y * s.spSize

	// Draw the map
	for x := s.viewPort.X; x < s.viewPort.Width+s.viewPort.X; x++ {
		for y := s.viewPort.Y; y < s.viewPort.Height+s.viewPort.Y; y++ {
			tile := gameMap.Tile(x, y)
			appear := tile.Appearance()
			drawX := x*s.spSize - offsetX
			drawY := y*s.spSize - offsetY

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
			colour := s.palette[palIndex]

			// Walls
			if appear.Graphic == "wall" {
				s.bank.Sprite("wall").Draw(screen, drawX, drawY, s.palette[PAL_INDEX_WALL], appear.InFOV, false, false)
				continue
			}

			// Draw the player
			if x == p.Pos().X && y == p.Pos().Y {
				s.bank.Sprite("player").Draw(screen, drawX, drawY, s.palette[PAL_INDEX_PLAYER], appear.InFOV, s.playerLeft, false)
				continue
			}

			if appear.Graphic == "floor" {
				s.bank.Sprite("floor").Draw(screen, drawX, drawY, s.palette[PAL_INDEX_FLOOR], appear.InFOV, false, false)
				continue
			}

			// Then items/monsters/stuff that might have a sprite
			itemSprite := s.bank.Sprite(appear.Graphic)
			if itemSprite != nil {
				itemSprite.Draw(screen, drawX, drawY, colour, appear.InFOV, false, false)
				continue
			}
		}
	}

	// Draw the status bar, it was at row VP_ROWS-1 but we added a row for the status bar
	statusText := fmt.Sprintf("%s    ♥%d/%d   ⌘%d   ▼%d", p.Name(), p.HP(), p.MaxHP(), p.Exp(), p.Level())
	graphics.DrawTextRow(screen, statusText, VP_ROWS, graphics.ColourStatus)

	// Events & messages
	for i, e := range s.events {
		graphics.DrawTextRow(screen, e.Text(), i, graphics.ColourLog)
	}

	// Sub-mode for multiple items
	if s.pickUpItem {
		t := s.game.Player().Tile()

		bodyText := ""
		for i, item := range t.ListItems() {
			bodyText += fmt.Sprintf("%d: %s\n", i+1, item.Name())
		}

		bodyText = strings.Trim(bodyText, "\n")
		graphics.DrawDialogBox(screen, VP_COLS-1, "Pickup an item, using number keys", bodyText)
	}
}
