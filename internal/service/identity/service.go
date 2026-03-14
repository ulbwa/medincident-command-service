package identity

import "context"

type Service interface {
	Get(ctx context.Context, id int64) (*Identity, error)
	UpdateHuman(ctx context.Context, id int64, human *Human) (*Identity, error)
}
