# go-http-server

Go で RPC スタイルの API を実装するサンプル

## 起動方法

```sh
docker build -t go-http-server .
docker run --rm -it -p 8000:8000 go-http-server
```

## 動作確認コマンド

```sh
curl -X POST localhost:8000/core/v1/task/get -d '{ "id": "id_01" }'
curl -X POST localhost:8000/core/v1/task/list
curl -X POST localhost:8000/core/v1/task/create -d '{ "title": "title_01" }'
curl -X POST localhost:8000/core/v1/task/update -d '{ "id": "id_01", "title": "title_01", "status": "DONE" }'
curl -X POST localhost:8000/core/v1/task/delete -d '{ "id": "id_01" }'
```

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

### API スキーマ定義

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

## アーキテクチャ

[The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)をベースに設計

### entity 層

エンティティとそれに紐づくビジネスルールをカプセル化する。
DB や外部 API の操作が必要なメソッドは interface のみ定義し、依存性逆転を行う。

usecase を見れば全体の処理の流れを追える状態にするため、極力複数のエンティティに跨るドメインサービスは使用しない。

### usecase 層

entity を使用して実際の一連の処理を記述する。

### interface 層

**inbound**

外部からのリクエストを受け取り、usecase を呼び出し、結果に応じてレスポンスを返す。

**outbound**

entity 層で定義した interface の実装を行う。DB や外部 API などとのやり取りを実装する。

### 各層の依存方向

```
interface/inbound -> usecase -> entity <- interface/outbound
```

## 開発フロー

`/core/v1/task/done`のように、先に interface/inbound と entity だけ作成し、usecase で mock データを返す状態を作る。

API 定義とモックサーバーを確認できる状態を先に作り、他のチームの作業に影響しないようにする。
