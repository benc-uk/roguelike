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

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// These are injected by the build system
var basePath string = "./"
var version string = "0.0.1-alpha_004"

// These are not
var windowIcon image.Image

// Left global for now, could be refactored
var game *engine.Game

const (
	wallPalIndex   = 0
	floorPalIndex  = 3
	playerPalIndex = 10
)

// Implements the ebiten.Game interface holds external state for the rendering game and IO
type EbitenGame struct {
	// Touch controls
	touches  map[ebiten.TouchID]*touch
	touchIDs []ebiten.TouchID
	taps     []tap

	// Core consts for rendering
	sz           int // const - size of each tile
	rows         int // const - number of rows to render
	cols         int // const - number of cols to render
	scrWidth     int // const - screen width
	scrHeight    int // const - screen height
	initialScale int // const - initial scale
	viewDist     int // const - view distance

	// Graphics
	bank     *graphics.SpriteBank
	palette  color.Palette
	viewPort core.Rect
}

func init() {
	log.Printf("Generic Dungeon Game v%s is starting...", version)

	iconImg, _, err := ebitenutil.NewImageFromFile(basePath + "assets/icon.png")
	if err != nil {
		log.Fatal(err)
	}
	windowIcon = iconImg
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

	// touch controls
	// What touches have just ended?
	g.taps = g.taps[:0]
	for id, t := range g.touches {
		if inpututil.IsTouchJustReleased(id) {

			// If this one has not been touched long (30 frames can be assumed to be 500ms), or moved far, then it's a tap.
			diff := core.DistanceF(t.originX, t.originY, t.currX, t.currY)
			if !t.wasPinch && !t.isPan && (t.duration <= 30 || diff < 2) {
				g.taps = append(g.taps, tap{
					X: t.currX,
					Y: t.currY,
				})
			}

			delete(g.touches, id)
		}
	}

	// What touches are new in this frame?
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		x, y := ebiten.TouchPosition(id)
		g.touches[id] = &touch{
			originX: x, originY: y,
			currX: x, currY: y,
		}
	}

	g.touchIDs = ebiten.AppendTouchIDs(g.touchIDs[:0])

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
	g.viewPort = core.NewRect(p.X-g.cols/2, p.Y-g.rows/2, g.cols, g.rows)

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

	offsetX := g.viewPort.X * g.sz
	offsetY := g.viewPort.Y * g.sz

	// Draw the map
	for x := g.viewPort.X; x < g.viewPort.Width+g.viewPort.X; x++ {
		for y := g.viewPort.Y; y < g.viewPort.Height+g.viewPort.Y; y++ {
			tile := gameMap.Tile(x, y)
			appear := tile.GetAppearance(gameMap)
			drawX := x*g.sz - offsetX
			drawY := y*g.sz - offsetY

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

	sz := bank.Size()
	ebitenGame := &EbitenGame{
		touches:      make(map[ebiten.TouchID]*touch),
		sz:           sz,
		rows:         16,
		cols:         20,
		scrWidth:     20 * sz,
		scrHeight:    16 * sz,
		initialScale: sz / 2,
		bank:         bank,
		palette:      palette,
		viewPort:     core.Rect{},
		viewDist:     6,
	}

	ebiten.SetWindowSize(ebitenGame.scrWidth*ebitenGame.initialScale, ebitenGame.scrHeight*ebitenGame.initialScale)
	ebiten.SetWindowPosition(0, 0)
	ebiten.SetWindowTitle("Generic Dungeon Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowIcon([]image.Image{windowIcon})

	// TODO: This is a lot of placeholder for now
	engine.LoadItemFactory()
	game = engine.NewGame()
	ebitenGame.UpdateViewAndFOV()

	if err := ebiten.RunGame(ebitenGame); err != nil {
		log.Fatal(err)
	}
}
