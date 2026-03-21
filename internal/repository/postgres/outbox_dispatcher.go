package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/jmoiron/sqlx"

	"github.com/ulbwa/medincident-command-service/internal/common/outbox"
	"github.com/ulbwa/medincident-command-service/internal/common/persistence"
)

// OutboxDispatcher writes domain events to the outbox_events table within the
// active transaction retrieved from the context. Events from all sources are
// merged and sorted by their monotonic Sequence before insertion, preserving
// the original cross-aggregate recording order.
type OutboxDispatcher struct{}

// NewOutboxDispatcher returns an outbox.EventDispatcher backed by PostgreSQL.
func NewOutboxDispatcher() *OutboxDispatcher {
	return &OutboxDispatcher{}
}

// sqlTxProvider is a narrow interface satisfied by the postgres transaction
// wrapper, used to extract the underlying *sqlx.Tx without importing the
// concrete type (which lives in the same package, but keeps the contract explicit).
type sqlTxProvider interface {
	SQLTx() *sqlx.Tx
}

// Dispatch collects events from all sources, sorts them by sequence, and inserts
// them into outbox_events within the provided transaction.
func (d *OutboxDispatcher) Dispatch(ctx context.Context, tx persistence.Transaction, sources ...outbox.EventSource) error {
	provider, ok := tx.(sqlTxProvider)
	if !ok {
		return fmt.Errorf("outbox dispatcher: transaction does not implement SQLTx()")
	}

	sqlTx := provider.SQLTx()

	var events []outbox.Event
	for _, src := range sources {
		events = append(events, src.PopEvents()...)
	}

	if len(events) == 0 {
		return nil
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Sequence < events[j].Sequence
	})

	for _, ev := range events {
		if err := d.insert(ctx, sqlTx, ev); err != nil {
			return err
		}
	}

	return nil
}

const insertOutboxEventQuery = `
INSERT INTO outbox_events (sequence, event_type, payload, created_at)
VALUES ($1, $2, $3, now())`

func (d *OutboxDispatcher) insert(ctx context.Context, sqlTx *sqlx.Tx, ev outbox.Event) error {
	payload, err := json.Marshal(ev.Payload)
	if err != nil {
		return fmt.Errorf("marshal outbox event payload: %w", err)
	}

	// %T produces the short qualified type name (e.g. "model.UserCreatedEvent").
	// This is sufficient as a type discriminator for the relay layer.
	eventType := fmt.Sprintf("%T", ev.Payload)

	if _, err := sqlTx.ExecContext(ctx, insertOutboxEventQuery, ev.Sequence, eventType, payload); err != nil {
		return fmt.Errorf("insert outbox event %s: %w", eventType, err)
	}

	return nil
}
