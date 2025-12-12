package task

import (
	"context"

	"github.com/kanaru-ssk/go-http-server/lib/tx"
)

type Repository interface {
	Get(ctx context.Context, tx tx.Tx, id string) (*Task, error)
	List(ctx context.Context, tx tx.Tx) ([]*Task, error)

	Create(ctx context.Context, tx tx.Tx, task *Task) error
	Update(ctx context.Context, tx tx.Tx, task *Task) error
	Delete(ctx context.Context, tx tx.Tx, id string) error
}
