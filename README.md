# go-rpc-server

Go で RPC スタイルの API を実装するサンプル

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

## API 設計

### 命名規則

全て POST で統一し、パスは以下の命名規則に従う。

```
POST /serviceName/v1/usecaseName/methodName
```

- 単語は lowerCamel
- serviceName はマイクロサービス化を想定して設定、モノリス時は`core`というサービス名を使う。
- usecaseName, methodName は`usecase`ディレクトリ配下の構造体名、メソッド名に合わせる。
- get, list から始まる methodName は安全性・冪等性を担保する。

### API スキーマ定義

OpenAPI などを使うとスキーマファイルの管理コストがかかるため用意しない。

代わりに、ソースコードを綺麗に保ち API のパス、リクエスト、レスポンスを読みやすくする。

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

// POST /core/v1/task/get
func (h *TaskHandler) HandleGetV1(w http.ResponseWriter, r *http.Request) {
	var request struct {
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

	// 500
	default:
		slog.ErrorContext(ctx, "httphandler.TaskHandler.HandleGetV1", "err", err)
		httpresponse.RenderJson(ctx, w, http.StatusInternalServerError, h.errorMapper.MapErrorResponse(errorresponse.ErrInternalServerError))
	}
}
```

## アーキテクチャ

### interface 層

外部からのリクエストを受け取り、usecase を呼び出し、結果に応じてレスポンスを返す。

### usecase 層

domain を使用して実際の一連の処理を記述する。

### domain 層

- entity
- value object
- factory
- repository

などを定義。

repository ではデータ永続化や取得の interface のみ定義し、domain 層が特定の DB や API などの外部ツールに依存しないようにする。

usecase を見れば全体の処理の流れを追えるようにするため、ドメインサービスは使用しない。

### infrastructure 層

domain 層で定義した interface の実装を行う。DB や外部 API などとのやり取りを実装する。

### 各層の依存方向

```
interface -> usecase -> domain <- infrastructure
```

## 開発フロー

domain 層と usecase 層から作成して、infrastructure 層、最後に interface 層を実装するフローが理想。

しかし、実際の開発現場は納期に追われ、フロントエンドの開発チームの作業を止めないように API スキーマを先に決めたり、mock サーバーを用意してフロントエンドから API との連携をテストできる状態にする必要がある。

その場合、`/core/v1/task/done`のように先に interface と domain の entity だけ作成し、usecase で mock データを返す状態を作る。

こうすることで、フロントエンドチームは API 定義とモックサーバーを確認でき、バックエンドチームの負担も増やさずに済む。
