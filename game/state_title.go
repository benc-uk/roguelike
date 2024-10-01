package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"roguelike/core"
	"roguelike/engine"
	"roguelike/game/controls"
	"roguelike/game/graphics"

	"github.com/hajimehoshi/ebiten/v2"
)

type TitleState struct {
	// Neatly encapsulate the state of the game
	*EbitenGame

	cursor     int
	quickStart bool
}

func (s *TitleState) Init() {
	s.cursor = 0
}

func (s *TitleState) PassEvent(e engine.GameEvent) {
}

func (s *TitleState) Update(heldKeys []ebiten.Key, tappedKeys []ebiten.Key) {
	if s.quickStart {
		log.Println("Quick starting, skipping title screen and starting new game")
		s.StartNewGame()
		s.quickStart = false
		return
	}

	for _, key := range tappedKeys {
		if controls.Up.IsKey(key) {
			s.cursor--
			if s.cursor < 0 {
				s.cursor = 0
				s.flashCount = 2
			}
		}

		if controls.Down.IsKey(key) {
			s.cursor++
			if s.cursor > 2 {
				s.cursor = 2
				s.flashCount = 2
			}
		}

		if controls.Select.IsKey(key) {
			switch s.cursor {
			case 0:
				s.StartNewGame()
			}
		}
	}
}

func (s *TitleState) Draw(screen *ebiten.Image) {
	graphics.DrawTextBox(screen, 2, 0, VP_COLS-1, VP_ROWS-4, graphics.ColourTitle)
	f := (math.Sin(float64(s.frameCount)*0.06) + 1) / 2
	graphics.DrawTextBox(screen, 5, 0, VP_COLS-1, 2, color.RGBA{0, 20, uint8((200 * f) + 55), 0xff})
	graphics.DrawTextRow(screen, fmt.Sprintf("%sGo WASM Roguelike", core.MakeStr(17, " ")), 6, graphics.ColourTrans)

	graphics.DrawTextRow(screen, fmt.Sprintf("%sNEW GAME", core.MakeStr(20, " ")), 9, graphics.ColourTrans)
	graphics.DrawTextRow(screen, fmt.Sprintf("%sLOAD GAME", core.MakeStr(20, " ")), 10, graphics.ColourTrans)
	graphics.DrawTextRow(screen, fmt.Sprintf("%sQUIT", core.MakeStr(20, " ")), 11, graphics.ColourTrans)
	graphics.DrawTextRow(screen, fmt.Sprintf("%d", s.seed), 17, graphics.ColourTrans)

	// Draw the cursor
	graphics.DrawTextRow(screen, fmt.Sprintf("%s⌦", core.MakeStr(18, " ")), s.cursor+9, graphics.ColourTrans)

	s1 := s.bank.Sprite("slime")
	s2 := s.bank.Sprite("potion")
	s1.Draw(screen, 5*12, 6*12, color.White, true, false, false)
	s1.Draw(screen, 6*12, 6*12, color.White, true, false, false)
	s1.Draw(screen, 7*12, 6*12, color.White, true, false, false)
	s2.Draw(screen, 17*12+6, 6*12, color.White, true, false, false)
	s2.Draw(screen, 18*12+6, 6*12, color.White, true, false, false)
	s2.Draw(screen, 19*12+6, 6*12, color.White, true, false, false)
}
