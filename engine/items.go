// nolint
package engine

import (
	"fmt"
	"roguelike/core"
	"slices"

	"gopkg.in/yaml.v3"
)

// ============================================================================
// Item entities are things like weapons, armour, potions etc
// ============================================================================

type equipLocation int

// TODO: Player can equip more than one ring!
const (
	EquipLocationNone equipLocation = iota
	equipLocationWeapon
	equipLocationMissile
	equipLocationBody
	equipLocationShield
	equipLocationHead
	equipLocationFeet
	equipLocationHands
	equipLocationRing
	equipLocationAmulet
)

func (el equipLocation) String() string {
	switch el {
	case EquipLocationNone:
		return "none"
	case equipLocationWeapon:
		return "weapon"
	case equipLocationMissile:
		return "missile"
	case equipLocationBody:
		return "body"
	case equipLocationShield:
		return "shield"
	case equipLocationHead:
		return "head"
	case equipLocationFeet:
		return "feet"
	case equipLocationHands:
		return "hands"
	case equipLocationRing:
		return "ring"
	case equipLocationAmulet:
		return "amulet"
	default:
		return "unknown"
	}
}

type Item struct {
	entityBase
	usable        bool          // Can be used by the player
	equipLocation equipLocation // Where the item can be equipped
	dropped       bool          // Previously dropped on the ground
	weight        int           // Weight of the item
	onUseScript   string        // Script to run when the item is used
}

func (i Item) Type() entityType {
	return entityTypeItem
}

func (i Item) String() string {
	return fmt.Sprintf("item_%s_%s at %s", i.name, i.id, i.pos)
}

func (i Item) BlocksLOS() bool {
	return false
}

func (i Item) BlocksMove() bool {
	return false
}

func (i Item) EquipLocation() equipLocation {
	return i.equipLocation
}

// ===== Item Generator =================================================================================================

type itemGenerator struct {
	genFunctions map[string](func() *Item)
	keys         []string
}

type yamlItem struct {
	Description   string `yaml:"description"`
	Name          string `yaml:"name"`
	Graphic       string `yaml:"graphic"`
	Colour        string `yaml:"colour"`
	Usable        bool   `yaml:"usable"`
	EquipLocation string `yaml:"equipLocation"`
	Weight        int    `yaml:"weight"`
	OnUseScript   string `yaml:"onUseScript"`
}

type yamlItemsFile struct {
	Items map[string]yamlItem `yaml:"items"`
}

func newItemGenerator(dataFile string) (*itemGenerator, error) {
	data, err := core.ReadFile(dataFile)
	if err != nil {
		return nil, err
	}

	var itemsFile yamlItemsFile
	err = yaml.Unmarshal(data, &itemsFile)
	if err != nil {
		return nil, err
	}

	gen := itemGenerator{
		genFunctions: make(map[string](func() *Item)),
		keys:         make([]string, 0),
	}

	for id, item := range itemsFile.Items {
		gen.genFunctions[id] = func() *Item {
			return &Item{
				entityBase: entityBase{
					id:         id,
					instanceID: core.RandId(6),
					desc:       item.Description,
					name:       item.Name,
					graphicId:  item.Graphic,
					colour:     item.Colour,
				},
				usable:      item.Usable,
				weight:      item.Weight,
				onUseScript: item.OnUseScript,
				equipLocation: func() equipLocation {
					switch item.EquipLocation {
					case "weapon":
						return equipLocationWeapon
					case "missile":
						return equipLocationMissile
					case "body":
						return equipLocationBody
					case "shield":
						return equipLocationShield
					case "head":
						return equipLocationHead
					case "feet":
						return equipLocationFeet
					case "hands":
						return equipLocationHands
					case "ring":
						return equipLocationRing
					case "amulet":
						return equipLocationAmulet
					default:
						return EquipLocationNone
					}
				}(),
			}
		}

		gen.keys = append(gen.keys, id)
	}

	// Sort the keys as the map iteration order above is random
	slices.Sort(gen.keys)

	return &gen, nil
}

// nolint
func (gen itemGenerator) createItem(id string) *Item {
	itemFunc, ok := gen.genFunctions[id]
	if !ok {
		return nil
	}

	// Create the item by invoking the generation function
	return itemFunc()
}

func (gen itemGenerator) createRandomItem() *Item {
	if len(gen.genFunctions) == 0 {
		return nil
	}

	// Get a random item, we have to use the keys slice as the map iteration order is random
	id := gen.keys[rng.IntN(len(gen.keys))]
	return gen.createItem(id)
}
