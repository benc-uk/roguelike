package main

import (
	"dungeon-run/game/graphics"
	"image/color"
	_ "image/png"
	"log"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// These are injected by the build system
var basePath string = "./"
var version string = "0.0.1-alpha_001"
var cIndex = 0 // TODO: Remove this
var fCount = 0 // TODO: Remove this
var bank *graphics.SpriteBank
var palette color.Palette

const (
	sz           = 12        // Sprite & tile size
	rows         = 12        // Rows of sprites/tiles on screen
	cols         = 16        // Columns of sprites/tiles on screen
	scrWidth     = cols * sz // Screen width in pixels
	scrHeight    = rows * sz // Screen height in pixels
	initialScale = 5         // Initial scale of the window
)

func init() {
	log.Printf("Dungeon Run v%s is starting...", version)

	var err error
	bank, err = graphics.NewSpriteBank(basePath+"assets/sprites/sprites.json", true)
	if err != nil {
		log.Fatal(err)
	}

	err = graphics.LoadPalettes(basePath + "assets/palettes.json")
	if err != nil {
		log.Fatal(err)
	}

	palette, err = graphics.GetRGBPalette("c64")
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	fCount++
	if fCount%10 == 0 {
		cIndex = rand.Intn(len(palette))
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	s1, err := bank.Sprite("Wall")
	if err != nil {
		log.Fatal(err)
	}
	s2, err := bank.Sprite("Slime")
	if err != nil {
		log.Fatal(err)
	}

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*sz), float64(y*sz))

			if x%2 == 0 && y%2 == 0 {
				op.ColorScale.ScaleWithColor(palette[11])
				screen.DrawImage(s1.Image(), op)
			} else {
				c := palette[cIndex]
				op.ColorScale.ScaleWithColor(c)

				screen.DrawImage(s2.Image(), op)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return scrWidth, scrHeight
}

func main() {
	ebiten.SetWindowSize(scrWidth*initialScale, scrHeight*initialScale)
	ebiten.SetWindowTitle("Dungeon Run")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
