package graphics

import (
	"dungeon-run/game/geom"
	_ "embed"
	"errors"
	"image"
	"log"
	"path"

	"encoding/json"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// HACK: Temp fix for not being able to load the JSON file from WASM
//
//go:embed dungeon.json
var dungeonJSON []byte

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

func NewSpriteBank(metaFile string) (*SpriteBank, error) {
	// Read the metaFile from disk
	// data, err := os.ReadFile(metaFile)
	// if err != nil {
	// 	return nil, err
	// }
	data := dungeonJSON

	// Parse the JSON data into a SpriteMetaFile struct
	var meta SpriteMetaFile
	err := json.Unmarshal(data, &meta)
	if err != nil {
		return nil, err
	}

	log.Printf("Creating sprite bank with %d sprites from %s\n", meta.Count, meta.Source)

	// Create a new SpriteBank and populate it with sprites
	spriteBank := &SpriteBank{
		sprites: make(map[string]*Sprite),
		size:    meta.Count,
	}

	// Load the source image file
	metaDir := path.Dir(metaFile)
	sheetImg, _, err := ebitenutil.NewImageFromFile(path.Join(metaDir, meta.Source))
	if err != nil {
		return nil, err
	}

	sz := meta.Size
	for _, entry := range meta.Sprites {
		// Sub image inside the sprite sheet where the sprite is located
		spriteImg := sheetImg.SubImage(image.Rect(entry.Pos.X, entry.Pos.Y, entry.Pos.X+sz, entry.Pos.Y+sz)).(*ebiten.Image)

		// TODO: We make a copy here just in case, it might be better to just use the subimage directly
		sprite := NewSprite(ebiten.NewImageFromImage(spriteImg), entry.Name)
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
