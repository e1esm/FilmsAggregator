package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"sync"
)

type TransactionManager struct {
	rwmx         sync.RWMutex
	transactions map[uuid.UUID]pgx.Tx
}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{transactions: make(map[uuid.UUID]pgx.Tx)}
}

func (tm *TransactionManager) Add(id uuid.UUID, tx pgx.Tx) {
	tm.rwmx.Lock()
	tm.transactions[id] = tx
	tm.rwmx.Unlock()
}

func (tm *TransactionManager) Get(id uuid.UUID) pgx.Tx {
	tm.rwmx.RLock()
	defer tm.rwmx.RUnlock()
	v, ok := tm.transactions[id]
	if ok {
		return v
	}
	return nil
}

func (tm *TransactionManager) Delete(id uuid.UUID) {
	tm.rwmx.Lock()
	delete(tm.transactions, id)
	tm.rwmx.Unlock()
}

func (tm *TransactionManager) VerifyAndGet(id uuid.UUID) (pgx.Tx, error) {
	tx := tm.Get(id)
	if tx == nil {
		return nil, fmt.Errorf("tracsaction wasn't started, neither was deleted")
	}
	return tx, nil
}
