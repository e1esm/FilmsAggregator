package postgres

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TransactionManager struct {
	transactions map[uuid.UUID]pgx.Tx
}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{transactions: make(map[uuid.UUID]pgx.Tx)}
}
