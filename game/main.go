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
	bank, err = graphics.NewSpriteBank(basePath + "assets/sprites.json")
	if err != nil {
		log.Fatal(err)
	}

	err = graphics.LoadPalettes(basePath + "assets/palettes.yaml")
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

type Game struct {
	touches  map[ebiten.TouchID]*touch
	touchIDs []ebiten.TouchID
	taps     []tap
}

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

		if tap.X < scrWidth/4 {
			move := engine.NewMoveAction(core.West)
			move.Execute(game.Player(), game.Map())

		} else if tap.X > scrWidth/4*3 {
			move := engine.NewMoveAction(core.East)
			move.Execute(game.Player(), game.Map())
		}

		if tap.Y < scrHeight/4 {
			move := engine.NewMoveAction(core.North)
			move.Execute(game.Player(), game.Map())
		} else if tap.Y > scrHeight/4*3 {
			move := engine.NewMoveAction(core.South)
			move.Execute(game.Player(), game.Map())
		}
	}

	game.UpdateFOV(*game.Player(), 6)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	gameMap := game.Map()
	p := game.Player()
	wallPalIndex := 0
	floorPalIndex := 3
	playerPalIndex := 10

	// viewPort is the area of the map that is visible centered on the player
	viewPort := core.NewRect(p.X-cols/2, p.Y-rows/2, cols, rows)

	if viewPort.X < 0 {
		viewPort.X = 0
	}
	if viewPort.Y < 0 {
		viewPort.Y = 0
	}
	if viewPort.X+viewPort.Width > gameMap.Size().Width {
		viewPort.X = gameMap.Size().Width - viewPort.Width
	}
	if viewPort.Y+viewPort.Height > gameMap.Size().Height {
		viewPort.Y = gameMap.Size().Height - viewPort.Height
	}

	offsetX := viewPort.X * sz
	offsetY := viewPort.Y * sz

	// Draw the map
	for x := viewPort.X; x < viewPort.Width+viewPort.X; x++ {
		for y := viewPort.Y; y < viewPort.Height+viewPort.Y; y++ {
			tile := gameMap.Tile(x, y)
			appear := tile.GetAppearance(gameMap)
			drawX := x*sz - offsetX
			drawY := y*sz - offsetY

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
			colour := palette[palIndex]

			// Walls
			if appear.Graphic == "wall" {
				bank.Sprite("wall").Draw(screen, drawX, drawY, palette[wallPalIndex], appear.InFOV)
				continue
			}

			// Draw the player in the center of the screen
			if x == p.X && y == p.Y {
				bank.Sprite("player").Draw(screen, drawX, drawY, palette[playerPalIndex], appear.InFOV)
				continue
			}

			if appear.Graphic == "floor" {
				bank.Sprite("floor").Draw(screen, drawX, drawY, palette[floorPalIndex], appear.InFOV)
				continue
			}

			// Then items/monsters/stuff
			itemSprite := bank.Sprite(appear.Graphic)
			if itemSprite != nil {
				itemSprite.Draw(screen, drawX, drawY, colour, appear.InFOV)
				continue
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

	ebitenGame := &Game{
		touches: make(map[ebiten.TouchID]*touch),
	}

	if err := ebiten.RunGame(ebitenGame); err != nil {
		log.Fatal(err)
	}
}
