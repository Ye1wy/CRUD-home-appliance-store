package postgres

import (
	crud_errors "CRUD-HOME-APPLIANCE-STORE/internal/errors"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AddressRepo struct {
	*basePostgresRepository
}

func NewAddressRepository(db DB, log *logger.Logger) *AddressRepo {
	baseRepo := newBasePostgresRepository(db, log)
	log.Debug("address repo is created")
	return &AddressRepo{baseRepo}
}

func (r *AddressRepo) Create(ctx context.Context, address *domain.Address) error {
	op := "repository.postgres.addressRepository.Create"
	sqlInsert := `INSERT 
		INTO address(country, city, street)
		VALUES (@country, @city, @street) 
		ON CONFLICT DO NOTHING
		RETURNING id`

	args := pgx.NamedArgs{
		"country": address.Country,
		"city":    address.City,
		"street":  address.Street,
	}

	err := r.db.QueryRow(ctx, sqlInsert, args).Scan(&address.Id)
	if err == nil {
		return nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Debug("failed to insert address", logger.Err(err), "op", op)
		return fmt.Errorf("%s: %v", op, err)
	}

	sqlSelect := `SELECT id FROM address WHERE country = @country AND city = @city AND street = @street`

	err = r.db.QueryRow(ctx, sqlSelect, args).Scan(&address.Id)
	if err != nil {
		r.logger.Error("failed to get existing address id", logger.Err(err), "op", op)
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (r *AddressRepo) Delete(ctx context.Context, id uuid.UUID) error {
	op := "repository.postgres.addressRepository.Delete"
	sqlStatement := `DELETE FROM address WHERE id = @id`
	arg := pgx.NamedArgs{
		"id": id,
	}

	ct, err := r.db.Exec(ctx, sqlStatement, arg)

	if ct.RowsAffected() == 0 {
		r.logger.Warn("No affected row, but error skiped", logger.Err(err), "op", op)
		return fmt.Errorf("%s: no affected row. all errors from exec is skiped. %w", op, crud_errors.ErrForeignKeyViolation)
	}

	return nil
}
