package user

import (
	"context"

	"github.com/ulbwa/medincident-command-service/internal/model"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	ExistsByID(ctx context.Context, id int64) (bool, error)
	Save(ctx context.Context, user *model.User) error
}
