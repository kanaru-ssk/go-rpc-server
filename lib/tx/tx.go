package tx

import "context"

// Tx はデータストア固有のトランザクションハンドルを表し、
// リポジトリが必要に応じて具象型へキャストして利用する。
type Tx any

type Manager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error
}
