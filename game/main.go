package main

import (
	"dungeon-run/game/graphics"
	"image/color"
	_ "image/png"
	"log"

	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

var sb *graphics.SpriteBank
var palette color.Palette
var basePath string = "./"

const (
	sz        = 12
	rows      = 16
	cols      = 22
	scrWidth  = cols * sz
	scrHeight = rows * sz
)

func init() {
	var err error
	sb, err = graphics.NewSpriteBank(basePath + "assets/sprites/sprites_new.json")
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
	}
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// block := sb.Sprite("Block").Image()
	// rock := sb.Sprite("Rock").Image()
	rock, err := sb.Sprite("Wall 1")
	if err != nil {
		log.Fatal(err)
	}
	dirt, err := sb.Sprite("Rat")
	if err != nil {
		log.Fatal(err)
	}

	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(x*sz), float64(y*sz))

			if x%2 == 0 && y%2 == 0 {
				opts.ColorScale.ScaleWithColor(palette[4])
				screen.DrawImage(rock.Image(), opts)
			} else {
				c := palette[rand.Intn(len(palette))]
				opts.ColorScale.ScaleWithColor(c)

				screen.DrawImage(dirt.Image(), opts)
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
