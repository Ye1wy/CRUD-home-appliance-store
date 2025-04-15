package psgrep

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/logger"
	"CRUD-HOME-APPLIANCE-STORE/internal/model/domain"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

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
	return nil, nil
}

func (r *productRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return nil, nil
}

func (r *productRepo) Update(ctx context.Context, product domain.Product) error {
	return nil
}

func (r *productRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}
