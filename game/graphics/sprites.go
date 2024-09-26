package graphics

import (
	"image"
	"log"
	"path"
	"roguelike/core"

	"encoding/json"
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Represents a single named sprite image
type Sprite struct {
	id           string
	image        *ebiten.Image
	size         core.Size
	paletteIndex int
}

func (s *Sprite) Draw(screen *ebiten.Image, x int, y int, colour color.Color, inFOV bool, flipX bool, flipY bool) {
	if s == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear

	if flipX {
		op.GeoM.Scale(-1, 1)
		x += s.size.Width
	}
	if flipY {
		op.GeoM.Scale(1, -1)
		y += s.size.Height
	}

	op.GeoM.Translate(float64(x), float64(y))

	op.ColorScale.ScaleWithColor(colour)

	if !inFOV {
		op.ColorScale.ScaleAlpha(0.5)
	}

	screen.DrawImage(s.image, op)
}

// ====================================================================================================================

// Holds a collection of sprites, indexed by name/id
type SpriteBank struct {
	sprites  map[string]*Sprite
	capacity int
	size     int
}

type spriteMetaFile struct {
	Size    int
	Count   int
	Source  string
	Sprites []spriteMetaEntry
}

type spriteMetaEntry struct {
	Id           string
	PaletteIndex int
	core.Pos
}

// Create a new SpriteBank from a JSON meta file and a source image file
func NewSpriteBank(metaFile string) (*SpriteBank, error) {
	data, err := core.ReadFile(metaFile)
	if err != nil {
		return nil, err
	}

	// Parse the JSON data into a SpriteMetaFile struct
	var meta spriteMetaFile
	err = json.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}

	// Create a new SpriteBank and populate it with sprites
	spriteBank := &SpriteBank{
		sprites:  make(map[string]*Sprite),
		capacity: meta.Count,
		size:     meta.Size,
	}

	// Load the source image file, relative to the meta file
	metaDir := path.Dir(metaFile)
	imgPath := path.Join(metaDir, meta.Source)

	log.Printf("Loading sprite sheet from: %s", imgPath)

	sheetImg, _, err := ebitenutil.NewImageFromFile(imgPath)
	if err != nil {
		return nil, err
	}

	sz := meta.Size
	for _, entry := range meta.Sprites {
		// Sub image inside the sprite sheet where the sprite is located
		subImage := sheetImg.SubImage(image.Rect(entry.Pos.X, entry.Pos.Y, entry.Pos.X+sz, entry.Pos.Y+sz)).(*ebiten.Image)

		// Logic to white out the sprite or not, used for monochrome sprites
		var spriteImg *ebiten.Image
		if entry.PaletteIndex > -1 {
			spriteImg = ebiten.NewImage(sz, sz)
			op := &ebiten.DrawImageOptions{}
			op.ColorScale.SetR(255)
			op.ColorScale.SetG(255)
			op.ColorScale.SetB(255)
			spriteImg.DrawImage(subImage, op)
		} else {
			// HACK: Clone the image, may not be needed
			spriteImg = ebiten.NewImageFromImage(subImage)
		}

		sprite := &Sprite{
			image:        spriteImg,
			size:         core.Size{Width: spriteImg.Bounds().Dx(), Height: spriteImg.Bounds().Dy()},
			id:           entry.Id,
			paletteIndex: entry.PaletteIndex,
		}

		spriteBank.sprites[sprite.id] = sprite
		spriteBank.capacity++
	}

	return spriteBank, nil
}

// Get a sprite from the SpriteBank by name, can return nil
func (sb *SpriteBank) Sprite(name string) *Sprite {
	sprite, ok := sb.sprites[name]
	if !ok {
		return nil
	}

	return sprite
}

func (sb *SpriteBank) Capacity() int {
	return sb.capacity
}

func (sb *SpriteBank) Size() int {
	return sb.size
}
