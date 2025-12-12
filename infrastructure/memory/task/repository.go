package task

import (
	"context"
	"fmt"

	"github.com/kanaru-ssk/go-http-server/domain/task"
	"github.com/kanaru-ssk/go-http-server/lib/tx"
)

func NewRepository(
	tasks map[string]*task.Task,
) task.Repository {
	return &repository{
		tasks: tasks,
	}
}

type repository struct {
	tasks map[string]*task.Task
}

func (r *repository) Get(ctx context.Context, _ tx.Tx, id string) (*task.Task, error) {
	t, ok := r.tasks[id]
	if !ok || t == nil {
		return nil, fmt.Errorf("task.repository.Get: %w", task.ErrNotFound)
	}

	return cloneTask(t), nil
}

func (r *repository) List(ctx context.Context, _ tx.Tx) ([]*task.Task, error) {
	list := make([]*task.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		if t == nil {
			continue
		}
		list = append(list, cloneTask(t))
	}

	return list, nil
}

func (r *repository) Create(ctx context.Context, _ tx.Tx, t *task.Task) error {
	if r.tasks == nil {
		r.tasks = make(map[string]*task.Task)
	}
	r.tasks[t.ID] = t

	return nil
}

func (r *repository) Update(ctx context.Context, _ tx.Tx, t *task.Task) error {
	cur, ok := r.tasks[t.ID]
	if !ok || cur == nil {
		return fmt.Errorf("task.repository.Update: %w", task.ErrNotFound)
	}
	r.tasks[t.ID] = t

	return nil
}

func (r *repository) Delete(ctx context.Context, _ tx.Tx, id string) error {
	if _, ok := r.tasks[id]; !ok {
		return task.ErrNotFound
	}

	delete(r.tasks, id)

	return nil
}

func cloneTask(t *task.Task) *task.Task {
	if t == nil {
		return nil
	}
	copy := *t
	return &copy
}
