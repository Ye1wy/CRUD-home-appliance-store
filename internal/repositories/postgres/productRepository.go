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

type ProductRepo struct {
	*basePostgresRepository
}

func NewProductRepository(conn pgx.Tx, logger *logger.Logger) *ProductRepo {
	repo := newBasePostgresRepository(conn, logger)
	logger.Debug("Postgres Product repository is created")
	return &ProductRepo{
		repo,
	}
}

func (r *ProductRepo) Create(ctx context.Context, product domain.Product) error {
	op := "repositories.postgres.productRepository.Create"
	sqlStatement := "INSERT INTO product(name, category, price, available_stock, supplier_id, image_id) VALUE (@name, @category, @price, @available_stock, @supplier_id, @image_id)"
	args := pgx.NamedArgs{
		"name":            product.Name,
		"category":        product.Category,
		"price":           product.Price,
		"available_stock": product.AvailableStock,
		"supplier_id":     product.SupplierId,
		"image_id":        product.ImageId,
	}

	_, err := r.db.Exec(ctx, sqlStatement, args)
	if err != nil {
		r.logger.Debug("Failed to create Product", "op", op)
		return fmt.Errorf("Product Repository: unable to insert row: %v", err)
	}

	return nil
}

func (r *ProductRepo) GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	op := "repositories.postgres.productRepository.GetAll"
	sqlStatement := "SELECT * FROM product LIMIT @limit OFFSET @offset"
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	rows, err := r.db.Query(ctx, sqlStatement, args)
	if err != nil {
		r.logger.Debug("query unvalable", logger.Err(err), "op", op)
		return nil, fmt.Errorf("Product Repo: query error %v", err)
	}

	var products []domain.Product

	for rows.Next() {
		var product domain.Product

		if err := rows.Scan(&product.Id, &product.Name, &product.Category,
			&product.Price, &product.AvailableStock, &product.LastUpdateDate,
			&product.SupplierId, &product.ImageId); err != nil {
			r.logger.Debug("Scan unable", logger.Err(err), "op", op)
			return nil, fmt.Errorf("Product Repo: Scan failed %v", err)
		}

		products = append(products, product)
	}

	if len(products) == 0 {
		r.logger.Debug("Product not found", "op", op)
		return nil, ErrNotFound
	}

	r.logger.Debug("All product is retrieved", "op", op)
	return products, nil
}

func (r *ProductRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	op := "repository.postgres.productRepository.GetById"
	sqlStatement := "SELECT * FROM product WHERE id=@id"
	arg := pgx.NamedArgs{
		"id": id,
	}

	row := r.db.QueryRow(ctx, sqlStatement, arg)
	product := domain.Product{}
	err := row.Scan(&product.Id, &product.Name, &product.Category,
		&product.Price, &product.AvailableStock, &product.LastUpdateDate,
		&product.SupplierId, &product.ImageId)
	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Debug("Product not found", "op", op)
		return nil, ErrNotFound
	}

	if err != nil {
		r.logger.Debug("Scan unable", logger.Err(err), "op", op)
		return nil, fmt.Errorf("Product Repo: Scan failed %v", err)
	}

	return &product, nil
}

func (r *ProductRepo) Update(ctx context.Context, id uuid.UUID, decrease int) error {
	op := "repository.postgres.productRepository.Update"
	sqlStatement := "UPDATE product SET available_stock = available_stock - @decrease WHERE id = @id"
	args := pgx.NamedArgs{
		"decrease": decrease,
	}

	tag, err := r.db.Exec(ctx, sqlStatement, args)
	if tag.RowsAffected() == 0 {
		r.logger.Debug("Product not found", "op", op)
		return ErrNotFound
	}

	if err != nil {
		r.logger.Debug("Execute sql statement for update stock is unable", logger.Err(err), "op", op)
		return fmt.Errorf("Product Repo: %v", err)
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
		r.logger.Debug("Execute sql statement for delete product is unable", logger.Err(err), "op", op)
		return fmt.Errorf("Product Repo: %v", err)
	}

	r.logger.Debug("Product deleted", "op", op)
	return nil
}
