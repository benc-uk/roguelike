package graphics

import (
	"dungeon-run/game/utils"
	"encoding/json"
	"errors"
	"image/color"
	"strconv"
)

var palettes map[string]palette

type palette struct {
	Name    string
	Colours []string
}

// LoadPalettes loads the palettes from the JSON file, call this before using any other functions in this package
func LoadPalettes(filename string) error {
	data, err := utils.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &palettes)
	if err != nil {
		return err
	}

	return nil
}

// GetPalettes returns the loaded palettes
func GetPalettes() (map[string]palette, error) {
	if palettes == nil {
		return nil, errors.New("palettes not loaded, please call LoadPalettes first")
	}

	return palettes, nil
}

// GetRGBPalette returns a color.Palette (array of color.RGBA) from the loaded palettes
func GetRGBPalette(id string) (color.Palette, error) {
	if palettes == nil {
		return nil, errors.New("palettes not loaded, please call LoadPalettes first")
	}

	pal, ok := palettes[id]
	if !ok {
		return nil, errors.New("palette with id " + id + " not found")
	}

	var p color.Palette = make(color.Palette, len(pal.Colours))
	for i, c := range pal.Colours {
		p[i] = HexStringToRGBA(c)
	}

	return p, nil
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
