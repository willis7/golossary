package handlers

import "sync"

// DefaultEventMux is the default EventMux used by Serve.
var DefaultEventMux = &defaultEventMux

var defaultEventMux EventMux

type Handler interface{}

// EventMux is an Event multiplexer.
// It matches the event type of each incoming event against a list of registered
// events and calls the handler for the event that
// most closely matches the event.
type EventMux struct {
	mu sync.RWMutex
	m  map[string]muxEntry
}

type muxEntry struct {
	explicit bool
	h        Handler
	event    string
}

// NewEventsMux allocates and returns a new EventsMux.
func NewEventsMux() *EventMux { return new(EventMux) }

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (mux *EventMux) Handle(event string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if event == "" {
		panic("events: invalid event " + event)
	}
	if handler == nil {
		panic("events: nil handler")
	}
	if mux.m[event].explicit {
		panic("events: multiple registrations for " + event)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	mux.m[event] = muxEntry{explicit: true, h: handler, event: event}
}

// Find a handler on a handler map given an event string.
// Most-specific (longest) pattern wins.
func (mux *EventMux) match(event string) (h Handler, pattern string) {
	// Check for exact match first.
	v, ok := mux.m[event]
	if ok {
		return v.h, v.event
	} else {
		panic("events: missing handler for " + event)
	}
}

func eventMatch(event string, inputEvent string) bool {
	if len(event) == 0 {
		// should not happen
		return false
	}

	return inputEvent == event
}
