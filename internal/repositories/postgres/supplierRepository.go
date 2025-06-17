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
	"github.com/jackc/pgx/v5/pgconn"
)

type SupplierRepo struct {
	*basePostgresRepository
}

func NewSupplierRepository(db DB, logger *logger.Logger) *SupplierRepo {
	repo := newBasePostgresRepository(db, logger)
	logger.Debug("postgres supplier repository is created")
	return &SupplierRepo{
		repo,
	}
}

func (r *SupplierRepo) Create(ctx context.Context, supplier *domain.Supplier) error {
	op := "repository.postgres.supplierRepository.Create"
	sqlStatement := "INSERT INTO supplier(name, address_id, phone_number) VALUES (@name, @address_id, @phone_number);"
	args := pgx.NamedArgs{
		"name":         supplier.Name,
		"address_id":   supplier.Address.Id,
		"phone_number": supplier.PhoneNumber,
	}

	_, err := r.db.Exec(ctx, sqlStatement, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			r.logger.Debug("Duplicate creation", "op", op)
			return fmt.Errorf("%s: unable to insert row: %w", op, crud_errors.ErrDuplicateKeyValue)
		}

		r.logger.Error("failed to create supplier", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unable to insert row: %v", op, err)
	}

	return nil
}

func (r *SupplierRepo) GetAll(ctx context.Context, limit, offset int) ([]domain.Supplier, error) {
	op := "repository.postgres.supplierRepository.GetAll"
	sqlStatement := `SELECT
		s.id,
		s.name,
		s.phone_number,
		a.id,
		a.country,
		a.city,
		a.street
		FROM supplier s
		LEFT JOIN address a ON s.address_id = a.id
		LIMIT @limit OFFSET @offset;`
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	rows, err := r.db.Query(ctx, sqlStatement, args)
	if err != nil {
		r.logger.Error("unable to query supplier: %v", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: unable to insert row: %v", op, err)
	}
	defer rows.Close()

	var suppliers []domain.Supplier

	for rows.Next() {
		supplier := domain.Supplier{
			Address: &domain.Address{},
		}

		err := rows.Scan(
			&supplier.Id,
			&supplier.Name,
			&supplier.PhoneNumber,
			&supplier.Address.Id,
			&supplier.Address.Country,
			&supplier.Address.City,
			&supplier.Address.Street,
		)
		if err != nil {
			r.logger.Warn("failed binding data", logger.Err(err), "op", op)
			continue
		}

		suppliers = append(suppliers, supplier)
	}

	if len(suppliers) == 0 {
		r.logger.Debug("supplier's not found", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	return suppliers, nil
}

func (r *SupplierRepo) GetByName(ctx context.Context, name string) (*domain.Supplier, error) {
	op := "repository.postgres.supplierRepository.GetById"
	sqlStatement := `SELECT
		s.id,
		s.name,
		s.phone_number,
		a.id,
		a.country,
		a.city,
		a.street
		FROM supplier s
		LEFT JOIN address a ON s.address_id = a.id
		WHERE s.name = @name;`
	arg := pgx.NamedArgs{
		"name": name,
	}

	row := r.db.QueryRow(ctx, sqlStatement, arg)
	supplier := domain.Supplier{
		Address: &domain.Address{},
	}
	err := row.Scan(
		&supplier.Id,
		&supplier.Name,
		&supplier.PhoneNumber,
		&supplier.Address.Id,
		&supplier.Address.Country,
		&supplier.Address.City,
		&supplier.Address.Street,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Debug("supplier not found", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	if err != nil {
		r.logger.Error("scan unable", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: scan failed: %v", op, err)
	}

	return &supplier, nil
}

func (r *SupplierRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	op := "repository.postgres.supplierRepository.GetById"
	sqlStatement := `SELECT
		s.id,
		s.name,
		s.phone_number,
		a.id,
		a.country,
		a.city,
		a.street
		FROM supplier s
		LEFT JOIN address a ON s.address_id = a.id
		WHERE s.id = @id;`
	arg := pgx.NamedArgs{
		"id": id,
	}

	row := r.db.QueryRow(ctx, sqlStatement, arg)
	supplier := domain.Supplier{
		Address: &domain.Address{},
	}
	err := row.Scan(
		&supplier.Id,
		&supplier.Name,
		&supplier.PhoneNumber,
		&supplier.Address.Id,
		&supplier.Address.Country,
		&supplier.Address.City,
		&supplier.Address.Street,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Debug("supplier not found", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	if err != nil {
		r.logger.Error("scan unable", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: scan failed: %v", op, err)
	}

	return &supplier, nil
}

func (r *SupplierRepo) Update(ctx context.Context, id, address uuid.UUID) error {
	op := "repositories.postgres.supplierRepository.Update"
	sqlStatement := "UPDATE supplier SET address_id=@address_id WHERE id=@id"
	arg := pgx.NamedArgs{
		"id":         id,
		"address_id": address,
	}

	r.logger.Debug("parameter", "id", id, "address_id", address)

	tag, err := r.db.Exec(ctx, sqlStatement, arg)
	if tag.RowsAffected() == 0 {
		r.logger.Debug("supplier not found", "op", op)
		return fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	if err != nil {
		r.logger.Error("failed execution update query", logger.Err(err), "op", op)
		return fmt.Errorf("%s: failed exec query: %v", op, err)
	}

	return nil
}

func (r *SupplierRepo) Delete(ctx context.Context, id uuid.UUID) error {
	op := "repository.postgres.supplierRepository.Delete"
	sqlStatement := "DELETE FROM supplier WHERE id=@id"
	arg := pgx.NamedArgs{
		"id": id,
	}

	_, err := r.db.Exec(ctx, sqlStatement, arg)
	if err != nil {
		r.logger.Error("execute sql statement is unable", logger.Err(err), "op", op)
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
