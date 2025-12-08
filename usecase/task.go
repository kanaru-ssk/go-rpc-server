package usecase

import (
	"context"
	"fmt"

	"github.com/kanaru-ssk/go-http-server/entity/task"
	"github.com/kanaru-ssk/go-http-server/lib/tx"
)

type TaskUseCase struct {
	txManager      tx.Manager
	taskFactory    *task.Factory
	taskRepository task.Repository
}

func NewTaskUseCase(txManager tx.Manager, taskFactory *task.Factory, taskRepository task.Repository) *TaskUseCase {
	return &TaskUseCase{
		txManager:      txManager,
		taskFactory:    taskFactory,
		taskRepository: taskRepository,
	}
}

func (u *TaskUseCase) Get(ctx context.Context, id string) (*task.Task, error) {
	pi, err := task.ParseID(id)
	if err != nil {
		return nil, fmt.Errorf("usecase.TaskUseCase.Get: %w", err)
	}
	task, err := u.taskRepository.Get(ctx, pi)
	if err != nil {
		return nil, fmt.Errorf("usecase.TaskUseCase.Get: %w", err)
	}
	return task, nil
}

func (u *TaskUseCase) List(ctx context.Context) ([]*task.Task, error) {
	tasks, err := u.taskRepository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("usecase.TaskUseCase.List: %w", err)
	}
	return tasks, nil
}

func (u *TaskUseCase) Create(ctx context.Context, title string) (*task.Task, error) {
	pt, err := task.ParseTitle(title)
	if err != nil {
		return nil, fmt.Errorf("usecase.TaskUseCase.Create: %w", err)
	}
	task := u.taskFactory.New(pt)
	if err := u.taskRepository.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("usecase.TaskUseCase.Create: %w", err)
	}
	return task, nil
}

func (u *TaskUseCase) Update(ctx context.Context, id, title string, status string) (*task.Task, error) {
	pi, err := task.ParseID(id)
	if err != nil {
		return nil, fmt.Errorf("usecase.TaskUseCase.Update: %w", err)
	}
	pt, err := task.ParseTitle(title)
	if err != nil {
		return nil, fmt.Errorf("usecase.TaskUseCase.Update: %w", err)
	}
	ps, err := task.ParseStatus(status)
	if err != nil {
		return nil, fmt.Errorf("usecase.TaskUseCase.Update: %w", err)
	}
	var task *task.Task
	if err := u.txManager.WithinTx(ctx, func(ctx context.Context) error {
		task, err = u.taskRepository.Get(ctx, pi)
		if err != nil {
			return err
		}
		task.Update(pt, ps)
		if err := u.taskRepository.Update(ctx, task); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("usecase.TaskUseCase.Update: %w", err)
	}

	return task, nil
}

func (u *TaskUseCase) Delete(ctx context.Context, id string) error {
	pi, err := task.ParseID(id)
	if err != nil {
		return fmt.Errorf("usecase.TaskUseCase.Delete: %w", err)
	}
	if err := u.taskRepository.Delete(ctx, pi); err != nil {
		return fmt.Errorf("usecase.TaskUseCase.Delete: %w", err)
	}
	return nil
}
