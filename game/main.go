package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand/v2"
	"os"
	"roguelike/core"
	"roguelike/engine"
	"roguelike/game/graphics"
	"slices"
	"strconv"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// These are injected by the build system
var basePath string = "./"
var version string = "0.0.1-alpha_011"

//go:embed icon.png
var iconBytes []byte // Icon for the window is embedded

// Left global for now, could be refactored
var game *engine.Game

const (
	PAL_INDEX_WALL   = 0
	PAL_INDEX_FLOOR  = 3
	PAL_INDEX_PLAYER = 10
	VP_ROWS          = 17 // Number of rows of tiles in the viewport, +1 for status bar
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

	events      []*engine.GameEvent
	eventLog    []string
	statusText  string
	playerLeft  bool
	frameCount  int64
	delayFrames int
}

func (g *EbitenGame) Update() error {
	// One closing save events to file, one line per event
	if ebiten.IsWindowBeingClosed() {
		eventFile, _ := os.OpenFile("events.txt", os.O_CREATE|os.O_WRONLY, 0644)
		_ = eventFile.Truncate(0)
		for _, evtText := range g.eventLog {
			_, _ = eventFile.WriteString(evtText + "\n")
		}
	}

	p := game.Player()
	g.statusText = fmt.Sprintf("%s    ♥%d/%d   ⌘%d   ▼%d", p.Name(), p.CurrentHP(), p.MaxHP(), p.Exp(), p.Level())

	var move *engine.MoveAction
	pressedKeys := inpututil.AppendPressedKeys(nil)
	justPressedKeys := inpututil.AppendJustPressedKeys(nil)

	// Touch controls - figure out taps and touches
	g.taps = handleTaps(g.taps, g.touches)
	handleTouches(g.touchIDs, g.touches)

	// Loop over taps (there should only be one for reasons)
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

	// Held keys require a delay before moving the player
	for _, key := range pressedKeys {
		if slices.Contains(controls["up"], key) && inpututil.KeyPressDuration(key) > 20 {
			move = engine.NewMoveAction(core.DirNorth)
		}
		if slices.Contains(controls["down"], key) && inpututil.KeyPressDuration(key) > 20 {
			move = engine.NewMoveAction(core.DirSouth)
		}
		if slices.Contains(controls["left"], key) && inpututil.KeyPressDuration(key) > 20 {
			move = engine.NewMoveAction(core.DirWest)
			g.playerLeft = true
		}
		if slices.Contains(controls["right"], key) && inpututil.KeyPressDuration(key) > 20 {
			move = engine.NewMoveAction(core.DirEast)
			g.playerLeft = false
		}
	}

	// Tapped keys (just pressed) reset the delayFrames and move the player immediately
	for _, key := range justPressedKeys {
		if slices.Contains(controls["up"], key) {
			move = engine.NewMoveAction(core.DirNorth)
			g.delayFrames = 0
		}
		if slices.Contains(controls["down"], key) {
			move = engine.NewMoveAction(core.DirSouth)
			g.delayFrames = 0
		}
		if slices.Contains(controls["left"], key) {
			move = engine.NewMoveAction(core.DirWest)
			g.playerLeft = true
			g.delayFrames = 0
		}
		if slices.Contains(controls["right"], key) {
			move = engine.NewMoveAction(core.DirEast)
			g.playerLeft = false
			g.delayFrames = 0
		}
	}

	// This stops the whole game from running too fast
	if g.delayFrames > 0 {
		g.delayFrames--
		return nil
	}

	if move != nil {
		result := move.Execute(*game)
		if !result.Success {
			return nil
		}

		// We translate the energy spent into frames to delay the game
		if result.EnergySpent > 0 {
			g.delayFrames = result.EnergySpent
		}

		g.viewPort = game.GetViewPort(VP_COLS, VP_ROWS)
		game.UpdateFOV(g.viewDist)

		// Handle events and age them
		for _, e := range g.events {
			e.Age++
		}

		// Remove old events
		for i := 0; i < len(g.events); i++ {
			e := g.events[i]
			if e.Age > MAX_EVENT_AGE {
				g.events = append(g.events[:i], g.events[i+1:]...)
			}
		}
	}

	return nil
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	g.frameCount++

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
				g.bank.Sprite("wall").Draw(screen, drawX, drawY, g.palette[PAL_INDEX_WALL], appear.InFOV, false, false)
				continue
			}

			// Draw the player
			if x == p.Pos().X && y == p.Pos().Y {
				g.bank.Sprite("player").Draw(screen, drawX, drawY, g.palette[PAL_INDEX_PLAYER], appear.InFOV, g.playerLeft, false)
				continue
			}

			if appear.Graphic == "floor" {
				g.bank.Sprite("floor").Draw(screen, drawX, drawY, g.palette[PAL_INDEX_FLOOR], appear.InFOV, false, false)
				continue
			}

			// Then items/monsters/stuff
			itemSprite := g.bank.Sprite(appear.Graphic)
			if itemSprite != nil {
				itemSprite.Draw(screen, drawX, drawY, colour, appear.InFOV, false, false)
				continue
			}
		}
	}

	// Draw the status bar, it was at row VP_ROWS-1 but we added a row for the status bar
	graphics.DrawTextRow(screen, g.statusText, VP_ROWS, color.RGBA{0x10, 0x50, 0x10, 0xff})

	for i, e := range g.events {
		graphics.DrawTextRow(screen, e.Text, i, color.RGBA{0x00, 0x00, 0x30, 0x30})
	}
}

func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.scrWidth, g.scrHeight
}

func main() {
	log.Printf("GoRogue v%s is starting...", version)

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
		scrHeight: (VP_ROWS + 1) * spSize, // Adds an extra row for status bar

		bank:     bank,
		palette:  palette,
		viewPort: core.NewRect(0, 0, VP_COLS, VP_ROWS),
		viewDist: 6,
	}

	ebiten.SetWindowSize(int(float64(ebitenGame.scrWidth)*INITIAL_SCALE), int(float64(ebitenGame.scrHeight)*INITIAL_SCALE))
	ebiten.SetWindowPosition(0, 0)
	ebiten.SetWindowTitle("GoRogue")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowIcon([]image.Image{icon})

	// Seed the game world from the command line or not
	var seed uint64
	if len(os.Args) > 1 && os.Args[1] != "" {
		seed, err = strconv.ParseUint(os.Args[1], 10, 64)
		if err != nil {
			log.Printf("Error parsing seed: %s", err)
		} else {

			log.Printf("Using provided game world seed: %d", seed)
		}
	}

	// Either no seed provided or it was invalid, better generate a random one
	if seed == 0 {
		seed = rand.Uint64N(100000000)
		log.Printf("Generated random seed: %d", seed)
	}

	game = engine.NewGame(basePath+"assets/datafiles", seed)

	game.AddEventListener(func(e engine.GameEvent) {
		ebitenGame.events = append(ebitenGame.events, &e)
		ebitenGame.eventLog = append(ebitenGame.eventLog, e.Text)
	})

	ebitenGame.viewPort = game.GetViewPort(VP_COLS, VP_ROWS)
	game.UpdateFOV(ebitenGame.viewDist)

	levelText := fmt.Sprintf("You are on level %d of %s", game.Map().Depth(), game.Map().Description())
	ebitenGame.events = append(ebitenGame.events, &engine.GameEvent{Type: "game_state", Text: "Version " + version})
	ebitenGame.events = append(ebitenGame.events, &engine.GameEvent{Type: "game_state", Text: levelText})

	// Phew - finally start the ebiten game loop with RunGame
	if err := ebiten.RunGame(ebitenGame); err != nil {
		log.Fatal(err)
	}
}
