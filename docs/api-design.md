# API 設計

## 命名規則

全て POST で統一し、パスは以下の命名規則に従う。

```
POST /serviceName/v1/useCaseName/methodName
```

- 単語は lowerCamel
- serviceName はマイクロサービス化を想定して設定、モノリス時は`core`というサービス名を使う。
- useCaseName, methodName は`usecase`ディレクトリ配下の構造体名、メソッド名に合わせる。
- get, list から始まる methodName は安全性・冪等性を担保する。

## API スキーマ定義

OpenAPI などを使うとスキーマファイルの管理コストがかかるため用意しない。

代わりに、ソースコードを綺麗に保ち API のパス、リクエスト、レスポンスを読みやすくする。

`cmd/httpserver/main.go`を見れば API の一覧が見れる。

```go
// cmd/httpserver/main.go
mux.HandleFunc("GET /healthz", handler.HandleGetHealthz)

mux.HandleFunc("POST /core/v1/task/get", taskHandler.HandleGetV1)
mux.HandleFunc("POST /core/v1/task/list", taskHandler.HandleListV1)
mux.HandleFunc("POST /core/v1/task/create", taskHandler.HandleCreateV1)
mux.HandleFunc("POST /core/v1/task/update", taskHandler.HandleUpdateV1)
mux.HandleFunc("POST /core/v1/task/delete", taskHandler.HandleDeleteV1)
mux.HandleFunc("POST /core/v1/task/done", taskHandler.HandleDoneV1)
```

handler のコードを追えばリクエスト、レスポンスの内容が分かる。

```go
// interface/inbound/http/handler/task.go

// POST /core/v1/task/get
func (h *TaskHandler) HandleGetV1(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ID string `json:"id"`
	}
	var successResponse response.Task
	var errorResponse response.Error

	ctx := r.Context()

	// 400
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		slog.WarnContext(ctx, "handler.TaskHandler.HandleGetV1", "err", err)
		errorResponse = response.MapError(response.ErrInvalidRequestBody)
		response.RenderJson(ctx, w, http.StatusBadRequest, errorResponse)
		return
	}

	t, err := h.taskUseCase.Get(ctx, request.ID)

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
```
