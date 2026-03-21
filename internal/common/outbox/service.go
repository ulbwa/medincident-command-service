package outbox

import (
	"context"

	"github.com/ulbwa/medincident-command-service/internal/common/persistence"
)

// Event is a domain event stamped with a process-local sequence number.
// The sequence is assigned at recording time using a process-wide atomic counter
// and is used exclusively to sort events from different aggregates into their
// original recording order before insertion. It is not persisted to the database;
// the outbox_events.id (bigserial) is the canonical global ordering field for
// the relay, as PostgreSQL sequences are non-transactional and monotonic across
// all service instances.
type Event struct {
	// Sequence is a process-local monotonic counter value assigned when the event
	// was recorded. Used only to sort events from multiple EventSource instances
	// before insertion — not stored in the database.
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
