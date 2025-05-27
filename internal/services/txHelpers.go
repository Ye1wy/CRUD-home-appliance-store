package services

import (
	crud_errors "CRUD-HOME-APPLIANCE-STORE/internal/errors"
	"CRUD-HOME-APPLIANCE-STORE/internal/uow"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func safeDeleteAddress(
	ctx context.Context,
	tx pgx.Tx,
	id uuid.UUID,
	deleteFunc func(ctx context.Context, id uuid.UUID) error,
	log *logger.Logger,
	op string,
	savepointName string,
) error {
	savepoint := fmt.Sprintf("SAVEPOINT %s;", savepointName)
	_, err := tx.Exec(ctx, savepoint)
	if err != nil {
		log.Debug("unable to set SAVEPOINT before delete address", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unable to set SAVEPOINT: %v", op, err)
	}

	if err := deleteFunc(ctx, id); err != nil {
		if errors.Is(err, crud_errors.ErrForeignKeyViolation) {
			backToSave := fmt.Sprintf("ROLLBACK TO SAVEPOINT %s;", savepointName)
			_, err := tx.Exec(ctx, backToSave)
			if err != nil {
				log.Debug("unable back to SAVEPOINT after try delete address", logger.Err(err), "op", op)
				return fmt.Errorf("%s: unable back to SAVEPOINT: %v", op, err)
			}

			return nil
		}

		log.Debug("unexpected error during delete: rollback to SAVEPOINT is unavailable", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unexpected error from delete address: %v", op, err)
	}

	return nil
}

func getReposiotry(tx uow.Transaction, name uow.RepositoryName, log *logger.Logger) (uow.Repository, error) {
	repo, err := tx.Get(name)
	if err != nil {
		return nil, fmt.Errorf("get address repository generator is unable: %v", err)
	}

	repoGen := repo.(uow.RepositoryGenerator)(tx.GetTX(), log)
	return repoGen, nil
}
