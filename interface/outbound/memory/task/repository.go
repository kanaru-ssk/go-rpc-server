package task

import (
	"context"
	"fmt"
	"sync"

	"github.com/kanaru-ssk/go-http-server/entity/task"
)

func NewRepository(
	mu *sync.RWMutex,
	tasks map[string]*task.Task,
) task.Repository {
	return &repository{
		mu:    mu,
		tasks: tasks,
	}
}

type repository struct {
	mu    *sync.RWMutex
	tasks map[string]*task.Task
}

func (r *repository) Get(ctx context.Context, id string) (*task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.tasks[id]
	if !ok || t == nil {
		return nil, fmt.Errorf("task.repository.Get: %w", task.ErrNotFound)
	}

	return t, nil
}

func (r *repository) List(ctx context.Context) ([]*task.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]*task.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		if t == nil {
			continue
		}
		list = append(list, t)
	}

	return list, nil
}

func (r *repository) Create(ctx context.Context, t *task.Task) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.tasks == nil {
		r.tasks = make(map[string]*task.Task)
	}
	r.tasks[t.ID] = t

	return nil
}

func (r *repository) Update(ctx context.Context, t *task.Task) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cur, ok := r.tasks[t.ID]
	if !ok || cur == nil {
		return fmt.Errorf("task.repository.Update: %w", task.ErrNotFound)
	}
	r.tasks[t.ID] = t

	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, ok := r.tasks[id]; !ok {
		return task.ErrNotFound
	}

	delete(r.tasks, id)

	return nil
}
