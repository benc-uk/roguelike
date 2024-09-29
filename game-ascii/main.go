package main

import (
	"roguelike/core"
	"roguelike/engine"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
)

var basePath string = "./"
var game *engine.Game
var viewPort core.Rect

const (
	VP_ROWS       = 16 // Number of rows of tiles in the viewport
	VP_COLS       = 48 // Number of columns of tiles in the viewport
	MAX_EVENT_AGE = 6  // Max number of events to store
)

func main() {
	game = engine.NewGame(basePath+"assets/datafiles", 1111)
	viewPort = game.GetViewPort(VP_COLS, VP_ROWS)
	game.UpdateFOV(6)
	area, _ := pterm.DefaultArea.WithFullscreen().Start()

	drawScreen(area)

	// Listen for key presses loops like a game loop
	_ = keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.CtrlC {
			return true, nil // Stop listener by returning true on Ctrl+C
		}

		var move *engine.MoveAction
		switch key.Code {
		case keys.Right:
			move = engine.NewMoveAction(core.DirEast)
		case keys.Left:
			move = engine.NewMoveAction(core.DirWest)
		case keys.Up:
			move = engine.NewMoveAction(core.DirNorth)
		case keys.Down:
			move = engine.NewMoveAction(core.DirSouth)
		}

		if move != nil {
			move.Execute(*game)
			game.UpdateFOV(6)
			viewPort = game.GetViewPort(VP_COLS, VP_ROWS)
		}

		drawScreen(area)

		return false, nil
	})
}

func drawScreen(area *pterm.AreaPrinter) {
	gameMap := game.Map()
	p := game.Player()

	screen := ""
	for y := viewPort.Y; y < viewPort.Height+viewPort.Y; y++ {
		for x := viewPort.X; x < viewPort.Width+viewPort.X; x++ {
			tile := gameMap.Tile(x, y)
			appear := tile.GetAppearance()

			if appear == nil {
				screen += " "
				continue
			}

			if appear.Graphic == "wall" {
				if appear.InFOV {
					screen += pterm.Blue("#")
				} else {
					screen += pterm.Gray("#")
				}
				continue
			} else {
				if x == p.X && y == p.Y {
					screen += pterm.Yellow("@")
					continue
				}

				symbol := " "
				symColor := pterm.FgWhite

				if appear.Graphic == "floor" {
					symbol = "."
				}

				if appear.Graphic == "potion" {
					symbol = "P"
					symColor = pterm.FgGreen
				}

				if appear.Graphic == "sword" {
					symbol = "!"
					symColor = pterm.FgRed
				}

				if appear.InFOV {
					screen += symColor.Sprint(symbol)
				} else {
					screen += pterm.Gray(symbol)
				}
			}
		}
		screen += "\n"
	}

	// Update the area with the current time.
	area.Update(screen)
}
