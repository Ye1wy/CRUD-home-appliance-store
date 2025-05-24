package postgres

import (
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

func (r *AddressRepo) Create(ctx context.Context, address domain.Address) (uuid.UUID, error) {
	op := "repository.postgres.addressRepository.Create"
	sqlInsert := `INSERT INTO address(country, city, street)
					 VALUES (@country, @city, @street) 
					 ON CONFLICT DO NOTHING
					 RETURNING id`

	args := pgx.NamedArgs{
		"country": address.Country,
		"city":    address.City,
		"street":  address.Street,
	}

	var id uuid.UUID
	err := r.db.QueryRow(ctx, sqlInsert, args).Scan(&id)
	if err == nil {
		return id, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Debug("failed to insert address", logger.Err(err), "op", op)
		return uuid.Nil, fmt.Errorf("%s: %v", op, err)
	}

	sqlSelect := `SELECT id FROM address WHERE country = @country AND city = @city AND street = @street`

	err = r.db.QueryRow(ctx, sqlSelect, args).Scan(&id)
	if err != nil {
		r.logger.Debug("failed to get existing address id", logger.Err(err), "op", op)
		return uuid.Nil, fmt.Errorf("%s: %v", op, err)
	}

	return id, nil
}

func (r *AddressRepo) Delete(ctx context.Context, id uuid.UUID) error {
	op := "repository.postgres.addressRepository.Delete"
	sqlStatement := `DELETE FROM address WHERE id = @id`
	arg := pgx.NamedArgs{
		"id": id,
	}

	_, err := r.db.Exec(ctx, sqlStatement, arg)
	if err != nil {
		r.logger.Error("execute sql statement for delete address is unable", logger.Err(err), "op", op)
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
