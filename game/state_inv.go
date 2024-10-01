package main

import (
	"fmt"
	"image/color"
	"roguelike/engine"
	"roguelike/game/controls"
	"roguelike/game/graphics"

	"github.com/hajimehoshi/ebiten/v2"
)

type InventoryState struct {
	// Neatly encapsulate the state of the game
	*EbitenGame

	// Internal vars for this state
	cursor    int
	inventory []*engine.Item
}

func (s *InventoryState) Init() {
	s.cursor = 0
	s.inventory = s.game.Player().Inventory()
}

func (s *InventoryState) PassEvent(e engine.GameEvent) {
}

func (s *InventoryState) Update(heldKeys []ebiten.Key, tappedKeys []ebiten.Key) {
	for _, key := range tappedKeys {
		if controls.Escape.IsKey(key) || controls.Inventory.IsKey(key) {
			s.state = GameStatePlaying
		}

		if controls.Up.IsKey(key) {
			s.cursor--
			if s.cursor < 0 {
				s.cursor = 0
			}
		}

		if controls.Down.IsKey(key) {
			s.cursor++
			if s.cursor >= len(s.inventory) {
				s.cursor = len(s.inventory) - 1
			}
		}

		if controls.Drop.IsKey(key) {
			a := engine.NewDropAction(s.inventory[s.cursor])
			_ = a.Execute(*s.game)
			s.state = GameStatePlaying
		}

		if controls.Select.IsKey(key) {
			a := engine.NewUseAction(s.inventory[s.cursor])

			result := a.Execute(*s.game)
			if result.Success {
				s.state = GameStatePlaying
			}
		}
	}
}

func (s *InventoryState) Draw(screen *ebiten.Image) {
	// Draw the inventory screen
	graphics.DrawTextBox(screen, 0, 0, VP_COLS-1, 2, graphics.ColourInv)
	graphics.DrawTextBox(screen, 2, 0, VP_COLS-1, VP_ROWS-2, graphics.ColourInv)

	countCarried := len(s.inventory)
	countMax := s.game.Player().BackpackSize()
	graphics.DrawTextRow(screen, fmt.Sprintf("  Backpack (%d/%d)", countCarried, countMax), 1, graphics.ColourTrans)

	// Draw the player's inventory
	for i, item := range s.inventory {

		equipLocDesc := ""
		if item.EquipLocation() != engine.EquipLocationNone {
			equipLocDesc = fmt.Sprintf("[%s]", item.EquipLocation())
		}

		graphics.DrawTextRow(screen, fmt.Sprintf("      %s", item.Name()), i+3, graphics.ColourTrans)
		graphics.DrawTextRow(screen, fmt.Sprintf("                          %s", equipLocDesc), i+3, graphics.ColourTrans)
		graphics.DrawTextRow(screen, "  ⌦", 3+s.cursor, graphics.ColourTrans)

		sprite := s.bank.Sprite(item.Graphic())
		if sprite != nil {
			sprite.Draw(screen, 24, (i+3)*12, color.White, true, false, false)
		}
	}
}
