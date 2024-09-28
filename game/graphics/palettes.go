package graphics

import (
	"image/color"
	"log"
	"roguelike/core"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type paletteMetaFile struct {
	Palettes map[string]paletteMetaEntry `yaml:"palettes"`
}

type paletteMetaEntry struct {
	Name    string   `yaml:"name"`
	Colours []string `yaml:"colours"`
}

type palette struct {
	Name    string
	Colours color.Palette
}

type PaletteSet struct {
	Palettes map[string]palette
}

func NewPallettSet(metaFile string) (*PaletteSet, error) {
	data, err := core.ReadFile(metaFile)
	if err != nil {
		return nil, err
	}

	paletteMeta := paletteMetaFile{}
	err = yaml.Unmarshal(data, &paletteMeta)
	if err != nil {
		return nil, err
	}

	paletteSet := &PaletteSet{
		Palettes: make(map[string]palette),
	}

	ids := make([]string, 0, len(paletteMeta.Palettes))
	for id, entry := range paletteMeta.Palettes {
		ids = append(ids, id)
		paletteSet.Palettes[id] = palette{
			Name:    entry.Name,
			Colours: GetRGBPalette(entry.Colours),
		}
	}

	log.Printf("Loaded palettes: %s", strings.Join(ids, ", "))

	return paletteSet, nil
}

// GetRGBPalette returns a color.Palette (array of color.RGBA) from the loaded palettes
func GetRGBPalette(pal []string) color.Palette {
	var p color.Palette = make(color.Palette, len(pal))
	for i, c := range pal {
		p[i] = HexStringToRGBA(c)
	}

	return p
}

// Utility function to convert a hex string to a color.RGBA
func HexStringToRGBA(hex string) color.RGBA {
	if hex[0] == '#' {
		hex = hex[1:]
	}

	var r uint64 = 0
	var g uint64 = 0
	var b uint64 = 0
	var a uint64 = 255
	if len(hex) == 6 {
		r, _ = strconv.ParseUint(hex[0:2], 16, 16)
		g, _ = strconv.ParseUint(hex[2:4], 16, 16)
		b, _ = strconv.ParseUint(hex[4:6], 16, 16)
	} else if len(hex) == 8 {
		r, _ = strconv.ParseUint(hex[0:2], 16, 16)
		g, _ = strconv.ParseUint(hex[2:4], 16, 16)
		b, _ = strconv.ParseUint(hex[4:6], 16, 16)
		a, _ = strconv.ParseUint(hex[6:8], 16, 16)
	}

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}
