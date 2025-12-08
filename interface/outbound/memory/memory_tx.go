package memory

import (
	"context"
	"sync"

	"github.com/kanaru-ssk/go-http-server/lib/tx"
)

// MemoryTxManager はメモリ内で簡易的にトランザクション境界を再現する。
type MemoryTxManager struct {
	mu *sync.RWMutex
}

func NewTxManager(mu *sync.RWMutex) tx.Manager {
	return &MemoryTxManager{mu: mu}
}

// WithinTx は排他制御を行いながら処理を実行する。
// ※実際のロールバックはできない（atomicityは保証しない）
func (m *MemoryTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return fn(ctx)
}
