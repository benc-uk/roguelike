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
	invCursor int
}

func (g *EbitenGame) SwitchStateInventory() {
	g.handlers[GameStateInventory].Init()
	g.state = GameStateInventory
}

func (s *InventoryState) Init() {
	s.invCursor = 0
}

func (s *InventoryState) PassEvent(e engine.GameEvent) {
}

func (s *InventoryState) Update(heldKeys []ebiten.Key, tappedKeys []ebiten.Key) {
	for _, key := range tappedKeys {
		if controls.Escape.IsKey(key) || controls.Inventory.IsKey(key) {
			s.state = GameStatePlaying
		}

		if controls.Up.IsKey(key) {
			s.invCursor--
			if s.invCursor < 0 {
				s.invCursor = 0
			}
		}

		if controls.Down.IsKey(key) {
			s.invCursor++
			if s.invCursor >= len(s.game.Player().Inventory()) {
				s.invCursor = len(s.game.Player().Inventory()) - 1
			}
		}

		if controls.Drop.IsKey(key) {
			p := s.game.Player()
			if len(p.Inventory()) > 0 {
				item := p.Inventory()[s.invCursor]

				p.DropItem(item)
				s.state = GameStatePlaying
			}
		}
	}
}

func (s *InventoryState) Draw(screen *ebiten.Image) {
	p := s.game.Player()

	// Draw the inventory screen
	graphics.DrawTextBox(screen, 0, 0, VP_COLS-1, 2, graphics.ColourInv)
	graphics.DrawTextBox(screen, 2, 0, VP_COLS-1, VP_ROWS-2, graphics.ColourInv)

	countCarried := len(p.Inventory())
	countMax := p.MaxItems()
	graphics.DrawTextRow(screen, fmt.Sprintf("  Backpack (%d/%d)", countCarried, countMax), 1, graphics.ColourTrans)

	// Draw the player's inventory
	for i, item := range p.Inventory() {
		curString := "    "
		if i == s.invCursor {
			curString = "  ⌦ "
		}

		equipLocDesc := ""
		if item.EquipLocation() != engine.EquipLocationNone {
			equipLocDesc = fmt.Sprintf("[%s]", item.EquipLocation())
		}

		graphics.DrawTextRow(screen, fmt.Sprintf("%s  %s", curString, item.Name()), i+3, graphics.ColourTrans)
		graphics.DrawTextRow(screen, fmt.Sprintf("                          %s", equipLocDesc), i+3, graphics.ColourTrans)

		sprite := s.bank.Sprite(item.Graphic())
		if sprite != nil {
			sprite.Draw(screen, 24, (i+3)*12, color.White, true, false, false)
		}
	}
}
