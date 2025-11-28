package taskresponse

import (
	"time"

	"github.com/kanaru-ssk/go-rpc-server/domain/task"
)

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

type TaskResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (m *Mapper) MapGetResponse(t *task.Task) *TaskResponse {
	return &TaskResponse{
		ID:        t.ID,
		Title:     string(t.Title),
		Status:    string(t.Status),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func (m *Mapper) MapListResponse(t []*task.Task) []*TaskResponse {
	result := make([]*TaskResponse, 0, len(t))
	for _, v := range t {
		result = append(result, m.MapGetResponse(v))
	}
	return result
}

func (m *Mapper) MapCreateResponse(t *task.Task) *TaskResponse {
	return &TaskResponse{
		ID:        t.ID,
		Title:     string(t.Title),
		Status:    string(t.Status),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func (m *Mapper) MapUpdateResponse(t *task.Task) *TaskResponse {
	return &TaskResponse{
		ID:        t.ID,
		Title:     string(t.Title),
		Status:    string(t.Status),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
