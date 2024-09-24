package engine

type GameEvent struct {
	Type   string
	Entity entity
	Text   string
	Age    int
}

const (
	EventGameState      = "game_state"
	EventItemPickup     = "item_pickup"
	EventItemUsed       = "item_used"
	EventItemDropped    = "item_dropped"
	EventCreatureKilled = "creature_killed"
)

type EventManager struct {
	// Log of game events
	eventListeners []func(GameEvent)
}

// Global event manager, yeah it's a singleton, sue me
var events EventManager = EventManager{}

func (em *EventManager) AddEventListener(listener func(GameEvent)) {
	em.eventListeners = append(em.eventListeners, listener)
}

func (em *EventManager) new(eventType string, entity entity, data string) {
	e := GameEvent{
		Type:   eventType,
		Entity: entity,
		Text:   data,
		Age:    0,
	}

	for _, listener := range em.eventListeners {
		listener(e)
	}
}
