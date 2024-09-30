package controls

import "github.com/hajimehoshi/ebiten/v2"

type control int

const (
	Up control = iota
	Down
	Left
	Right
	Inventory
	Drop
	Get
	Escape
	Select
)

// TODO: Move this to some sort of config file
var controls = map[control][]ebiten.Key{
	Up:        {ebiten.KeyW, ebiten.KeyUp},
	Down:      {ebiten.KeyS, ebiten.KeyDown},
	Left:      {ebiten.KeyA, ebiten.KeyLeft},
	Right:     {ebiten.KeyD, ebiten.KeyRight},
	Inventory: {ebiten.KeyI},
	Drop:      {ebiten.KeyD},
	Get:       {ebiten.KeyG},
	Escape:    {ebiten.KeyEscape},
	Select:    {ebiten.KeyEnter, ebiten.KeySpace},
}

func (c control) Keys() []ebiten.Key {
	return controls[c]
}

func (c control) IsKeys(keys []ebiten.Key) bool {
	for _, key := range keys {
		for _, controlKey := range controls[c] {
			if key == controlKey {
				return true
			}
		}
	}

	return false
}

func (c control) IsKey(key ebiten.Key) bool {
	for _, controlKey := range controls[c] {
		if key == controlKey {
			return true
		}
	}

	return false
}
