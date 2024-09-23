package graphics

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var fontFace text.Face

func SetFontFace(face text.Face) {
	fontFace = face
}

func DrawTextRow(screen *ebiten.Image, textStr string, row int, bgCol color.RGBA) {
	const rowH = 12
	offset := row * rowH

	vector.DrawFilledRect(screen, 0, float32(offset), 2000, rowH, bgCol, false)

	op := &text.DrawOptions{}
	op.GeoM.Translate(2, float64(offset-1))
	text.Draw(screen, textStr, fontFace, op)
}
