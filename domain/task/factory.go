package task

import (
	"time"

	"github.com/kanaru-ssk/go-http-server/lib/id"
)

type Factory struct {
	idGenerator id.Generator
}

func NewFactory(
	idGenerator id.Generator,
) *Factory {
	return &Factory{
		idGenerator: idGenerator,
	}
}

func (f *Factory) New(title string) *Task {
	id := f.idGenerator.NewID()
	return &Task{
		ID:        id,
		Title:     title,
		Status:    StatusTodo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
