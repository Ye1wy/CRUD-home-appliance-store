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

type ProductRepo struct {
	*basePostgresRepository
}

func NewProductRepository(db DB, logger *logger.Logger) *ProductRepo {
	repo := newBasePostgresRepository(db, logger)
	logger.Debug("Postgres Product repository is created")
	return &ProductRepo{
		repo,
	}
}

func (r *ProductRepo) Create(ctx context.Context, product *domain.Product) error {
	op := "repositories.postgres.productRepository.Create"
	sqlStatement := `
	INSERT INTO product(name, category, price, available_stock,  supplier_id, image_id) 
	VALUES (@name, @category, @price, @available_stock, @supplier_id, @image_id)
	RETURNING id;`
	args := pgx.NamedArgs{
		"name":            product.Name,
		"category":        product.Category,
		"price":           product.Price,
		"available_stock": product.AvailableStock,
		"supplier_id":     product.Supplier.Id,
		"image_id":        product.Image.Id,
	}

	err := r.db.QueryRow(ctx, sqlStatement, args).Scan(&product.Id)
	if err != nil {
		r.logger.Error("failed to create product", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unable to insert row: %v", op, err)
	}

	return nil
}

func (r *ProductRepo) GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	op := "repositories.postgres.productRepository.GetAll"
	sqlStatement := `SELECT
		p.id,
		p.name,
		p.category,
		p.price,
		p.available_stock,
		s.id,
		s.name,
		s.phone_number,
		a.id,
		a.country,
		a.city,
		a.street,
		i.id,
		i.data
		FROM product p
		LEFT JOIN supplier s ON p.supplier_id = s.id
		LEFT JOIN address a ON s.address_id = a.id
		LEFT JOIN image i ON p.image_id = i.id 
		LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	rows, err := r.db.Query(ctx, sqlStatement, args)
	if err != nil {
		r.logger.Error("query unvalable", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: query error: %v", op, err)
	}

	var products []domain.Product

	for rows.Next() {
		var (
			product                                                            domain.Product
			imageData                                                          []byte
			supplierAddressId                                                  *uuid.UUID
			supplierAddressCountry, supplierAddressCity, supplierAddressStreet *string
		)

		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Category,
			&product.Price,
			&product.AvailableStock,
			&product.Supplier.Id,
			&product.Supplier.Name,
			&product.Supplier.PhoneNumber,
			&supplierAddressId,
			&supplierAddressCountry,
			&supplierAddressCity,
			&supplierAddressStreet,
			&product.Image.Id,
			&imageData,
		)

		if err != nil {
			r.logger.Warn("scan unable", logger.Err(err), "op", op)
			continue
		}

		if supplierAddressId != nil && imageData != nil {
			product.Supplier.Address = &domain.Address{
				Id:      *supplierAddressId,
				Country: *supplierAddressCountry,
				City:    *supplierAddressCity,
				Street:  *supplierAddressStreet,
			}
			product.Image.Data = imageData

		} else {
			if imageData == nil {
				r.logger.Error("WRONG! Unthinkable, a image without data, this can't be", "op", op)
				return nil, fmt.Errorf("%s: Image data is %w", op, crud_errors.ErrProductImageDataEmpty)
			}

			r.logger.Error("WRONG! Unthinkable, a supplier without an address, this can't be", "op", op)
			return nil, fmt.Errorf("%s: Supplier Address is %w", op, crud_errors.ErrProductSupplerAddressEmpty)
		}

		products = append(products, product)
	}

	if len(products) == 0 {
		r.logger.Debug("product not found", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	return products, nil
}

func (r *ProductRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	op := "repository.postgres.productRepository.GetById"
	sqlStatement := `SELECT
		p.id,
		p.name,
		p.category,
		p.price,
		p.available_stock,
		s.id,
		s.name,
		s.phone_number,
		a.id,
		a.country,
		a.city,
		a.street,
		i.id,
		i.data
		FROM product p
		LEFT JOIN supplier s ON p.supplier_id = s.id
		LEFT JOIN address a ON s.address_id = a.id
		LEFT JOIN image i ON p.image_id = i.id 
		WHERE p.id = @id`
	arg := pgx.NamedArgs{
		"id": id,
	}

	row := r.db.QueryRow(ctx, sqlStatement, arg)
	var (
		product                                                            domain.Product
		imageData                                                          []byte
		supplierAddressId                                                  *uuid.UUID
		supplierAddressCountry, supplierAddressCity, supplierAddressStreet *string
	)
	err := row.Scan(
		&product.Id,
		&product.Name,
		&product.Category,
		&product.Price,
		&product.AvailableStock,
		&product.Supplier.Id,
		&product.Supplier.Name,
		&product.Supplier.PhoneNumber,
		&supplierAddressId,
		&supplierAddressCountry,
		&supplierAddressCity,
		&supplierAddressStreet,
		&product.Image.Id,
		&imageData,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Debug("product not found", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	if supplierAddressId != nil && imageData != nil {
		product.Supplier.Address = &domain.Address{
			Id:      *supplierAddressId,
			Country: *supplierAddressCountry,
			City:    *supplierAddressCity,
			Street:  *supplierAddressStreet,
		}
		product.Image.Data = imageData

	} else {
		if imageData == nil {
			r.logger.Error("WRONG! Unthinkable, a image without data, this can't be", "op", op)
			return nil, fmt.Errorf("%s: Image data is %w", op, crud_errors.ErrProductImageDataEmpty)
		}

		r.logger.Error("WRONG! Unthinkable, a supplier without an address, this can't be", "op", op)
		return nil, fmt.Errorf("%s: Supplier Address is %w", op, crud_errors.ErrProductSupplerAddressEmpty)
	}

	if err != nil {
		r.logger.Error("scan unable", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: scan failed: %v", op, err)
	}

	return &product, nil
}

func (r *ProductRepo) Update(ctx context.Context, id uuid.UUID, decrease int) error {
	op := "repository.postgres.productRepository.Update"
	sqlStatement := "UPDATE product SET available_stock = available_stock - @decrease WHERE id = @id"
	args := pgx.NamedArgs{
		"decrease": decrease,
		"id":       id,
	}

	tag, err := r.db.Exec(ctx, sqlStatement, args)
	if tag.RowsAffected() == 0 {
		r.logger.Debug("product not found", "op", op)
		return fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	if err != nil {
		r.logger.Error("execute sql statement for update stock is unable", logger.Err(err), "op", op)
		return fmt.Errorf("%s: failed exec query: %v", op, err)
	}

	return nil
}

func (r *ProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
	op := "repository.postgres.productRepository.Delete"
	sqlStatement := "DELETE FROM product WHERE id=@id"
	arg := pgx.NamedArgs{
		"id": id,
	}

	_, err := r.db.Exec(ctx, sqlStatement, arg)
	if err != nil {
		r.logger.Error("execute sql statement for delete product is unable", logger.Err(err), "op", op)
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
