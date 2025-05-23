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

type ImageRepository struct {
	*basePostgresRepository
}

func NewImageRepository(db DB, logger *logger.Logger) *ImageRepository {
	repo := newBasePostgresRepository(db, logger)
	return &ImageRepository{
		repo,
	}
}

func (r *ImageRepository) Create(ctx context.Context, image domain.Image) error {
	op := "repository.postgres.imageRepository.Create"
	sqlStatement := "INSERT INTO image(image) VALUES (@image)"
	args := pgx.NamedArgs{"image": image.Image}

	_, err := r.db.Exec(ctx, sqlStatement, args)
	if err != nil {
		r.logger.Debug("failed to create image", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unable to insert row: %v", op, err)
	}

	return nil
}

func (r *ImageRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Image, error) {
	op := "repostiory.postgres.imageRepository.GetAll"
	sqlStatement := "SELECT * FROM image LIMIT @limit OFFSET @offset"
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	rows, err := r.db.Query(ctx, sqlStatement, args)
	if err != nil {
		r.logger.Debug("failed to get all", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: query error: %v", op, err)
	}
	defer rows.Close()

	var images []domain.Image

	for rows.Next() {
		var image domain.Image

		if err := rows.Scan(&image.Id, &image.Image); err != nil {
			return nil, fmt.Errorf("%s: failed to bind data %v", op, err)
		}

		images = append(images, image)
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	return images, nil
}

func (r *ImageRepository) GetById(ctx context.Context, id uuid.UUID) (*domain.Image, error) {
	op := "repository.postgres.imageRepositoru.GetById"
	sqlStatement := "SELECT * FROM image WHERE id = $id"
	arg := pgx.NamedArgs{
		"id": id,
	}

	row := r.db.QueryRow(ctx, sqlStatement, arg)
	image := &domain.Image{}
	err := row.Scan(&image.Id, &image.Image)
	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Debug("image not found", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	if err != nil {
		r.logger.Debug("failed get image by id", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: scan failed: %v", op, err)
	}

	return image, nil
}

func (r *ImageRepository) Update(ctx context.Context, image *domain.Image) error {
	op := "repository.postgres.imageRepository.Update"
	sqlStatement := "UPDATE image SET image = @image WHERE id = @id"
	args := pgx.NamedArgs{
		"id":    image.Id,
		"image": image.Image,
	}

	tag, err := r.db.Exec(ctx, sqlStatement, args)
	if tag.RowsAffected() == 0 {
		r.logger.Debug("image not found", "op", op)
		return fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	if err != nil {
		r.logger.Debug("failed update image by id", logger.Err(err), "op", op)
		return fmt.Errorf("%s: failed exec query: %v", op, err)
	}

	return nil
}

func (r *ImageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	op := "repository.postgres.productRepository.Delete"
	sqlStatement := "DELETE FROM image WHERE id=@id"
	arg := pgx.NamedArgs{
		"id": id,
	}

	_, err := r.db.Exec(ctx, sqlStatement, arg)
	if err != nil {
		r.logger.Debug("failed delete image by id", logger.Err(err), "op", op)
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
