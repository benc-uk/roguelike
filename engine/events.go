package engine

// ============================================================================
// Events are used to communicate to listeners to changes in the game state
// ============================================================================

type GameEventType string

const (
	EventMiscMessage    = "misc_message"
	EventGameState      = "game_state"
	EventItemPickup     = "item_pickup"
	EventItemMultiple   = "item_pickup_multi"
	EventItemSkipped    = "item_pickup_skipped"
	EventItemUsed       = "item_used"
	EventItemDropped    = "item_dropped"
	EventCreatureKilled = "creature_killed"
	EventPackFull       = "player_pack_full"
)

type GameEvent struct {
	eventType GameEventType
	entity    entity
	text      string
	Age       int
}

func (e GameEvent) Type() GameEventType {
	return e.eventType
}

func (e GameEvent) Text() string {
	return e.text
}

func (e GameEvent) Entity() entity {
	return e.entity
}

func (e GameEvent) SameAs(other *GameEvent) bool {
	if other == nil {
		return false
	}

	if e.entity == nil || other.entity != nil {
		return e.eventType == other.eventType && e.text == other.text
	}

	return e.eventType == other.eventType && e.text == other.text && e.entity == other.entity
}

type eventManager struct {
	// Log of game events
	eventListeners []EventListener
}

type EventListener func(GameEvent)

// Global event manager, yeah it's a singleton, sue me
var events eventManager = eventManager{}

func (em *eventManager) addEventListeners(listener ...EventListener) {
	em.eventListeners = append(em.eventListeners, listener...)
}

func (em *eventManager) new(eventType GameEventType, entity entity, text string) {
	e := GameEvent{
		eventType: eventType,
		entity:    entity,
		text:      text,
		Age:       0,
	}

	for _, listener := range em.eventListeners {
		listener(e)
	}
}
