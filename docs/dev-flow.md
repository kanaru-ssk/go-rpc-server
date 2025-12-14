# 開発フロー

## 1. interface

`/core/v1/task/done`のように、先に interface だけ作成して mock データを返す状態を作る。

```go
// POST /core/v1/task/done
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
```

API 定義とモックサーバーを確認できる状態を先に作り、他のチームの作業に影響しないようにする。

## 2. domain, usecase

domain, usecase で DB や外部 API に依存しないロジックを実装する。

## 3. infrastructure

domain で定義した interface に従って実際の DB 操作や外部 API との通信などの実装をする。

## 4. cmd/httpserver

各レイヤーのインスタンスを作成し、`dependencyInjection`で依存性注入を行い、mock から実際の API 動作に切り替える。

## 補足

スケジュールや他チームの作業に余裕があれば、domain, usecase を先に実装し、DB や外部 API、API 形式などの決定は後回しにしても良い。
