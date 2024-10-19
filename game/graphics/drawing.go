package graphics

import (
	"roguelike/core"
	"strings"

	"github.com/hajimehoshi/bitmapfont/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/mitchellh/go-wordwrap"
)

var fontFace = text.NewGoXFace(bitmapfont.Face)
var tileSz = 12
var tileSzH = tileSz / 2
var scrRows = 25
var scrCols = 80

func InitGraphics(tileSize int, screenRows, screenCols int) {
	tileSz = tileSize
	tileSzH = tileSize / 2
	scrRows = screenRows
	scrCols = screenCols
}

// Draw a row of text across the screen
func DrawTextRow(screen *ebiten.Image, textStr string, row int) {
	offset := row * tileSz

	vector.DrawFilledRect(screen, 0, float32(offset), 2000, float32(tileSz), BgColour, false)

	op := &text.DrawOptions{}
	op.GeoM.Translate(2, float64(offset-2))
	op.ColorScale.ScaleWithColor(FgColour)
	text.Draw(screen, textStr, fontFace, op)
}

// Draw a filled rectangle with a border around it
func DrawBox(screen *ebiten.Image, row, x int, width, heightRows int) {
	vector.DrawFilledRect(screen, float32(x*tileSz), float32(row*tileSz+tileSzH), float32(width*tileSz), float32(heightRows*tileSz), BgColour, false)
	vector.StrokeRect(screen, float32(x*tileSz), float32(row*tileSz+tileSzH), float32(width*tileSz), float32(heightRows*tileSz), 2, FgColour, false)
}

// Draw simple text based dialog box
func DrawDialogBox(screen *ebiten.Image, title string, body string) {
	old := BgColour
	BgColour = ColourDialog

	// 15 is a magic number that seems to work well for the font face
	body = wordwrap.WrapString(body, uint(scrCols+15))

	bodyLines := strings.Split(body, "\n")
	bodyLineCount := len(bodyLines)

	height := bodyLineCount + 3
	topOffset := (scrRows - height) / 2
	width := scrCols - 4

	// Draw the main outlines of the box, with a border and title bar
	DrawBox(screen, topOffset, 2, width, 2)
	DrawBox(screen, topOffset+2, 2, width, height-2)

	// Draw the title text
	BgColour = ColourTrans
	DrawTextRow(screen, "     "+title, topOffset+1)

	// Draw the body text
	for i, line := range bodyLines {
		DrawTextRow(screen, "     "+line, topOffset+i+3)
	}

	BgColour = old
}

// Draw text that wraps to the next line when it reaches the wrapWidth
func DrawWrappedText(screen *ebiten.Image, text string, startRow int, pad int, wrapWidth int) {
	text = wordwrap.WrapString(text, uint(wrapWidth))
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		DrawTextRow(screen, core.MakeStr(pad, " ")+line, startRow+i)
	}
}
