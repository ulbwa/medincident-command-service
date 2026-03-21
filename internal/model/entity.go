package model

import (
	"sync/atomic"

	"github.com/ulbwa/medincident-command-service/internal/common/outbox"
)

// globalSeq is a process-wide monotonic counter used to preserve the original
// recording order of domain events across multiple aggregates. Each call to
// recordEvent claims the next value atomically, so events from different
// entities can be sorted by outbox.Event.Sequence even after their slices are
// merged by the EventDispatcher.
var globalSeq atomic.Uint64

// Entity is the base type for aggregates and domain entities that produce domain events.
// Embed it by value into aggregate structs to gain EventSource behaviour without
// boilerplate. It satisfies the outbox.EventSource interface implicitly.
//
// Usage:
//
//	type Order struct {
//	    model.Entity
//	    id     uuid.UUID
//	    status OrderStatus
//	}
//
//	func (o *Order) Cancel() error {
//	    if o.status == OrderStatusCancelled {
//	        return ErrAlreadyCancelled
//	    }
//	    o.status = OrderStatusCancelled
//	    o.recordEvent(OrderCancelledEvent{ID: o.id})
//	    return nil
//	}
type Entity struct {
	events []outbox.Event
}

// recordEvent appends a domain event to the entity's pending event list and
// stamps it with a globally unique sequence number. The sequence is assigned
// atomically, so events recorded across multiple aggregates within the same
// use-case can later be sorted into their original recording order.
//
// Call this from within domain methods after a successful state transition.
func (e *Entity) recordEvent(payload any) {
	e.events = append(e.events, outbox.Event{
		Sequence: globalSeq.Add(1),
		Payload:  payload,
	})
}

// PopEvents returns all pending domain events in recording order and clears
// the internal list. Satisfies the outbox.EventSource interface.
func (e *Entity) PopEvents() []outbox.Event {
	if len(e.events) == 0 {
		return nil
	}
	events := e.events
	e.events = nil
	return events
}
