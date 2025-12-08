package response

import (
	"time"

	"github.com/kanaru-ssk/go-http-server/entity/task"
)

type Task struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func MapTask(t *task.Task) Task {
	return Task{
		ID:        t.ID,
		Title:     string(t.Title),
		Status:    string(t.Status),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func MapTaskList(t []*task.Task) []Task {
	result := make([]Task, 0, len(t))
	for _, v := range t {
		result = append(result, MapTask(v))
	}
	return result
}
