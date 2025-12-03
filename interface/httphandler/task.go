package httphandler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/kanaru-ssk/go-rpc-server/domain/task"
	"github.com/kanaru-ssk/go-rpc-server/interface/httpresponse"
	"github.com/kanaru-ssk/go-rpc-server/interface/response/errorresponse"
	"github.com/kanaru-ssk/go-rpc-server/interface/response/taskresponse"
	"github.com/kanaru-ssk/go-rpc-server/usecase"
)

type TaskHandler struct {
	taskUsecase *usecase.TaskUsecase
	taskMapper  *taskresponse.Mapper
	errorMapper *errorresponse.Mapper
}

func NewTaskHandler(
	taskUsecase *usecase.TaskUsecase,
	taskMapper *taskresponse.Mapper,
	errorMapper *errorresponse.Mapper,
) *TaskHandler {
	return &TaskHandler{
		taskUsecase: taskUsecase,
		taskMapper:  taskMapper,
		errorMapper: errorMapper,
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
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleGetV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))
		return
	}

	t, err := h.taskUsecase.Get(ctx, request.ID)
	switch {

	// 200
	case err == nil:
		httpresponse.RenderJson(ctx, w, http.StatusOK, h.taskMapper.MapGetResponse(t))

	// 400
	case errors.Is(err, task.ErrInvalidID):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleGetV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))

	// 404
	case errors.Is(err, task.ErrNotFound):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleGetV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusNotFound, h.errorMapper.MapErrorResponse(errorresponse.ErrNotFound))

	// 500
	default:
		slog.ErrorContext(ctx, "httphandler.TaskHandler.HandleGetV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusInternalServerError, h.errorMapper.MapErrorResponse(errorresponse.ErrInternalServerError))
	}
}

// POST /core/v1/task/list
func (h *TaskHandler) HandleListV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	t, err := h.taskUsecase.List(ctx)
	switch {

	// 200
	case err == nil:
		httpresponse.RenderJson(ctx, w, http.StatusOK, h.taskMapper.MapListResponse(t))

	// 500
	default:
		slog.ErrorContext(ctx, "httphandler.TaskHandler.HandleListV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusInternalServerError, h.errorMapper.MapErrorResponse(errorresponse.ErrInternalServerError))
	}
}

// POST /core/v1/task/create
func (h *TaskHandler) HandleCreateV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request struct {
		Title string `json:"title"`
	}

	// 400
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleCreateV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))
		return
	}

	t, err := h.taskUsecase.Create(ctx, request.Title)
	switch {

	// 200
	case err == nil:
		httpresponse.RenderJson(ctx, w, http.StatusOK, h.taskMapper.MapCreateResponse(t))

	// 400
	case errors.Is(err, task.ErrInvalidTitle):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleCreateV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))

	// 500
	default:
		slog.ErrorContext(ctx, "httphandler.TaskHandler.HandleCreateV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusInternalServerError, h.errorMapper.MapErrorResponse(errorresponse.ErrInternalServerError))
	}
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
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleUpdateV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))
		return
	}

	t, err := h.taskUsecase.Update(ctx, request.ID, request.Title, request.Status)
	switch {

	// 200
	case err == nil:
		httpresponse.RenderJson(ctx, w, http.StatusOK, h.taskMapper.MapUpdateResponse(t))

	// 400
	case errors.Is(err, task.ErrInvalidID),
		errors.Is(err, task.ErrInvalidTitle),
		errors.Is(err, task.ErrInvalidStatus):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleUpdateV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))

	// 404
	case errors.Is(err, task.ErrNotFound):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleUpdateV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusNotFound, h.errorMapper.MapErrorResponse(errorresponse.ErrNotFound))

	// 500
	default:
		slog.ErrorContext(ctx, "httphandler.TaskHandler.HandleUpdateV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusInternalServerError, h.errorMapper.MapErrorResponse(errorresponse.ErrInternalServerError))
	}
}

// POST /core/v1/task/delete
func (h *TaskHandler) HandleDeleteV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request struct {
		ID string `json:"id"`
	}

	// 400
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleDeleteV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))
		return
	}

	err := h.taskUsecase.Delete(ctx, request.ID)
	switch {

	// 204
	case err == nil:
		w.WriteHeader(http.StatusNoContent)

	// 400
	case errors.Is(err, task.ErrInvalidID):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleDeleteV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))

	// 404
	case errors.Is(err, task.ErrNotFound):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleDeleteV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusNotFound, h.errorMapper.MapErrorResponse(errorresponse.ErrNotFound))

	// 500
	default:
		slog.ErrorContext(ctx, "httphandler.TaskHandler.HandleDeleteV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusInternalServerError, h.errorMapper.MapErrorResponse(errorresponse.ErrInternalServerError))
	}
}

// POST /core/v1/task/done
func (h *TaskHandler) HandleDoneV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request struct {
		ID string `json:"id"`
	}

	// 400
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleDoneV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))
		return
	}

	err := h.taskUsecase.Done(ctx, request.ID)
	switch {

	// 204
	case err == nil:
		w.WriteHeader(http.StatusNoContent)

	// 400
	case errors.Is(err, task.ErrInvalidID):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleDoneV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))

	// 404
	case errors.Is(err, task.ErrNotFound):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleDoneV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusNotFound, h.errorMapper.MapErrorResponse(errorresponse.ErrNotFound))

	// 500
	default:
		slog.ErrorContext(ctx, "httphandler.TaskHandler.HandleDoneV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusInternalServerError, h.errorMapper.MapErrorResponse(errorresponse.ErrInternalServerError))
	}
}
