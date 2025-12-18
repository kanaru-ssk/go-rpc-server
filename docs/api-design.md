# API 設計

## 概要

RPC スタイルの HTTP API

## 命名規則

パスは以下の命名規則に従う。

```
POST /<serviceName>/<version>/<useCaseName>/<methodName>

例: POST /core/v1/task/done
```

- 単語は `lowerCamelCase` で統一
- `<serviceName>` はマイクロサービス化を想定して設定、モノリス時は`core`というサービス名を使う。
- `<useCaseName>`, `<methodName>` は`usecase`ディレクトリ配下の構造体名、メソッド名に合わせる。
- `<methodName>` は常にパスに指定し、処理の種別（読み取り/更新など）はメソッド名のプレフィックスで判断できるようにする。取得処理は`get*`, `list*` のプレフィックスを付ける。

### 通常の CRUD のメソッド名

```
POST /core/v1/task/create
POST /core/v1/task/get
POST /core/v1/task/update
POST /core/v1/task/delete
```

## API スキーマ定義

OpenAPI などを使うとスキーマファイルの管理コストがかかるため用意しない。

代わりに、ソースコードから API のパス、リクエスト、レスポンスを読み取りやすくする。

`cmd/httpserver/main.go`を見れば全てのエンドポイントが分かる。

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
// interface/http/handler/task.go

// POST /core/v1/task/get
func (h *TaskHandler) HandleGetV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request struct {
		ID string `json:"id"`
	}

	ctx := r.Context()

	// 型のパースだけinterfaceレイヤーで行う。
	// 詳細なバリデーションはusecase以降で行う。
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
		{Err: task.ErrInvalidID, StatusCode: http.StatusBadRequest, ErrorCode: response.ErrInvalidRequestBody},
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
```
