package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/kanaru-ssk/go-http-server/domain/task"
	"github.com/kanaru-ssk/go-http-server/interface/http/response"
	"github.com/kanaru-ssk/go-http-server/usecase"
)

type TaskHandler struct {
	taskUseCase *usecase.TaskUseCase
}

func NewTaskHandler(
	taskUseCase *usecase.TaskUseCase,
) *TaskHandler {
	return &TaskHandler{
		taskUseCase: taskUseCase,
	}
}

// POST /core/v1/task/get
func (h *TaskHandler) HandleGetV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request struct {
		ID string `json:"id"`
	}

	// 400
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleGetV1: request json decode error", "err", err)
		response.RenderJson(ctx, w, http.StatusBadRequest, response.ErrorJson{ErrorCode: response.ErrInvalidRequestBody})
		return
	}

	t, err := h.taskUseCase.Get(ctx, request.ID)

	// 200
	if err == nil {
		response.RenderJson(ctx, w, http.StatusOK, response.MapTask(t))
		return
	}

	maps := []response.ErrResMap{
		// 400
		{Err: task.ErrInvalidID, StatusCode: http.StatusBadRequest, ErrorCode: response.ErrTaskInvalidID},
		// 404
		{Err: task.ErrNotFound, StatusCode: http.StatusNotFound, ErrorCode: response.ErrNotFound},
	}

	for _, m := range maps {
		if errors.Is(err, m.Err) {
			slog.ErrorContext(ctx, "handler.TaskHandler.HandleGetV1", "err", err)
			response.RenderJson(ctx, w, m.StatusCode, response.ErrorJson{ErrorCode: m.ErrorCode})
			return
		}
	}

	// 500
	slog.ErrorContext(ctx, "handler.TaskHandler.HandleGetV1", "err", err)
	response.RenderJson(ctx, w, http.StatusInternalServerError, response.ErrorJson{ErrorCode: response.ErrInternalServerError})
}

// POST /core/v1/task/list
func (h *TaskHandler) HandleListV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	t, err := h.taskUseCase.List(ctx)

	// 200
	if err == nil {
		response.RenderJson(ctx, w, http.StatusOK, response.MapTaskList(t))
		return
	}

	// 500
	slog.ErrorContext(ctx, "handler.TaskHandler.HandleListV1", "err", err)
	response.RenderJson(ctx, w, http.StatusInternalServerError, response.ErrorJson{ErrorCode: response.ErrInternalServerError})
}

// POST /core/v1/task/create
func (h *TaskHandler) HandleCreateV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request struct {
		Title string `json:"title"`
	}

	// 400
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleCreateV1: request json decode error", "err", err)
		response.RenderJson(ctx, w, http.StatusBadRequest, response.ErrorJson{ErrorCode: response.ErrInvalidRequestBody})
		return
	}

	t, err := h.taskUseCase.Create(ctx, request.Title)

	// 200
	if err == nil {
		response.RenderJson(ctx, w, http.StatusOK, response.MapTask(t))
		return
	}

	maps := []response.ErrResMap{
		// 400
		{Err: task.ErrInvalidTitle, StatusCode: http.StatusBadRequest, ErrorCode: response.ErrTaskInvalidTitle},
	}

	for _, m := range maps {
		if errors.Is(err, m.Err) {
			slog.ErrorContext(ctx, "handler.TaskHandler.HandleCreateV1", "err", err)
			response.RenderJson(ctx, w, m.StatusCode, response.ErrorJson{ErrorCode: m.ErrorCode})
			return
		}
	}

	// 500
	slog.ErrorContext(ctx, "handler.TaskHandler.HandleCreateV1", "err", err)
	response.RenderJson(ctx, w, http.StatusInternalServerError, response.ErrorJson{ErrorCode: response.ErrInternalServerError})
}

// POST /core/v1/task/update
func (h *TaskHandler) HandleUpdateV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Status string `json:"status"`
	}

	// 400
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleUpdateV1: request json decode error", "err", err)
		response.RenderJson(ctx, w, http.StatusBadRequest, response.ErrorJson{ErrorCode: response.ErrInvalidRequestBody})
		return
	}

	t, err := h.taskUseCase.Update(ctx, request.ID, request.Title, request.Status)

	// 200
	if err == nil {
		response.RenderJson(ctx, w, http.StatusOK, response.MapTask(t))
		return
	}

	maps := []response.ErrResMap{
		// 400
		{Err: task.ErrInvalidID, StatusCode: http.StatusBadRequest, ErrorCode: response.ErrTaskInvalidID},
		{Err: task.ErrInvalidTitle, StatusCode: http.StatusBadRequest, ErrorCode: response.ErrTaskInvalidTitle},
		{Err: task.ErrInvalidStatus, StatusCode: http.StatusBadRequest, ErrorCode: response.ErrTaskInvalidStatus},
		// 404
		{Err: task.ErrNotFound, StatusCode: http.StatusNotFound, ErrorCode: response.ErrNotFound},
	}

	for _, m := range maps {
		if errors.Is(err, m.Err) {
			slog.ErrorContext(ctx, "handler.TaskHandler.HandleUpdateV1", "err", err)
			response.RenderJson(ctx, w, m.StatusCode, response.ErrorJson{ErrorCode: m.ErrorCode})
			return
		}
	}

	// 500
	slog.ErrorContext(ctx, "handler.TaskHandler.HandleUpdateV1", "err", err)
	response.RenderJson(ctx, w, http.StatusInternalServerError, response.ErrorJson{ErrorCode: response.ErrInternalServerError})
}

// POST /core/v1/task/delete
func (h *TaskHandler) HandleDeleteV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request struct {
		ID string `json:"id"`
	}

	// 400
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleDeleteV1: request json decode error", "err", err)
		response.RenderJson(ctx, w, http.StatusBadRequest, response.ErrorJson{ErrorCode: response.ErrInvalidRequestBody})
		return
	}

	err := h.taskUseCase.Delete(ctx, request.ID)

	// 204
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	maps := []response.ErrResMap{
		// 400
		{Err: task.ErrInvalidID, StatusCode: http.StatusBadRequest, ErrorCode: response.ErrTaskInvalidID},
		// 404
		{Err: task.ErrNotFound, StatusCode: http.StatusNotFound, ErrorCode: response.ErrNotFound},
	}

	for _, m := range maps {
		if errors.Is(err, m.Err) {
			slog.ErrorContext(ctx, "handler.TaskHandler.HandleDeleteV1", "err", err)
			response.RenderJson(ctx, w, m.StatusCode, response.ErrorJson{ErrorCode: m.ErrorCode})
			return
		}
	}

	// 500
	slog.ErrorContext(ctx, "handler.TaskHandler.HandleDeleteV1", "err", err)
	response.RenderJson(ctx, w, http.StatusInternalServerError, response.ErrorJson{ErrorCode: response.ErrInternalServerError})
}

// POST /core/v1/task/done
func (h *TaskHandler) HandleDoneV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request struct {
		ID string `json:"id"`
	}

	// 400
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleDoneV1: request json decode error", "err", err)
		response.RenderJson(ctx, w, http.StatusBadRequest, response.ErrorJson{ErrorCode: response.ErrInvalidRequestBody})
		return
	}

	// 他のチームの作業に影響しないように、開発中はmockデータを返しておく
	// 200
	response.RenderJson(ctx, w, http.StatusOK, response.Task{
		ID:        "id",
		Title:     "title",
		Status:    "TODO",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	// 400

	// 404

	// 500
}
