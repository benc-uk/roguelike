package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"roguelike/core"
	"roguelike/engine"
	"roguelike/game/graphics"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// These are injected by the build system
var basePath string = "./"
var version string = "0.0.1-alpha_003"

// These are not
var windowIcon image.Image
var bank *graphics.SpriteBank
var palette color.Palette
var game *engine.Game

const (
	sz           = 12        // Sprite & tile size
	rows         = 16        // Rows of sprites/tiles on screen
	cols         = 20        // Columns of sprites/tiles on screen
	scrWidth     = cols * sz // Screen width in pixels
	scrHeight    = rows * sz // Screen height in pixels
	initialScale = 5         // Initial scale of the window
)

func init() {
	log.Printf("Generic Dungeon Game v%s is starting...", version)

	var err error
	bank, err = graphics.NewSpriteBank(basePath + "assets/sprites/sprites.json")
	if err != nil {
		log.Fatal(err)
	}

	err = graphics.LoadPalettes(basePath + "assets/palettes.json")
	if err != nil {
		log.Fatal(err)
	}

	palette, err = graphics.GetRGBPalette("default")
	if err != nil {
		log.Fatal(err)
	}

	iconImg, _, err := ebitenutil.NewImageFromFile(basePath + "assets/icon.png")
	if err != nil {
		log.Fatal(err)
	}
	windowIcon = iconImg
}

type Game struct{}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		move := engine.NewMoveAction(core.North)
		move.Execute(game.Player(), game.Map())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		move := engine.NewMoveAction(core.South)
		move.Execute(game.Player(), game.Map())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		move := engine.NewMoveAction(core.West)
		move.Execute(game.Player(), game.Map())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		move := engine.NewMoveAction(core.East)
		move.Execute(game.Player(), game.Map())
	}

	// touch controls
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		// click on left side of screen
		if x < scrWidth/4 {
			move := engine.NewMoveAction(core.West)
			move.Execute(game.Player(), game.Map())
		}
		if x > scrWidth*3/4 {
			move := engine.NewMoveAction(core.East)
			move.Execute(game.Player(), game.Map())
		}

		// click on top side of screen
		if y < scrHeight/4 {
			move := engine.NewMoveAction(core.North)
			move.Execute(game.Player(), game.Map())
		}
		if y > scrHeight*3/4 {
			move := engine.NewMoveAction(core.South)
			move.Execute(game.Player(), game.Map())
		}
	}

	game.UpdateFOV(*game.Player(), 6)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	wallSprite, err := bank.Sprite("wall")
	if err != nil {
		log.Fatal(err)
	}
	playerSprite, err := bank.Sprite("warrior")
	if err != nil {
		log.Fatal(err)
	}

	gameMap := game.Map()
	p := game.Player()

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			tile := gameMap.Tile(x, y)
			appear := tile.GetAppearance(gameMap)

			// Unseen areas are blank/not drawn
			if appear == nil {
				continue
			}

			if appear.Details == "wall" {
				wallSprite.Draw(screen, x*sz, y*sz, palette, !appear.InFOV)
				continue
			}

			// Floors are drawn first
			if !appear.InFOV {
				vector.DrawFilledRect(screen, float32(x*sz), float32(y*sz), float32(sz), float32(sz), color.RGBA{15, 15, 15, 255}, false)
			} else {
				vector.DrawFilledRect(screen, float32(x*sz), float32(y*sz), float32(sz), float32(sz), color.RGBA{30, 30, 30, 255}, false)
			}

			// Player after floors
			if p.X == x && p.Y == y {
				playerSprite.Draw(screen, x*sz, y*sz, palette, false)
				continue
			}

			// Then items
			itemSprite, _ := bank.Sprite(appear.Details)
			if itemSprite != nil {
				palIndex := itemSprite.PaletteIndex()

				// Check for hints, which can override the palette index
				if len(appear.Hints) > 0 {
					for _, hint := range appear.Hints {
						hintParts := strings.Split(hint, "::")
						if len(hintParts) == 2 && hintParts[0] == "colour" {
							palIndex, _ = strconv.Atoi(hintParts[1])
							break
						}
					}
				}

				itemSprite.DrawWithColour(screen, x*sz, y*sz, palette, palIndex, !appear.InFOV)
			}

		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return scrWidth, scrHeight
}

func main() {
	ebiten.SetWindowSize(scrWidth*initialScale, scrHeight*initialScale)
	ebiten.SetWindowPosition(0, 0)
	ebiten.SetWindowTitle("Generic Dungeon Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowIcon([]image.Image{windowIcon})

	// TODO: This is a lot of placeholder for now
	engine.LoadItemFactory()
	game = engine.NewGame()

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
