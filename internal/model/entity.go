package model

// Entity provides basic functionality for domain entities,
// specifically around tracking domain events for CQRS and Event Sourcing.
type Entity struct {
	events []any
}

func (e *Entity) recordEvent(event any) {
	e.events = append(e.events, event)
}

// PopEvents returns all recorded domain events and clears the internal list.
func (e *Entity) PopEvents() []any {
	events := e.events
	e.events = nil
	return events
}
