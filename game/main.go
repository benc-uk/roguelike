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
var cIndex = 0
var fCount = 0
var sb *graphics.SpriteBank
var palette color.Palette

const (
	sz        = 12
	rows      = 16
	cols      = 22
	scrWidth  = cols * sz
	scrHeight = rows * sz
)

func init() {
	log.Printf("Dungeon Run v%s is starting...", version)

	var err error
	sb, err = graphics.NewSpriteBank(basePath+"assets/sprites/sprites.json", true)
	if err != nil {
		log.Fatal(err)
	}

	// Set up 16 color palette
	palette = color.Palette{
		color.RGBA{0, 255, 0, 255},
		color.RGBA{29, 43, 83, 255},
		color.RGBA{126, 37, 83, 255},
		color.RGBA{26, 37, 255, 255},
		color.RGBA{128, 128, 128, 255},
		color.RGBA{0, 0, 0, 255},
		color.RGBA{255, 0, 77, 255},
		color.RGBA{255, 163, 0, 255},
		color.RGBA{255, 240, 36, 255},
		color.RGBA{0, 231, 86, 255},
		color.RGBA{41, 173, 255, 255},
		color.RGBA{131, 118, 156, 255},
		color.RGBA{255, 119, 168, 255},
		color.RGBA{255, 204, 170, 255},
		color.RGBA{255, 255, 255, 255},
		color.RGBA{0, 122, 76, 255},
	}
}

type Game struct{}

func (g *Game) Update() error {
	fCount++
	if fCount%4 == 0 {
		cIndex = rand.Intn(len(palette))
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	s1, err := sb.Sprite("Wall")
	if err != nil {
		log.Fatal(err)
	}
	s2, err := sb.Sprite("Slime")
	if err != nil {
		log.Fatal(err)
	}

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(x*sz), float64(y*sz))

			if x%2 == 0 && y%2 == 0 {
				opts.ColorScale.ScaleWithColor(palette[4])
				screen.DrawImage(s1.Image(), opts)
			} else {
				c := palette[cIndex]
				opts.ColorScale.ScaleWithColor(c)

				screen.DrawImage(s2.Image(), opts)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return scrWidth, scrHeight
}

func main() {
	ebiten.SetWindowSize(scrWidth*4, scrHeight*4)
	ebiten.SetWindowTitle("Dungeon Run")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
