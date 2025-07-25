package postgres

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/uow"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// All the presented structures are wrappers of the pgx driver
// to abstract from the specific implementation of the driver.

// Implements part of the pgx driver transaction
type pgxTx struct {
	tx pgx.Tx
}

func NewPgxTx(tx pgx.Tx) uow.Tx {
	return &pgxTx{
		tx: tx,
	}
}

func (p *pgxTx) Begin(ctx context.Context) (uow.Tx, error) {
	tx, err := p.tx.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("pgxTx begin: %v", err)
	}

	return &pgxTx{
		tx: tx,
	}, nil
}

func (p *pgxTx) Rollback(ctx context.Context) error {
	if err := p.tx.Rollback(ctx); err != nil {
		return fmt.Errorf("pgxTx rollback: %v", err)
	}

	return nil
}

func (p *pgxTx) Commit(ctx context.Context) error {
	if err := p.tx.Commit(ctx); err != nil {
		return fmt.Errorf("pgxTx commit: %v", err)
	}

	return nil
}

type Row struct {
}
