package graphics

import (
	"errors"
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
	name         string
	image        *ebiten.Image
	size         core.Size
	paletteIndex int
}

// Getter for Sprite image
func (s *Sprite) Image() *ebiten.Image {
	return s.image
}

// Getter for Sprite size
func (s *Sprite) Size() core.Size {
	return s.size
}

// Getter for Sprite name
func (s *Sprite) Name() string {
	return s.name
}

// Getter for Sprite palette index
func (s *Sprite) PaletteIndex() int {
	return s.paletteIndex
}

func (s *Sprite) Draw(screen *ebiten.Image, x int, y int, palette color.Palette, dim bool) {
	s.DrawWithColour(screen, x, y, palette, s.paletteIndex, dim)
}

func (s *Sprite) DrawWithColour(screen *ebiten.Image, x int, y int, palette color.Palette, index int, dim bool) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(palette[index])
	if dim {
		op.ColorScale.Scale(0.3, 0.3, 0.3, 1)
	}
	screen.DrawImage(s.image, op)
}

// ====================================================================================================================

// Holds a collection of sprites, indexed by name/id
type SpriteBank struct {
	sprites map[string]*Sprite
	size    int
}

type spriteMetaFile struct {
	Size    int
	Count   int
	Source  string
	Sprites []spriteMetaEntry
}

type spriteMetaEntry struct {
	Name         string
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
		sprites: make(map[string]*Sprite),
		size:    meta.Count,
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
			name:         entry.Name,
			paletteIndex: entry.PaletteIndex,
		}

		spriteBank.sprites[sprite.name] = sprite
		spriteBank.size++
	}

	return spriteBank, nil
}

// Get a sprite from the SpriteBank by name
func (sb *SpriteBank) Sprite(name string) (*Sprite, error) {
	sprite, ok := sb.sprites[name]
	if !ok {
		return nil, errors.New("sprite " + name + " not found")
	}

	return sprite, nil
}

func (sb *SpriteBank) Capacity() int {
	return sb.size
}
