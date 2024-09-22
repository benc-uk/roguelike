package engine

type GameEvent struct {
	Type   string
	Source entity
	Data   string
}

type EventManager struct {
	// Log of game events
	// log            []GameEvent
	eventListeners []func(GameEvent)
}

// Global event manager, yeah it's a singleton, sue me
var events EventManager

const (
	MAX_LOG_SIZE = 100
)

func init() {
	events = EventManager{}
}

func (em *EventManager) AddEventListener(listener func(GameEvent)) {
	em.eventListeners = append(em.eventListeners, listener)
}

func (em *EventManager) new(eventType string, source entity, data string) {
	e := GameEvent{
		Type:   eventType,
		Source: source,
		Data:   data,
	}

	for _, listener := range em.eventListeners {
		listener(e)
	}

	// // Log the last 100 events
	// em.log = append(em.log, e)
	// if len(em.log) > MAX_LOG_SIZE {
	// 	em.log = events.log[1:]
	// }
}
