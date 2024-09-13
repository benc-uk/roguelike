package engine

// TODO: This is placeholder code
func LoadItemFactory() {
	itemFactory = make(map[string](func() *Item))

	itemFactory["sword"] = func() *Item {
		return &Item{
			entityBase: entityBase{
				id:        "sword",
				desc:      "A rusty sword",
				graphicId: "sword",
				colour:    "2",
			},
			consumable: false,
		}
	}

	itemFactory["potion"] = func() *Item {
		return &Item{
			entityBase: entityBase{
				id:        "potion",
				desc:      "A refreshing looking blue potion",
				graphicId: "potion",
				colour:    "9",
			},
			consumable: true,
		}
	}

	itemFactory["potion_poison"] = func() *Item {
		return &Item{
			entityBase: entityBase{
				id:        "potion_poison",
				desc:      "A bubbling pink potion",
				graphicId: "potion",
				colour:    "11",
			},
			consumable: true,
		}
	}

	itemFactory["door"] = func() *Item {
		return &Item{
			entityBase: entityBase{
				id:         "door",
				desc:       "A sturdy wooden door",
				graphicId:  "door",
				blocksMove: false,
				blocksLOS:  true,
				colour:     "5",
			},
			consumable: false,
		}
	}

	itemFactory["rat"] = func() *Item {
		return &Item{
			entityBase: entityBase{
				id:         "rat",
				desc:       "A small, scurrying rat",
				graphicId:  "rat",
				blocksMove: true,
				colour:     "4",
			},
			consumable: false,
		}
	}
}
