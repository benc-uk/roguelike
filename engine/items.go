// nolint
package engine

import (
	"fmt"
	"log"
	"math/rand/v2"
	"roguelike/core"
	"slices"

	"github.com/dop251/goja"
	"gopkg.in/yaml.v3"
)

// ============================================================================
// Item entities are things like weapons, armour, potions etc
// ============================================================================

type Item struct {
	entityBase
	usable        bool          // Can be used by the player
	equipLocation equipLocation // Where the item can be equipped
	dropped       bool          // Previously dropped on the ground
	weight        int           // Weight of the item
	onUseScript   string        // Script to run when the item is used
	rarity        rarity        // Rarity of the item
	consumable    bool          // Consumable items are removed from the player's inventory when used
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

func (i Item) IsEquipment() bool {
	return i.equipLocation != EquipLocationNone
}

func (i Item) Usable() bool {
	return i.usable
}

func (i Item) Weight() int {
	return i.weight
}

func (i Item) Rarity() rarity {
	return i.rarity
}

func (i Item) ItemType() string {
	if i.usable {
		return "Usable / Consumable"
	}

	if i.equipLocation != EquipLocationNone {
		return "Equipment"
	}

	return "Unknown"
}

func (i Item) use(g Game) bool {
	if i.onUseScript == "" || !i.usable {
		return false
	}

	vm := goja.New()

	// Push some globals into the vm for use in the scriptlets
	vm.Set("player", g.player)
	vm.Set("chance", func(chance int) bool {
		return rand.IntN(100) < chance
	})
	vm.Set("rand", func(min, max int) int {
		return rand.IntN(max-min) + min
	})

	result, err := vm.RunString(i.onUseScript)
	if err != nil {
		log.Printf("Error running item script: %s", err)
		events.new(EventMiscMessage, &i, err.Error())
		return false
	}

	if result != nil {
		if msgText, ok := result.Export().(string); ok {
			events.new(EventMiscMessage, &i, msgText)
		} else {
			events.new(EventMiscMessage, &i, "The item's effect is unknown")
		}
	}

	// Items have single use
	g.player.backpack.Remove(&i)

	return true
}

// ===== Item Generator =================================================================================================

// The itemGenerator is a kind of factory for creating items
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
	Consumable    bool   `yaml:"consumable"`
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

	for id, entry := range itemsFile.Items {
		// TODO: Needs to be rarity aware
		gen.genFunctions[id] = func() *Item {
			i := &Item{
				entityBase: entityBase{
					id:         id,
					instanceID: core.RandId(6),
					desc:       entry.Description,
					name:       entry.Name,
					graphicId:  entry.Graphic,
					colour:     entry.Colour,
				},
				usable:      entry.Usable,
				weight:      entry.Weight,
				onUseScript: entry.OnUseScript,
				equipLocation: func() equipLocation {
					switch entry.EquipLocation {
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
						return equipLocationFinger
					case "neck":
						return equipLocationNeck
					default:
						return EquipLocationNone
					}
				}(),
			}

			return i
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

func (gen itemGenerator) createRandomItem(rarity rarity) *Item {
	if len(gen.genFunctions) == 0 {
		return nil
	}

	//	generatorAtRarity := make([]func() *Item, 0)

	// Get a random item, we have to use the keys slice as the map iteration order is random
	id := gen.keys[rng.IntN(len(gen.keys))]
	return gen.createItem(id)
}

// ====== Item Equip Location =================================================================================================

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
	equipLocationFinger
	equipLocationNeck
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
	case equipLocationFinger:
		return "ring"
	case equipLocationNeck:
		return "neck"
	default:
		return "none"
	}
}

// ====== Item Rarity =================================================================================================

type rarity int

const (
	rarityVeryCommon rarity = iota
	rarityCommon
	rarityUncommon
	rarityRare
	rarityVeryRare
	rarityEpic
	rarityLegendary
)

func (r rarity) String() string {
	switch r {
	case rarityVeryCommon:
		return "very common"
	case rarityCommon:
		return "common"
	case rarityUncommon:
		return "uncommon"
	case rarityRare:
		return "rare"
	case rarityVeryRare:
		return "very rare"
	case rarityEpic:
		return "epic"
	case rarityLegendary:
		return "legendary"
	default:
		return "unknown"
	}
}
