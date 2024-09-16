package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"roguelike/core"
	"roguelike/engine"
	"roguelike/game/graphics"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// These are injected by the build system
var basePath string = "./"
var version string = "0.0.1-alpha_005"

//go:embed icon.png
var iconBytes []byte // Icon for the window is embedded

// Left global for now, could be refactored
var game *engine.Game

const (
	wallPalIndex   = 0
	floorPalIndex  = 3
	playerPalIndex = 10
	rows           = 16
	cols           = 24
)

// Implements the ebiten.Game interface
// Holds external state for the rendering & running of the game
type EbitenGame struct {
	//state GameState

	// Core consts for rendering the window
	spSize       int // const - size of each tile in pixels
	scrWidth     int // const - screen width in pixels
	scrHeight    int // const - screen height in pixels
	initialScale int // const - initial scale

	// Graphics
	bank    *graphics.SpriteBank // Sprite bank holds all the sprites
	palette color.Palette        // Current palette

	// Viewport & FOV
	viewPort core.Rect // The area of the map that is visible
	viewDist int       // View distance in tiles (const)

	// Weird crap for touch controls
	touches  map[ebiten.TouchID]*touch
	touchIDs []ebiten.TouchID
	taps     []tap
}

func init() {

}

func (g *EbitenGame) Update() error {
	p := game.Player()
	updateViewPort := false

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		move := engine.NewMoveAction(core.North)
		move.Execute(p, game.Map())
		updateViewPort = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		move := engine.NewMoveAction(core.South)
		move.Execute(p, game.Map())
		updateViewPort = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		move := engine.NewMoveAction(core.West)
		move.Execute(p, game.Map())
		updateViewPort = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		move := engine.NewMoveAction(core.East)
		move.Execute(p, game.Map())
		updateViewPort = true
	}

	// Touch controls - figure out taps and touches
	g.taps = handleTaps(g.taps, g.touches)
	handleTouches(g.touchIDs, g.touches)

	// Loop over taps (there should only be one)
	for _, tap := range g.taps {

		if tap.X < g.scrWidth/4 {
			move := engine.NewMoveAction(core.West)
			move.Execute(p, game.Map())
			updateViewPort = true

		} else if tap.X > g.scrWidth/4*3 {
			move := engine.NewMoveAction(core.East)
			move.Execute(p, game.Map())
			updateViewPort = true
		}

		if tap.Y < g.scrHeight/4 {
			move := engine.NewMoveAction(core.North)
			move.Execute(p, game.Map())
			updateViewPort = true
		} else if tap.Y > g.scrHeight/4*3 {
			move := engine.NewMoveAction(core.South)
			move.Execute(p, game.Map())
			updateViewPort = true
		}
	}

	if updateViewPort {
		g.UpdateViewAndFOV()
	}

	return nil
}

func (g *EbitenGame) UpdateViewAndFOV() {
	gameMap := game.Map()
	p := game.Player()
	// ViewPort is the area of the map that is visible centered on the player
	g.viewPort = core.NewRect(p.X-cols/2, p.Y-rows/2, cols, rows)

	if g.viewPort.X < 0 {
		g.viewPort.X = 0
	}
	if g.viewPort.Y < 0 {
		g.viewPort.Y = 0
	}
	if g.viewPort.X+g.viewPort.Width > gameMap.Size().Width {
		g.viewPort.X = gameMap.Size().Width - g.viewPort.Width
	}
	if g.viewPort.Y+g.viewPort.Height > gameMap.Size().Height {
		g.viewPort.Y = gameMap.Size().Height - g.viewPort.Height
	}

	game.UpdateFOV(*p, g.viewDist)
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	// Clear
	screen.Fill(color.Black)

	gameMap := game.Map()
	p := game.Player()

	offsetX := g.viewPort.X * g.spSize
	offsetY := g.viewPort.Y * g.spSize

	// Draw the map
	for x := g.viewPort.X; x < g.viewPort.Width+g.viewPort.X; x++ {
		for y := g.viewPort.Y; y < g.viewPort.Height+g.viewPort.Y; y++ {
			tile := gameMap.Tile(x, y)
			appear := tile.GetAppearance(gameMap)
			drawX := x*g.spSize - offsetX
			drawY := y*g.spSize - offsetY

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
			colour := g.palette[palIndex]

			// Walls
			if appear.Graphic == "wall" {
				g.bank.Sprite("wall").Draw(screen, drawX, drawY, g.palette[wallPalIndex], appear.InFOV)
				continue
			}

			// Draw the player in the center of the screen
			if x == p.X && y == p.Y {
				g.bank.Sprite("player").Draw(screen, drawX, drawY, g.palette[playerPalIndex], appear.InFOV)
				continue
			}

			if appear.Graphic == "floor" {
				g.bank.Sprite("floor").Draw(screen, drawX, drawY, g.palette[floorPalIndex], appear.InFOV)
				continue
			}

			// Then items/monsters/stuff
			itemSprite := g.bank.Sprite(appear.Graphic)
			if itemSprite != nil {
				itemSprite.Draw(screen, drawX, drawY, colour, appear.InFOV)
				continue
			}
		}
	}
}

func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.scrWidth, g.scrHeight
}

func main() {
	log.Printf("Generic Dungeon Game v%s is starting...", version)

	// Create image for window icon from embedded bytes
	buf := bytes.NewBuffer(iconBytes)
	icon, _, err := image.Decode(buf)
	if err != nil {
		log.Fatal(err)
	}

	bank, err := graphics.NewSpriteBank(basePath + "assets/sprites.json")
	if err != nil {
		log.Fatal(err)
	}

	err = graphics.LoadPalettes(basePath + "assets/palettes.yaml")
	if err != nil {
		log.Fatal(err)
	}

	palette, err := graphics.GetRGBPalette("default")
	if err != nil {
		log.Fatal(err)
	}

	spSize := bank.Size()
	ebitenGame := &EbitenGame{
		touches:   make(map[ebiten.TouchID]*touch),
		spSize:    spSize,
		scrWidth:  cols * spSize,
		scrHeight: rows * spSize,

		initialScale: spSize / 2,
		bank:         bank,
		palette:      palette,
		viewPort:     core.NewRect(0, 0, cols, rows),
		viewDist:     6,
	}

	ebiten.SetWindowSize(ebitenGame.scrWidth*ebitenGame.initialScale, ebitenGame.scrHeight*ebitenGame.initialScale)
	ebiten.SetWindowPosition(0, 0)
	ebiten.SetWindowTitle("Generic Dungeon Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowIcon([]image.Image{icon})

	// TODO: This is a lot of placeholder for now
	engine.LoadItemFactory()
	game = engine.NewGame()
	ebitenGame.UpdateViewAndFOV()

	if err := ebiten.RunGame(ebitenGame); err != nil {
		log.Fatal(err)
	}
}
