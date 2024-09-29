package graphics

import (
	"image/color"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var fontFace = text.NewGoXFace(bitmapfont.Face)
var tileSz = 12
var tileSzH = tileSz / 2

func SetTileSize(h int) {
	tileSz = h
	tileSzH = h / 2
}

func DrawTextRow(screen *ebiten.Image, textStr string, row int, bgCol color.RGBA) {
	offset := row * tileSz

	vector.DrawFilledRect(screen, 0, float32(offset), 2000, float32(tileSz), bgCol, false)

	op := &text.DrawOptions{}
	op.GeoM.Translate(2, float64(offset-2))
	text.Draw(screen, textStr, fontFace, op)
}

// Use box drawing characters to draw a text box
func DrawTextBox(screen *ebiten.Image, row, x int, width, heightRows int, bgCol color.RGBA) {
	vector.DrawFilledRect(screen, float32(x+tileSzH), float32(row*tileSz+tileSzH), float32(width*tileSz), float32(heightRows*tileSz), bgCol, false)
	vector.StrokeRect(screen, float32(x+tileSzH), float32(row*tileSz+tileSzH), float32(width*tileSz), float32(heightRows*tileSz), 2, color.White, false)
}
