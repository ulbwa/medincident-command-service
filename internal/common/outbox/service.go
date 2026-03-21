package outbox

import (
	"context"

	"github.com/ulbwa/medincident-command-service/internal/common/persistence"
)

// Event is a domain event stamped with a monotonically increasing sequence number.
// The sequence is assigned at recording time using a process-wide atomic counter,
// so events from different aggregates can be merged and sorted into their original
// recording order before being dispatched.
type Event struct {
	// Sequence is a process-wide monotonic counter value assigned when the event
	// was recorded. Use it to sort events from multiple EventSource instances.
	Sequence uint64
	// Payload is the domain event value (e.g. OrderCancelledEvent{}).
	Payload any
}

// EventSource is implemented by domain entities that produce domain events.
// PopEvents returns all pending events in recording order and clears the list.
type EventSource interface {
	PopEvents() []Event
}

// EventDispatcher writes domain events produced by the given sources to the
// outbox table within the provided transaction. Implementations must merge
// events from all sources sorted by Event.Sequence to preserve the original
// cross-aggregate recording order.
type EventDispatcher interface {
	Dispatch(ctx context.Context, tx persistence.Transaction, sources ...EventSource) error
}
