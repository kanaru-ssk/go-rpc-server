package task

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanaru-ssk/go-http-server/domain/task"
	"github.com/kanaru-ssk/go-http-server/infrastructure/postgres"
	"github.com/kanaru-ssk/go-http-server/lib/tx"
)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(
	pool *pgxpool.Pool,
) task.Repository {
	return &repository{
		pool: pool,
	}
}

//go:embed sql/select_task.sql
var selectTaskSql string

func (r *repository) Get(ctx context.Context, ltx tx.Tx, id string) (*task.Task, error) {
	tx, err := postgres.UnwrapTx(ltx, r.pool)
	if err != nil {
		return nil, fmt.Errorf("task.repository.Get: postgres.UnwrapTx: %w", err)
	}

	var t task.Task
	if err := tx.
		QueryRow(ctx, selectTaskSql, id).
		Scan(
			&t.ID,
			&t.Title,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task.repository.Get: %w: %w", task.ErrNotFound, err)
		}
		return nil, fmt.Errorf("ask.repository.Get: %w", err)
	}

	return &t, nil
}

//go:embed sql/select_all_task.sql
var selectAllTaskSql string

func (r *repository) List(ctx context.Context, ltx tx.Tx) ([]*task.Task, error) {
	tx, err := postgres.UnwrapTx(ltx, r.pool)
	if err != nil {
		return nil, fmt.Errorf("task.repository.List: postgres.UnwrapTx: %w", err)
	}

	rows, err := tx.Query(ctx, selectAllTaskSql)
	if err != nil {
		return nil, fmt.Errorf("task.repository.List: tx.Query: %w", err)
	}
	defer rows.Close()

	var list []*task.Task
	for rows.Next() {
		var t task.Task
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("task.repository.List: rows.Scan: %w", err)
		}
		list = append(list, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("task.repository.List: rows.Err: %w", err)
	}
	return list, nil
}

//go:embed sql/insert_task.sql
var insertTaskSQL string

func (r *repository) Create(ctx context.Context, ltx tx.Tx, t *task.Task) error {
	tx, err := postgres.UnwrapTx(ltx, r.pool)
	if err != nil {
		return fmt.Errorf("task.repository.Create: postgres.UnwrapTx: %w", err)
	}

	_, err = tx.Exec(ctx, insertTaskSQL,
		t.ID,
		t.Title,
		t.Status,
		t.CreatedAt,
		t.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("task.repository.Create: tx.Exec: %w", err)
	}

	return nil
}

//go:embed sql/update_task.sql
var updateTaskSQL string

func (r *repository) Update(ctx context.Context, ltx tx.Tx, t *task.Task) error {
	tx, err := postgres.UnwrapTx(ltx, r.pool)
	if err != nil {
		return fmt.Errorf("task.repository.Update: postgres.UnwrapTx: %w", err)
	}

	_, err = tx.Exec(ctx, updateTaskSQL,
		t.ID,
		t.Title,
		t.Status,
		t.CreatedAt,
		t.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("task.repository.Update: %w", err)
	}

	return nil
}

//go:embed sql/delete_task.sql
var deleteTaskSQL string

func (r *repository) Delete(ctx context.Context, ltx tx.Tx, id string) error {
	tx, err := postgres.UnwrapTx(ltx, r.pool)
	if err != nil {
		return fmt.Errorf("task.repository.Delete: postgres.UnwrapTx: %w", err)
	}

	cmd, err := tx.Exec(ctx, deleteTaskSQL, id)
	if err != nil {
		return fmt.Errorf("task.repository.Delete: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("task.repository.Delete: %w", task.ErrNotFound)
	}

	return nil
}
