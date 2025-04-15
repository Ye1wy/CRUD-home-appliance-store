package psgrep

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"CRUD-HOME-APPLIANCE-STORE/pkg/logger"
	"context"
	"errors"
	"fmt"
	"runtime"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var ErrProductNotFound = errors.New("Product not found in database")

type WriteProductRepo interface {
	Create(ctx context.Context, product domain.Product) error
	Update(ctx context.Context, id uuid.UUID, decrease int) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type productRepo struct {
	*basePostgresRepository
}

func NewProductRepository(conn *pgx.Conn, logger *logger.Logger) *productRepo {
	repo := newBasePostgresRepository(conn, logger)
	logger.Debug("Postgres Product repository is created")
	return &productRepo{
		repo,
	}
}

func (r *productRepo) Create(ctx context.Context, product domain.Product) error {
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

func (r *productRepo) GetAll(ctx context.Context, limit, offset int) ([]domain.Product, error) {
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
		return nil, ErrProductNotFound
	}

	r.logger.Debug("All product is retrieved", "op", op)
	return products, nil
}

func (r *productRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
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
		return nil, ErrProductNotFound
	}

	if err != nil {
		r.logger.Debug("Scan unable", logger.Err(err), "op", op)
		return nil, fmt.Errorf("Product Repo: Scan failed %v", err)
	}

	return &product, nil
}

func (r *productRepo) Update(ctx context.Context, id uuid.UUID, decrease int) error {
	op := "repository.postgres.productRepository.Update"
	sqlStatement := "UPDATE product SET available_stock = available_stock - @decrease WHERE id = @id"
	args := pgx.NamedArgs{
		"decrease": decrease,
	}

	tag, err := r.db.Exec(ctx, sqlStatement, args)
	if tag.RowsAffected() == 0 {
		r.logger.Debug("Product not found", "op", op)
		return ErrProductNotFound
	}

	if err != nil {
		r.logger.Debug("Execute sql statement for update stock is unable", logger.Err(err), "op", op)
		return fmt.Errorf("Product Repo: %v", err)
	}

	return nil
}

func (r *productRepo) Delete(ctx context.Context, id uuid.UUID) error {
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

func (r *productRepo) UnitOfWork(ctx context.Context, fn func(WriteProductRepo) error) error {
	op := "repository.postgres.productRepository.UnitOfWork"
	trx, err := r.db.Begin(ctx)
	if err != nil {
		r.logger.Debug("Transaction begin error", logger.Err(err), "op", op)
		return fmt.Errorf("Client unit of work %v", err)
	}

	defer func(ctx context.Context) {
		if p := recover(); p != nil {
			if p := recover(); p != nil {
				r.logger.Error("Panic catch in unit of work, rolling back", "op", op)
				_ = trx.Rollback(ctx)

				switch e := p.(type) {
				case runtime.Error:
					r.logger.Error("Runtime panic", "error", e, "op", op)
					panic(e)
				case error:
					err = fmt.Errorf("panic err: %v", err)
					return
				default:
					r.logger.Error("Not handle error", "error", e, "op", op)
					panic(e)
				}
			}

			if err != nil {
				r.logger.Error("Transaction error!", logger.Err(err), "op", op)
				_ = trx.Rollback(ctx)
			} else {
				r.logger.Debug("Transaction commited", "op", op)
				_ = trx.Commit(ctx)
			}
		}
	}(ctx)

	newStore := NewProductRepository(r.db, r.logger)
	return fn(newStore)
}
