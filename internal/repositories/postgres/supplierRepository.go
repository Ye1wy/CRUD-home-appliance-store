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

func (r *SupplierRepo) Create(ctx context.Context, supplier domain.Supplier) error {
	op := "repository.postgres.supplierRepository.Create"
	sqlStatement := "INSERT INTO supplier(name, address_id, phone_number) VALUE (@name, @address_id, @phone_number)"
	args := pgx.NamedArgs{
		"name":         supplier.Name,
		"address_id":   supplier.AddressId,
		"phone_number": supplier.PhoneNumber,
	}

	_, err := r.db.Exec(ctx, sqlStatement, args)
	if err != nil {
		r.logger.Debug("failed to create supplier", logger.Err(err), "op", op)
		return fmt.Errorf("supplier repository: unable to insert row: %v", err)
	}

	return nil
}

func (r *SupplierRepo) GetAll(ctx context.Context, limit, offset int) ([]domain.Supplier, error) {
	op := "repository.postgres.supplierRepository.GetAll"
	sqlStatement := "SELECT * FROM client LIMIT @limit OFFSET @offset"
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	rows, err := r.db.Query(ctx, sqlStatement, args)
	if err != nil {
		r.logger.Debug("unable to query supplier: %v", logger.Err(err), "op", op)
		return nil, fmt.Errorf("supplier repository: error to get all client: %v", err)
	}
	defer rows.Close()

	var suppliers []domain.Supplier

	for rows.Next() {
		var supplier domain.Supplier

		if err := rows.Scan(&supplier.Id, &supplier.Name, &supplier.AddressId, &supplier.PhoneNumber); err != nil {
			return nil, fmt.Errorf("supplier repository: failed to bind data in GetAll: %v", err)
		}

		suppliers = append(suppliers, supplier)
	}

	if len(suppliers) == 0 {
		r.logger.Debug("supplier's not found", "op", op)
		return nil, ErrNotFound
	}

	r.logger.Debug("all supplier retrived", "op", op)
	return suppliers, nil
}

func (r *SupplierRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	op := "repository.postgres.supplierRepository.GetById"
	sqlStatement := "SELECT * FROM product WHERE id=@id"
	arg := pgx.NamedArgs{
		"id": id,
	}

	row := r.db.QueryRow(ctx, sqlStatement, arg)
	supplier := domain.Supplier{}
	err := row.Scan(&supplier.Id, &supplier.Name, &supplier.AddressId, &supplier.PhoneNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Debug("supplier not found", "op", op)
		return nil, ErrNotFound
	}

	if err != nil {
		r.logger.Debug("Scan unable", logger.Err(err), "op", op)
		return nil, fmt.Errorf("supplier repository: ccan failed %v", err)
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
		return ErrNotFound
	}

	if err != nil {
		r.logger.Debug("failed execution update query", logger.Err(err), "op", op)
		return fmt.Errorf("supplier repository: failed exec update query %v", err)
	}

	r.logger.Debug("data updated", "op", op)
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
		r.logger.Debug("execute sql statement for delete product is unable", logger.Err(err), "op", op)
		return fmt.Errorf("supplier repo: %v", err)
	}

	r.logger.Debug("supplier deleted", "op", op)
	return nil
}
