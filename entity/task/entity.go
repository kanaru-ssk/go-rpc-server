package task

import (
	"fmt"
	"strings"
	"time"
)

type Status string

const (
	StatusTodo Status = "TODO"
	StatusDone Status = "DONE"
)

type Task struct {
	ID        string
	Title     string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t *Task) UpdateTitle(title string) error {
	pt, err := ParseTitle(title)
	if err != nil {
		return fmt.Errorf("task.UpdateTitle: %w", err)
	}
	t.Title = pt
	return nil
}

func (t *Task) UpdateStatus(status Status) error {
	ps, err := ParseStatus(string(status))
	if err != nil {
		return fmt.Errorf("task.UpdateStatus: %w", err)
	}
	t.Status = ps
	return nil
}

func ParseID(id string) (string, error) {
	s := strings.TrimSpace(id)
	if s == "" {
		return "", fmt.Errorf("task.ParseID: %w", ErrInvalidID)
	}
	return s, nil
}

func ParseTitle(title string) (string, error) {
	s := strings.TrimSpace(title)
	if s == "" {
		return "", fmt.Errorf("task.ParseTitle: %w", ErrInvalidTitle)
	}
	return s, nil
}

func ParseStatus(status string) (Status, error) {
	switch s := Status(status); s {
	case StatusTodo, StatusDone:
		return s, nil
	default:
		return "", fmt.Errorf("task.ParseStatus: %w", ErrInvalidStatus)
	}
}
