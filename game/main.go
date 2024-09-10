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
	rows         = 12        // Rows of sprites/tiles on screen
	cols         = 16        // Columns of sprites/tiles on screen
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
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		move := engine.NewMoveAction(core.North)
		move.Execute(game.Player(), game.Map())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		move := engine.NewMoveAction(core.South)
		move.Execute(game.Player(), game.Map())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		move := engine.NewMoveAction(core.West)
		move.Execute(game.Player(), game.Map())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		move := engine.NewMoveAction(core.East)
		move.Execute(game.Player(), game.Map())
	}

	game.UpdateFOV(*game.Player())
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	wallSprite, err := bank.Sprite("Wall")
	if err != nil {
		log.Fatal(err)
	}
	playerSprite, err := bank.Sprite("Warrior")
	if err != nil {
		log.Fatal(err)
	}

	gameMap := game.Map()
	p := game.Player()

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			tile := gameMap.Tile(x, y)

			// Player takes precedence
			if p.X == x && p.Y == y {
				playerSprite.Draw(screen, x*sz, y*sz, palette, false)
				continue
			}

			// Then walls & floors
			appear := tile.GetAppearance()

			if appear.Details == "blank" {
				continue
			}

			if appear.Details == "wall" {
				wallSprite.Draw(screen, x*sz, y*sz, palette, !appear.InFOV)
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
