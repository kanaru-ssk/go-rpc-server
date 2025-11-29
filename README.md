# go-rpc-server

Go で RPC スタイルの API を実装するサンプル

## API 設計

### 命名規則

全て POST で統一し、パスは以下の命名規則に従う。

```
POST /serviceName/v1/useCaseName/methodName
```

- 単語は lowerCamel
- serviceName はマイクロサービス化を想定して設定、モノリス時は`core`というサービス名を使う。
- useCaseName, methodName は`usecase`ディレクトリ配下の構造体名、メソッド名に合わせる。
- get, list から始まる methodName は安全性・冪等性を担保する。

### API 定義

OpenAPI などのスキーマは書かない代わりに、ソースコードを綺麗に保ち API のパス、リクエスト、レスポンスを読みやすくする。

`main.go`を見れば API の一覧が見れる。

```go
// main.go
mux.HandleFunc("/healthz", httphandler.HandleGetHealthz)

mux.HandleFunc("/core/v1/task/get", taskHandler.HandleGetV1)
mux.HandleFunc("/core/v1/task/list", taskHandler.HandleListV1)
mux.HandleFunc("/core/v1/task/create", taskHandler.HandleCreateV1)
mux.HandleFunc("/core/v1/task/update", taskHandler.HandleUpdateV1)
mux.HandleFunc("/core/v1/task/delete", taskHandler.HandleDeleteV1)
```

handler のコードを追えばリクエスト、レスポンスの内容が分かる。

```go
// interface/httphandler/task.go
func (h *TaskHandler) HandleGetV1(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
	}

	ctx := r.Context()

	// 405
	if r.Method != http.MethodPost {
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleGetV1", "err", errorresponse.ErrMethodNotAllowed)
		httpresponse.RenderJson(ctx, w, http.StatusMethodNotAllowed, h.errorMapper.MapErrorResponse(errorresponse.ErrMethodNotAllowed))
		return
	}

	// 400
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleGetV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))
		return
	}

	t, err := h.taskUseCase.Get(ctx, req.ID)
	switch {

	// 200
	case err == nil:
		httpresponse.RenderJson(ctx, w, http.StatusOK, h.taskMapper.MapGetResponse(t))

	// 400
	case errors.Is(err, task.ErrInvalidID):
		slog.WarnContext(ctx, "httphandler.TaskHandler.HandleGetV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusBadRequest, h.errorMapper.MapErrorResponse(errorresponse.ErrInvalidRequestBody))

	// 500
	default:
		slog.ErrorContext(ctx, "httphandler.TaskHandler.HandleGetV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusInternalServerError, h.errorMapper.MapErrorResponse(errorresponse.ErrInternalServerError))
	}
}
```

## 起動方法

```sh
docker build -t go-rpc-server .
docker run --rm -it -p 8000:8000 go-rpc-server
```

## 動作確認コマンド

```sh
curl -X POST localhost:8000/core/v1/task/get -d '{ "id": "id_01" }'
curl -X POST localhost:8000/core/v1/task/list
curl -X POST localhost:8000/core/v1/task/create -d '{ "title": "title_01" }'
curl -X POST localhost:8000/core/v1/task/update -d '{ "id": "id_01", "title": "title_01", "status": "DONE" }'
curl -X POST localhost:8000/core/v1/task/delete -d '{ "id": "id_01" }'
```
