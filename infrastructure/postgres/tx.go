package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanaru-ssk/go-http-server/lib/tx"
)

var _ tx.Manager = (*Manager)(nil)

type Manager struct {
	pool *pgxpool.Pool
}

func NewManager(pool *pgxpool.Pool) *Manager {
	return &Manager{pool: pool}
}

func (m *Manager) WithinTx(ctx context.Context, fn func(ctx context.Context, tx tx.Tx) error) error {
	conn, err := m.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("postgres.Manager.WithinTx: Acquire: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("postgres.Manager.WithinTx: Begin: %w", err)
	}

	if err := fn(ctx, tx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("postgres.Manager.WithinTx: %v (rollback error: %w)", err, rbErr)
		}
		return fmt.Errorf("postgres.Manager.WithinTx: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("postgres.Manager.WithinTx: Commit: %w", err)
	}

	return nil
}

type Tx interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func UnwrapTx(tx tx.Tx, pool *pgxpool.Pool) (Tx, error) {
	// txがnilの時はpoolを使用
	if tx == nil {
		return pool, nil
	}
	pgxTx, ok := tx.(pgx.Tx)
	if !ok {
		return nil, fmt.Errorf("postgres.UnwrapTx: unexpected type %T", tx)
	}
	return pgxTx, nil
}
