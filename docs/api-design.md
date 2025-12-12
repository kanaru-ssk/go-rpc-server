# API 設計

## 概要

RPC スタイルの HTTP API

## 命名規則

パスは以下の命名規則に従う。

```
/<serviceName>/<version>/<useCaseName>/<methodName>

例: /core/v1/task/get
```

- 単語は `lowerCamelCase` で統一
- `<serviceName>` はマイクロサービス化を想定して設定、モノリス時は`core`というサービス名を使う。
- `<useCaseName>`, `<methodName>` は`usecase`ディレクトリ配下の構造体名、メソッド名に合わせる。
- `<methodName>`は常にパスに指定し、HTTP メソッドは安全性や冪等性などを判断するために使い分ける。

## HTTP の規約

[HTTP | MDN](https://developer.mozilla.org/docs/Web/HTTP) に従う。

以下は特に確認すべきもの

- [HTTP リクエストメソッド](https://developer.mozilla.org/docs/Web/HTTP/Reference/Methods)
- [HTTP レスポンスステータスコード](https://developer.mozilla.org/docs/Web/HTTP/Reference/Status)

`GET`, `DELETE`, `PUT`, `POST`を以下のルールで使い分ける。

| メソッド | 安全性 | 冪等性 | キャッシュ | リクエスト値 | 独自ルール                                                    |
| :------: | :----: | :----: | :--------: | :----------: | :------------------------------------------------------------ |
|   GET    |   O    |   O    |     O      |    query     | パス内のメソッド名に get または list のプレフィックスを付ける |
|  DELETE  |   X    |   O    |     X      |    query     | 物理削除で使用                                                |
|   PUT    |   X    |   O    |     X      |     body     | 部分更新でも冪等であれば PUT を使う                           |
|   POST   |   X    |   X    |    X\*     |     body     |                                                               |

- 安全かつ冪等な処理は GET
- 物理削除を行い、冪等な処理は DELETE
- 安全ではないが冪等性な処理は PUT
- それ以外は全て POST

## API スキーマ定義

OpenAPI などを使うとスキーマファイルの管理コストがかかるため用意しない。

代わりに、ソースコードを綺麗に保ち API のパス、リクエスト、レスポンスを読みやすくする。

`cmd/httpserver/main.go`を見れば全てのエンドポイントが分かる。

```go
// cmd/httpserver/main.go
mux.HandleFunc("GET /healthz", handler.HandleGetHealthz)

mux.HandleFunc("GET /core/v1/task/get", taskHandler.HandleGetV1)
mux.HandleFunc("GET /core/v1/task/list", taskHandler.HandleListV1)
mux.HandleFunc("POST /core/v1/task/create", taskHandler.HandleCreateV1)
mux.HandleFunc("PUT /core/v1/task/update", taskHandler.HandleUpdateV1)
mux.HandleFunc("DELETE /core/v1/task/delete", taskHandler.HandleDeleteV1)
mux.HandleFunc("PUT /core/v1/task/done", taskHandler.HandleDoneV1)
```

handler のコードを追えばリクエスト、レスポンスの内容が分かる。

```go
// interface/http/handler/task.go

// GET /core/v1/task/get
func (h *TaskHandler) HandleGetV1(w http.ResponseWriter, r *http.Request) {
	// 1番初めにリクエストとレスポンスの変数を定義してI/Oを読み取りやすくする。
	var query struct {
		ID string `query:"id"`
	}
	var successResponse response.Task
	var errorResponse response.Error

	ctx := r.Context()

	// 型のパースだけinterfaceレイヤーで行う。
	// 詳細なバリデーションはusecase以降で行う。
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
```

## その他の API 形式との比較

それぞれのメリット、デメリットはこのリポジトリの設計との比較。

**RESTful**

メリット

- 採用実績が多く、チームメンバー間での認識齟齬が生まれにくい

デメリット

- リソース指向の URI を使用する。同じリソースに対して多様な操作を実装する場合、拡張性に課題がある。

[**JSON-RPC**](https://www.jsonrpc.org/)

メリット

- コミュニティで仕様を規定しているため、独自ルールよりはチームメンバー間での認識齟齬が生まれにくい

デメリット

- パスにメソッドを含まれないのでログの検索性が下がる
- HTTP メソッドが POST 固定のため、安全性、冪等性、キャッシュ可否が読み取りにくい

**gRPC**

メリット

- protocol buffer というバイナリ形式で通信するため、データサイズを圧縮できる。

デメリット

- proto ファイルを元に生成したクライアントコードの使用が前提になるため、JavaScript の fetch API や、curl コマンドなどが使用できず、開発体験が悪い。
- HTTP メソッドが POST 固定のため、安全性、冪等性、キャッシュ可否が読み取りにくい

**GraphQL**

メリット

- クライアントから必要なデータを指定できるため、オーバーフェッチ、アンダーフェッチを回避できる。

デメリット

- Apollo などのライブラリを利用する前提になるため依存が増える。
- HTTP メソッドが POST 固定のため、安全性、冪等性、キャッシュ可否が読み取りにくい

**結論**

- gRPC, GraphQL のように依存が増えるものは開発体験が下がるので採用しない。
- RESTful は優れた設計だが、複雑なアプリケーションではメソッドの拡張性に課題がある。
- JSON-RPC は全て POST を使用するのため安全性、冪等性、キャッシュ可否が読み取りにくい
