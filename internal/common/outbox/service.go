package outbox

import (
	"context"

	"github.com/ulbwa/medincident-command-service/internal/common/persistence"
)

type EventSource interface {
	PopEvents() []any
}

type EventDispatcher interface {
	Dispatch(ctx context.Context, tx persistence.Transaction, sources ...EventSource) error
}
