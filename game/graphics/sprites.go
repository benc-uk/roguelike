package graphics

import (
	"dungeon-run/game/geom"
	"dungeon-run/game/utils"
	"errors"
	"image"
	"log"
	"path"

	"encoding/json"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Sprite struct {
	name  string
	image *ebiten.Image
	size  geom.Size
}

type SpriteBank struct {
	sprites map[string]*Sprite
	size    int
}

type SpriteMetaFile struct {
	Size    int
	Count   int
	Source  string
	Sprites []SpriteMetaEntry
}

type SpriteMetaEntry struct {
	Name string
	geom.Pos
}

func NewSprite(img *ebiten.Image, name string) *Sprite {
	return &Sprite{
		image: img,
		size:  geom.Size{Width: img.Bounds().Dx(), Height: img.Bounds().Dy()},
		name:  name,
	}
}

func NewSpriteBank(metaFile string, whiteOut bool) (*SpriteBank, error) {
	data, err := utils.ReadFile(metaFile)
	if err != nil {
		return nil, err
	}

	// Parse the JSON data into a SpriteMetaFile struct
	var meta SpriteMetaFile
	err = json.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}

	// Create a new SpriteBank and populate it with sprites
	spriteBank := &SpriteBank{
		sprites: make(map[string]*Sprite),
		size:    meta.Count,
	}

	// Load the source image file
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
		 spriteImg := sheetImg.SubImage(image.Rect(entry.Pos.X, entry.Pos.Y, entry.Pos.X+sz, entry.Pos.Y+sz)).(*ebiten.Image)

		// Logic to white out the sprite or not
		var newImg *ebiten.Image
		if whiteOut {
			newImg = ebiten.NewImage(sz, sz)
			op := &ebiten.DrawImageOptions{}
			op.ColorScale.SetR(255)
			op.ColorScale.SetG(255)
			op.ColorScale.SetB(255)
			newImg.DrawImage(spriteImg, op)
		} else {
			// Clone the image if we don't want to white out the sprite
			newImg = ebiten.NewImageFromImage(spriteImg)
		}

		sprite := NewSprite(newImg, entry.Name)
		spriteBank.AddSprite(sprite)
	}

	return spriteBank, nil
}

func (s *Sprite) Image() *ebiten.Image {
	return s.image
}

func (s *Sprite) Size() geom.Size {
	return s.size
}

func (s *Sprite) Name() string {
	return s.name
}

func (sb *SpriteBank) AddSprite(s *Sprite) {
	sb.sprites[s.name] = s
	sb.size++
}

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
