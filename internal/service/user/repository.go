package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/ulbwa/medincident-command-service/internal/model"
)

// Repository is the write-side port for user persistence.
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
	ExistsByIdentityID(ctx context.Context, identityID string) (bool, error)
	Save(ctx context.Context, user *model.User) error
}
