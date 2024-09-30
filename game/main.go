package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand/v2"
	"path"
	"roguelike/core"
	"roguelike/engine"
	"roguelike/game/audio"
	"roguelike/game/graphics"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// These are injected by the build system
var basePath string = "./"
var version string = "0.0.1-alpha_014"

//go:embed icon.png
var iconBytes []byte // Icon for the window is embedded

const (
	PAL_INDEX_WALL   = 0
	PAL_INDEX_FLOOR  = 3
	PAL_INDEX_PLAYER = 10
	VP_ROWS          = 17 // Number of rows of tiles in the viewport, +1 for status bar
	VP_COLS          = 26 // Number of columns of tiles in the viewport
	MAX_EVENT_AGE    = 8  // Max number of events to store
	INITIAL_SCALE    = 4
	ASSETS_DIR       = "assets"
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

// GameStateHander is an interface for processing the game in a given state
type GameStateHander interface {
	Update(heldKeys []ebiten.Key, tappedKeys []ebiten.Key)
	Draw(screen *ebiten.Image)
	Init()
	PassEvent(evt engine.GameEvent)
}

// Implements the ebiten.Game interface
// Holds external state for the rendering & running of the game
type EbitenGame struct {
	game  *engine.Game
	state GameState // nolint

	// Core consts for rendering the window
	spSize    int // const - size of each tile in pixels
	scrWidth  int // const - screen width in pixels
	scrHeight int // const - screen height in pixels

	// Graphics
	bank       *graphics.SpriteBank // Sprite bank holds all the sprites
	paletteSet *graphics.PaletteSet
	palette    color.Palette // Current palette

	// Viewport & FOV
	viewPort core.Rect // The area of the map that is visible
	viewDist int       // View distance in tiles (const)

	// Weird crap for touch controls
	touches  map[ebiten.TouchID]*touch
	touchIDs []ebiten.TouchID
	taps     []tap

	events   []*engine.GameEvent
	eventLog []string

	frameCount int64

	// Audio
	effect *audio.Effects

	// State handlers
	handlers map[GameState]GameStateHander
}

func (g *EbitenGame) Update() error {
	heldKeys := inpututil.AppendPressedKeys(nil)
	tappedKeys := inpututil.AppendJustPressedKeys(nil)

	// Touch controls - figure out taps and touches
	g.taps = handleTaps(g.taps, g.touches)
	handleTouches(g.touchIDs, g.touches)

	// Based on the current state, call the appropriate update handler
	g.handlers[g.state].Update(heldKeys, tappedKeys)

	return nil
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	g.frameCount++
	screen.Fill(color.Black)

	// Based on the current state, call the appropriate draw handler
	g.handlers[g.state].Draw(screen)

	// debug version
	if g.frameCount < 240 {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("v%s", version), g.scrWidth-130, g.scrHeight-14)
	}
}

func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.scrWidth, g.scrHeight
}

func main() {
	log.Printf("GoRogue v%s is starting...", version)

	// Arguments and flags
	var seed uint64
	var disableAudio bool
	flag.Uint64Var(&seed, "seed", 0, "Seed for the game world")
	flag.BoolVar(&disableAudio, "noaudio", false, "Disable audio")
	flag.Parse()

	// Window icon uses embedded bytes
	buf := bytes.NewBuffer(iconBytes)
	icon, _, err := image.Decode(buf)
	if err != nil {
		log.Fatal(err)
	}

	// Load the metadata files for graphics, palettes and sounds
	bank, err := graphics.NewSpriteBank(path.Join(basePath, ASSETS_DIR, "sprites.yaml"))
	if err != nil {
		log.Fatal(err)
	}

	palSet, err := graphics.NewPallettSet(path.Join(basePath, ASSETS_DIR, "palettes.yaml"))
	if err != nil {
		log.Fatal(err)
	}

	sound, err := audio.NewEffects(path.Join(basePath, ASSETS_DIR, "sounds.yaml"), !disableAudio)
	if err != nil {
		log.Fatal(err)
	}

	palette := palSet.Palettes["default"].Colours

	spSize := bank.Size()
	graphics.SetTileSize(spSize)
	ebitenGame := &EbitenGame{
		game:      nil,
		state:     GameStatePlaying,
		touches:   make(map[ebiten.TouchID]*touch),
		spSize:    spSize,
		scrWidth:  VP_COLS * spSize,
		scrHeight: (VP_ROWS + 1) * spSize, // Adds an extra row for status bar

		bank:       bank,
		paletteSet: palSet,
		palette:    palette,
		viewPort:   core.NewRect(0, 0, VP_COLS, VP_ROWS),
		viewDist:   6,
		effect:     sound,
	}

	// Build a map of state handlers for each game state
	ebitenGame.handlers = map[GameState]GameStateHander{
		GameStatePlaying: &PlayingState{
			EbitenGame: ebitenGame,
		},

		GameStateInventory: &InventoryState{
			EbitenGame: ebitenGame,
		},
	}

	ebiten.SetWindowSize(int(float64(ebitenGame.scrWidth)*INITIAL_SCALE), int(float64(ebitenGame.scrHeight)*INITIAL_SCALE))
	ebiten.SetWindowPosition(0, 0)
	ebiten.SetWindowTitle("GoRogue")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowIcon([]image.Image{icon})

	// Either no seed provided or it was invalid, better generate a random one
	if seed == 0 {
		seed = rand.Uint64N(100000000)
		log.Printf("Generated random seed: %d", seed)
	}

	listener := func(e engine.GameEvent) {
		ebitenGame.events = append(ebitenGame.events, &e)
		ebitenGame.eventLog = append(ebitenGame.eventLog, e.Text())

		if e.Type() == engine.EventCreatureKilled {
			ebitenGame.effect.Play("hurt")
		}

		if e.Type() == engine.EventItemPickup {
			ebitenGame.effect.Play("pickup")
		}

		// Pass all events to state handlers
		ebitenGame.handlers[ebitenGame.state].PassEvent(e)
	}

	game := engine.NewGame(basePath+"assets/datafiles", seed, listener)

	// game.AddEventListener()

	ebitenGame.viewPort = game.GetViewPort(VP_COLS, VP_ROWS)
	game.UpdateFOV(ebitenGame.viewDist)

	// IMPORTANT! Link both structs
	ebitenGame.game = game

	// Phew - finally start the ebiten game loop with RunGame
	if err := ebiten.RunGame(ebitenGame); err != nil {
		log.Fatal(err)
	}
}
