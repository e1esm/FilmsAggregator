package postgres

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Status int

const (
	SUCCESS Status = iota
	FAIL
)

func TestTransactionManager_AddAndGet(t *testing.T) {
	t.Parallel()
	tm := TransactionManager{transactions: make(map[uuid.UUID]pgx.Tx)}
	type test struct {
		name Status
		id   uuid.UUID
		tx   *pgxpool.Tx
		want pgx.Tx
	}

	tests := []test{
		{
			name: SUCCESS,
			id:   uuid.New(),
			tx:   &pgxpool.Tx{},
			want: &pgxpool.Tx{},
		},
		{
			name: FAIL,
			id:   uuid.New(),
			tx:   nil,
			want: &pgxpool.Tx{},
		},
	}

	for _, tmT := range tests {
		tm.Add(tmT.id, tmT.tx)
		got := tm.Get(tmT.id)
		if tmT.name == SUCCESS {
			assert.Equal(t, tmT.want, got)
		} else {
			assert.NotEqual(t, tmT.want, got)
		}
	}
}

func TestTransactionManager_Delete(t *testing.T) {
	tm := TransactionManager{transactions: make(map[uuid.UUID]pgx.Tx, 0)}
	t.Parallel()
	type test struct {
		name Status
		id   uuid.UUID
		tx   *pgxpool.Tx
		want pgx.Tx
	}

	tests := []test{
		{
			name: SUCCESS,
			id:   uuid.New(),
			tx:   &pgxpool.Tx{},
			want: nil,
		},
		{
			name: FAIL,
			id:   uuid.New(),
			tx:   nil,
			want: &pgxpool.Tx{},
		},
	}

	for _, tmT := range tests {
		tm.Add(tmT.id, tmT.tx)
		tm.Delete(tmT.id)
		got := tm.Get(tmT.id)
		if tmT.name == SUCCESS {
			assert.Equal(t, got, tmT.want)
		} else {
			assert.NotEqual(t, got, tmT.want)
		}
	}
}
