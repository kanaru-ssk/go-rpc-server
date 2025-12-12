package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/kanaru-ssk/go-http-server/domain/task"
	"github.com/kanaru-ssk/go-http-server/interface/http/response"
	querydecoder "github.com/kanaru-ssk/go-http-server/lib/query"
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

// GET /core/v1/task/get
func (h *TaskHandler) HandleGetV1(w http.ResponseWriter, r *http.Request) {
	var query struct {
		ID string `query:"id"`
	}
	var successResponse response.Task
	var errorResponse response.Error

	ctx := r.Context()

	if err := querydecoder.Decode(r.URL.Query(), &query); err != nil {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleGetV1", "err", err)
		errorResponse = response.MapError(response.ErrInvalidRequestBody)
		response.RenderJson(ctx, w, http.StatusBadRequest, errorResponse)
		return
	}

	t, err := h.taskUseCase.Get(ctx, query.ID)

	// 200
	if err == nil {
		successResponse = response.MapTask(t)
		response.RenderJson(ctx, w, http.StatusOK, successResponse)
		return
	}

	// 400
	if errors.Is(err, task.ErrInvalidID) {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleGetV1", "err", err)
		errorResponse = response.MapError(response.ErrInvalidRequestBody)
		response.RenderJson(ctx, w, http.StatusBadRequest, errorResponse)
		return
	}

	// 404
	if errors.Is(err, task.ErrNotFound) {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleGetV1", "err", err)
		errorResponse = response.MapError(response.ErrNotFound)
		response.RenderJson(ctx, w, http.StatusNotFound, errorResponse)
		return
	}

	// 500
	slog.ErrorContext(ctx, "handler.TaskHandler.HandleGetV1", "err", err)
	errorResponse = response.MapError(response.ErrInternalServerError)
	response.RenderJson(ctx, w, http.StatusInternalServerError, errorResponse)
}

// GET /core/v1/task/list
func (h *TaskHandler) HandleListV1(w http.ResponseWriter, r *http.Request) {
	var successResponse []response.Task
	var errorResponse response.Error

	ctx := r.Context()

	t, err := h.taskUseCase.List(ctx)

	// 200
	if err == nil {
		successResponse = response.MapTaskList(t)
		response.RenderJson(ctx, w, http.StatusOK, successResponse)
		return
	}

	// 500
	slog.ErrorContext(ctx, "handler.TaskHandler.HandleListV1", "err", err)
	errorResponse = response.MapError(response.ErrInternalServerError)
	response.RenderJson(ctx, w, http.StatusInternalServerError, errorResponse)
}

// POST /core/v1/task/create
func (h *TaskHandler) HandleCreateV1(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title string `json:"title"`
	}
	var successResponse response.Task
	var errorResponse response.Error

	ctx := r.Context()

	// 400
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleCreateV1", "err", err)
		errorResponse = response.MapError(response.ErrInvalidRequestBody)
		response.RenderJson(ctx, w, http.StatusBadRequest, errorResponse)
		return
	}

	t, err := h.taskUseCase.Create(ctx, body.Title)

	// 200
	if err == nil {
		successResponse = response.MapTask(t)
		response.RenderJson(ctx, w, http.StatusOK, successResponse)
		return
	}

	// 400
	if errors.Is(err, task.ErrInvalidTitle) {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleCreateV1", "err", err)
		errorResponse = response.MapError(response.ErrInvalidRequestBody)
		response.RenderJson(ctx, w, http.StatusBadRequest, errorResponse)
		return
	}

	// 500
	slog.ErrorContext(ctx, "handler.TaskHandler.HandleCreateV1", "err", err)
	errorResponse = response.MapError(response.ErrInternalServerError)
	response.RenderJson(ctx, w, http.StatusInternalServerError, errorResponse)
}

// PUT /core/v1/task/update
func (h *TaskHandler) HandleUpdateV1(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID     string `json:"id"`
		Title  string `json:"title"`
		Status string `json:"status"`
	}
	var successResponse response.Task
	var errorResponse response.Error

	ctx := r.Context()

	// 400
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleUpdateV1", "err", err)
		errorResponse = response.MapError(response.ErrInvalidRequestBody)
		response.RenderJson(ctx, w, http.StatusBadRequest, errorResponse)
		return
	}

	t, err := h.taskUseCase.Update(ctx, body.ID, body.Title, body.Status)

	// 200
	if err == nil {
		successResponse = response.MapTask(t)
		response.RenderJson(ctx, w, http.StatusOK, successResponse)
		return
	}

	// 400
	if errors.Is(err, task.ErrInvalidID) ||
		errors.Is(err, task.ErrInvalidTitle) ||
		errors.Is(err, task.ErrInvalidStatus) {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleUpdateV1", "err", err)
		errorResponse = response.MapError(response.ErrInvalidRequestBody)
		response.RenderJson(ctx, w, http.StatusBadRequest, errorResponse)
		return
	}

	// 404
	if errors.Is(err, task.ErrNotFound) {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleUpdateV1", "err", err)
		errorResponse = response.MapError(response.ErrNotFound)
		response.RenderJson(ctx, w, http.StatusNotFound, errorResponse)
		return
	}

	// 500
	slog.ErrorContext(ctx, "handler.TaskHandler.HandleUpdateV1", "err", err)
	errorResponse = response.MapError(response.ErrInternalServerError)
	response.RenderJson(ctx, w, http.StatusInternalServerError, errorResponse)
}

// DELETE /core/v1/task/delete
func (h *TaskHandler) HandleDeleteV1(w http.ResponseWriter, r *http.Request) {
	var query struct {
		ID string `query:"id"`
	}
	var errorResponse response.Error

	ctx := r.Context()

	if err := querydecoder.Decode(r.URL.Query(), &query); err != nil {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleDeleteV1", "err", err)
		errorResponse = response.MapError(response.ErrInvalidRequestBody)
		response.RenderJson(ctx, w, http.StatusBadRequest, errorResponse)
		return
	}

	err := h.taskUseCase.Delete(ctx, query.ID)

	// 204
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// 400
	if errors.Is(err, task.ErrInvalidID) {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleDeleteV1", "err", err)
		errorResponse = response.MapError(response.ErrInvalidRequestBody)
		response.RenderJson(ctx, w, http.StatusBadRequest, errorResponse)
		return
	}

	// 404
	if errors.Is(err, task.ErrNotFound) {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleDeleteV1", "err", err)
		errorResponse = response.MapError(response.ErrNotFound)
		response.RenderJson(ctx, w, http.StatusNotFound, errorResponse)
		return
	}

	// 500
	slog.ErrorContext(ctx, "handler.TaskHandler.HandleDeleteV1", "err", err)
	errorResponse = response.MapError(response.ErrInternalServerError)
	response.RenderJson(ctx, w, http.StatusInternalServerError, errorResponse)
}

// PUT /core/v1/task/done
func (h *TaskHandler) HandleDoneV1(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ID string `json:"id"`
	}
	var successResponse response.Task
	var errorResponse response.Error

	ctx := r.Context()

	// 400
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleDoneV1", "err", err)
		errorResponse = response.MapError(response.ErrInvalidRequestBody)
		response.RenderJson(ctx, w, http.StatusBadRequest, errorResponse)
		return
	}

	// 他のチームの作業に影響しないように、開発中はmockデータを返しておく
	// 204
	successResponse = response.Task{
		ID:        "id",
		Title:     "title",
		Status:    "TODO",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	response.RenderJson(ctx, w, http.StatusOK, successResponse)

	// 400

	// 404

	// 500
}
