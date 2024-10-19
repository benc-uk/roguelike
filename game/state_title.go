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

	if newGameRect := core.NewRect(120, 9*s.spSize, 60, s.spSize); s.TouchData.DidTapIn(newGameRect) {
		s.StartNewGame()
	}
}

func (s *TitleState) Draw(screen *ebiten.Image) {
	graphics.FgColour = graphics.ColourWhite
	graphics.BgColour = graphics.ColourTitle

	graphics.DrawBox(screen, 2, 1, VP_COLS-2, VP_ROWS-4)
	f := (math.Sin(float64(s.frameCount)*0.06) + 1) / 2
	graphics.BgColour = color.RGBA{0, 0, uint8(f * 255), 0xff}
	graphics.DrawBox(screen, 5, 1, VP_COLS-2, 2)

	graphics.BgColour = graphics.ColourTrans
	graphics.DrawTextRow(screen, fmt.Sprintf("%sGo WASM Roguelike", core.MakeStr(17, " ")), 6)

	graphics.DrawTextRow(screen, fmt.Sprintf("%sNEW GAME", core.MakeStr(20, " ")), 9)
	graphics.DrawTextRow(screen, fmt.Sprintf("%sLOAD GAME", core.MakeStr(20, " ")), 10)
	graphics.DrawTextRow(screen, fmt.Sprintf("%sQUIT", core.MakeStr(20, " ")), 11)
	graphics.DrawTextRow(screen, fmt.Sprintf("%d", s.seed), 17)

	s1 := s.bank.Sprite("slime")
	s2 := s.bank.Sprite("potion")
	s1.Draw(screen, 5*12, 6*12, color.White, true, false, false)
	s1.Draw(screen, 6*12, 6*12, color.White, true, false, false)
	s1.Draw(screen, 7*12, 6*12, color.White, true, false, false)
	s2.Draw(screen, 17*12+6, 6*12, color.White, true, false, false)
	s2.Draw(screen, 18*12+6, 6*12, color.White, true, false, false)
	s2.Draw(screen, 19*12+6, 6*12, color.White, true, false, false)

	// Draw the cursor
	graphics.FgColour = graphics.ColourCursor
	graphics.DrawTextRow(screen, fmt.Sprintf("%s⌦", core.MakeStr(18, " ")), s.cursor+9)
}
