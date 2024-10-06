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
	cursor         int
	inv            []*engine.Item
	item           *engine.Item
	describingItem bool
}

func (s *InventoryState) Init() {
	s.cursor = 0
	s.inv = s.game.Player().Inventory()
	s.item = s.inv[s.cursor]
	s.describingItem = false
}

func (s *InventoryState) PassEvent(e engine.GameEvent) {
}

func (s *InventoryState) Update(heldKeys []ebiten.Key, tappedKeys []ebiten.Key) {
	for _, key := range tappedKeys {
		if controls.Escape.IsKey(key) || controls.Inventory.IsKey(key) {
			s.state = gameStatePlaying
		}

		if s.describingItem {
			if controls.Info.IsKey(key) || controls.Escape.IsKey(key) {
				s.describingItem = false
			}
			continue
		}

		if controls.Up.IsKey(key) {
			s.cursor--
			if s.cursor < 0 {
				s.cursor = 0
			}
		}

		if controls.Down.IsKey(key) {
			s.cursor++
			if s.cursor >= len(s.inv) {
				s.cursor = len(s.inv) - 1
			}
		}

		// Update the item we're looking at
		s.item = s.inv[s.cursor]

		if controls.Drop.IsKey(key) {
			a := engine.NewDropAction(s.item)
			_ = a.Execute(*s.game)
			s.state = gameStatePlaying
		}

		if controls.Info.IsKey(key) {
			s.describingItem = true
		}

		if controls.Select.IsKey(key) {
			var a engine.Action
			if s.item.IsEquipment() {
				a = engine.NewEquipAction(s.item)
			} else {
				a = engine.NewUseAction(s.item)
			}

			if a == nil {
				continue
			}

			result := a.Execute(*s.game)
			if result.Success {
				s.state = gameStatePlaying
			}
		}
	}
}

func (s *InventoryState) Draw(screen *ebiten.Image) {
	// Draw the inventory screen
	graphics.BgColour = graphics.ColourInv
	graphics.FgColour = graphics.ColourWhite

	graphics.DrawBox(screen, 0, 1, VP_COLS-2, 2)
	graphics.DrawBox(screen, 2, 1, VP_COLS-2, VP_ROWS-2)

	graphics.BgColour = graphics.ColourTrans

	// In mode where we're describing an item
	if s.describingItem {
		graphics.DrawTextRow(screen, fmt.Sprintf("   %s", s.item.NameTitle()), 1)
		text := "Type: " + s.item.ItemType() + "\n"
		text += "Rarity: " + s.item.Rarity().String() + "\n"
		if s.item.IsEquipment() {
			text += "Equip Location: " + s.item.EquipLocation().String() + "\n"
		}
		text += "\n" + s.item.Description()
		graphics.DrawWrappedText(screen, text, 3, 3, VP_COLS+18)

		// Skip over the normal inventory drawing
		return
	}

	countCarried := len(s.inv)
	countMax := s.game.Player().BackpackSize()
	graphics.DrawTextRow(screen, fmt.Sprintf("   Backpack (%d/%d)", countCarried, countMax), 1)

	// Draw the player's inventory item by item
	for i, item := range s.inv {
		// Just in case we have a nil item (shouldn't happen)
		if item == nil {
			continue
		}

		graphics.FgColour = graphics.ColourWhite

		extra1 := ""
		if item.IsEquipment() {
			graphics.FgColour = color.RGBA{76, 147, 230, 255}
			if s.game.Player().IsEquipped(item) {
				extra1 = fmt.Sprintf("[%s]", item.EquipLocation())
			}
		}

		extra2 := ""

		sprite := s.bank.Sprite(item.Graphic())
		if sprite != nil {
			sprite.Draw(screen, 30, (i+3)*12, graphics.FgColour, true, false, false)
		}

		graphics.DrawTextRow(screen, fmt.Sprintf("       %s %s %s", item.NameTitle(), extra1, extra2), i+3)
		graphics.FgColour = graphics.ColourTurq
		graphics.DrawTextRow(screen, "   ⌦", 3+s.cursor)
	}
}
