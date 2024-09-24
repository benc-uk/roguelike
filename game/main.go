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

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// These are injected by the build system
var basePath string = "./"
var version string = "0.0.1-alpha_008"

//go:embed icon.png
var iconBytes []byte // Icon for the window is embedded

// Left global for now, could be refactored
var game *engine.Game

const (
	PAL_INDEX_WALL   = 0
	PAL_INDEX_FLOOR  = 3
	PAL_INDEX_PLAYER = 10
	VP_ROWS          = 18 // Number of rows of tiles in the viewport
	VP_COLS          = 26 // Number of columns of tiles in the viewport
	MAX_EVENT_AGE    = 6  // Max number of events to store
	INITIAL_SCALE    = 4
)

type GameState int

const (
	GameStateTitle GameState = iota
	GameStatePlaying
	GameStateInventory
	GameStateCharacter
	GameStateGameOver
	GameStatePlayerGen
)

// Implements the ebiten.Game interface
// Holds external state for the rendering & running of the game
type EbitenGame struct {
	state GameState // nolint

	// Core consts for rendering the window
	spSize    int // const - size of each tile in pixels
	scrWidth  int // const - screen width in pixels
	scrHeight int // const - screen height in pixels

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

	events     []*engine.GameEvent
	statusText string
}

func (g *EbitenGame) Update() error {
	p := game.Player()

	var move *engine.MoveAction
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		move = engine.NewMoveAction(core.DirNorth)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		move = engine.NewMoveAction(core.DirSouth)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		move = engine.NewMoveAction(core.DirWest)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		move = engine.NewMoveAction(core.DirEast)
	}

	// Touch controls - figure out taps and touches
	g.taps = handleTaps(g.taps, g.touches)
	handleTouches(g.touchIDs, g.touches)

	// Loop over taps (there should only be one)
	for _, tap := range g.taps {
		if tap.X < g.scrWidth/4 {
			move = engine.NewMoveAction(core.DirWest)
		} else if tap.X > g.scrWidth/4*3 {
			move = engine.NewMoveAction(core.DirEast)
		}

		if tap.Y < g.scrHeight/4 {
			move = engine.NewMoveAction(core.DirNorth)
		} else if tap.Y > g.scrHeight/4*3 {
			move = engine.NewMoveAction(core.DirSouth)
		}
	}

	if move != nil {
		move.Execute(p, game.Map())
		g.viewPort = game.GetViewPort(VP_COLS, VP_ROWS)
		game.UpdateFOV(g.viewDist)

		// Handle events and age them
		for _, e := range g.events {
			e.Age++
		}

		// Remove old events
		for i, e := range g.events {
			if e.Age > MAX_EVENT_AGE {
				g.events = append(g.events[:i], g.events[i+1:]...)
			}
		}
	}

	g.statusText = "❤️18/45   $5   ▼1"

	return nil
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
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
				g.bank.Sprite("wall").Draw(screen, drawX, drawY, g.palette[PAL_INDEX_WALL], appear.InFOV)
				continue
			}

			// Draw the player
			if x == p.X && y == p.Y {
				g.bank.Sprite("player").Draw(screen, drawX, drawY, g.palette[PAL_INDEX_PLAYER], appear.InFOV)
				continue
			}

			if appear.Graphic == "floor" {
				g.bank.Sprite("floor").Draw(screen, drawX, drawY, g.palette[PAL_INDEX_FLOOR], appear.InFOV)
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

	graphics.DrawTextRow(screen, g.statusText, VP_ROWS-1, color.RGBA{0x10, 0x50, 0x10, 0xff})

	for i, e := range g.events {
		graphics.DrawTextRow(screen, e.Text, i, color.RGBA{0x00, 0x00, 0x30, 0x30})
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

	// Using hajimehoshi/bitmapfont/v3 for now
	graphics.SetFontFace(text.NewGoXFace(bitmapfont.Face))

	spSize := bank.Size()
	ebitenGame := &EbitenGame{
		state:     GameStateTitle,
		touches:   make(map[ebiten.TouchID]*touch),
		spSize:    spSize,
		scrWidth:  VP_COLS * spSize,
		scrHeight: VP_ROWS * spSize,

		bank:     bank,
		palette:  palette,
		viewPort: core.NewRect(0, 0, VP_COLS, VP_ROWS),
		viewDist: 6,
	}

	ebiten.SetWindowSize(int(float64(ebitenGame.scrWidth)*INITIAL_SCALE), int(float64(ebitenGame.scrHeight)*INITIAL_SCALE))
	ebiten.SetWindowPosition(0, 0)
	ebiten.SetWindowTitle("Generic Dungeon Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowIcon([]image.Image{icon})

	game = engine.NewGame(basePath + "assets/datafiles")
	game.AddEventListener(func(e engine.GameEvent) {
		ebitenGame.events = append(ebitenGame.events, &e)
	})
	ebitenGame.viewPort = game.GetViewPort(VP_COLS, VP_ROWS)
	game.UpdateFOV(ebitenGame.viewDist)

	ebitenGame.events = append(ebitenGame.events, &engine.GameEvent{Type: "game_state", Text: "You have entered level 1!"})

	// HACK: Removed for now to test map generation
	if err := ebiten.RunGame(ebitenGame); err != nil {
		log.Fatal(err)
	}
}
