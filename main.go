package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var floorTile *ebiten.Image
var barrelTile *ebiten.Image

const (
	scrWidth  = 320
	scrHeight = 240
)

func init() {
	var err error
	floorTile, _, err = ebitenutil.NewImageFromFile("./assets/tilesets/dungeon/stone.png")
	if err != nil {
		log.Fatal(err)
	}
	barrelTile, _, err = ebitenutil.NewImageFromFile("./assets/tilesets/dungeon/barrel.png")
	if err != nil {
		log.Fatal(err)
	}
}

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// draw 20 x 20 tiles of the image
	for x := 0; x < 20; x++ {
		for y := 0; y < 15; y++ {

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*16), float64(y*16))
			screen.DrawImage(floorTile, op)

			if x == 0 && y == 0 || x == 19 && y == 14 || x == 0 && y == 14 || x == 19 && y == 0 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*16), float64(y*16))
				screen.DrawImage(barrelTile, op)
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
