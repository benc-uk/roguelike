package main

import (
	"fmt"
	"image/color"
	"roguelike/game/controls"
	"roguelike/game/graphics"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *EbitenGame) SwitchStateInventory() {
	g.state = GameStateInventory
	g.invCursor = 0
}

func (g *EbitenGame) UpdateInv(heldKeys []ebiten.Key, tappedKeys []ebiten.Key) {
	for _, key := range tappedKeys {
		if controls.Escape.IsKey(key) || controls.Inventory.IsKey(key) {
			g.state = GameStatePlaying
		}

		if controls.Up.IsKey(key) {
			g.invCursor--
			if g.invCursor < 0 {
				g.invCursor = 0
			}
		}

		if controls.Down.IsKey(key) {
			g.invCursor++
			if g.invCursor >= len(g.game.Player().Inventory()) {
				g.invCursor = len(g.game.Player().Inventory()) - 1
			}
		}

		if controls.Drop.IsKey(key) {
			p := g.game.Player()
			if len(p.Inventory()) > 0 {
				p.DropItem(g.invCursor)
				g.state = GameStatePlaying
			}
		}
	}
}

func (g *EbitenGame) DrawInv(screen *ebiten.Image) {
	p := g.game.Player()

	// Draw the inventory screen
	countCarried := len(p.Inventory())
	countMax := p.MaxItems()
	graphics.DrawTextRow(screen, fmt.Sprintf(" Inventory (%d/%d)", countCarried, countMax), 0, color.RGBA{0x30, 0x00, 0x30, 0xff})

	// Draw the player's inventory
	for i, item := range p.Inventory() {
		curString := "  "
		if i == g.invCursor {
			curString = "⌦ "
		}

		graphics.DrawTextRow(screen, fmt.Sprintf("%s  %s", curString, item.Name()), i+1, color.RGBA{0x30, 0x00, 0x30, 0xff})

		sprite := g.bank.Sprite(item.Graphic())
		if sprite != nil {
			sprite.Draw(screen, 12, (i+1)*12, color.White, true, false, false)
		}

	}
}
