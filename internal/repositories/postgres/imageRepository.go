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

type ImageRepo struct {
	*basePostgresRepository
}

func NewImageRepository(db DB, logger *logger.Logger) *ImageRepo {
	repo := newBasePostgresRepository(db, logger)
	return &ImageRepo{
		repo,
	}
}

func (r *ImageRepo) Create(ctx context.Context, image *domain.Image) error {
	op := "repository.postgres.imageRepository.Create"
	sqlInsert := `INSERT
		INTO image(title, data)
		VALUES (@title, @image)
		RETURNING id;`
	args := pgx.NamedArgs{
		"title": image.Title,
		"image": image.Data,
	}

	err := r.db.QueryRow(ctx, sqlInsert, args).Scan(&image.Id)
	if err != nil {
		r.logger.Error("failed to create image", logger.Err(err), "op", op)
		return fmt.Errorf("%s: unable to insert row: %v", op, err)
	}

	return nil
}

func (r *ImageRepo) GetAll(ctx context.Context, limit, offset int) ([]domain.Image, error) {
	op := "repostiory.postgres.imageRepository.GetAll"
	sqlStatement := `SELECT * FROM image LIMIT @limit OFFSET @offset;`
	r.logger.Debug("check limit and offset", "limit", limit, "offset", offset)
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	rows, err := r.db.Query(ctx, sqlStatement, args)
	if err != nil {
		r.logger.Error("failed to get all", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: query error: %v", op, err)
	}
	defer rows.Close()

	var images []domain.Image

	for rows.Next() {
		var image domain.Image

		if err := rows.Scan(&image.Id, &image.Title, &image.Data); err != nil {
			r.logger.Warn("failed to bind data", logger.Err(err), "op", op)
			continue
		}

		r.logger.Debug("extracted data", "image id", image.Id, "image title", image.Title, "op", op)

		images = append(images, image)
	}

	if len(images) == 0 {
		r.logger.Debug("No content", "op", op)
		return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	return images, nil
}

func (r *ImageRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.Image, error) {
	op := "repository.postgres.imageRepositoru.GetById"
	sqlStatement := "SELECT * FROM image WHERE id = @id"
	arg := pgx.NamedArgs{
		"id": id,
	}

	row := r.db.QueryRow(ctx, sqlStatement, arg)
	image := domain.Image{}
	err := row.Scan(&image.Id, &image.Title, &image.Data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Debug("image not found", "op", op)
			return nil, fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
		}

		r.logger.Error("failed get image by id", logger.Err(err), "op", op)
		return nil, fmt.Errorf("%s: scan failed: %v", op, err)
	}

	return &image, nil
}

func (r *ImageRepo) Update(ctx context.Context, image *domain.Image) error {
	op := "repository.postgres.imageRepository.Update"
	sqlStatement := `UPDATE image SET title = @title, data = @image WHERE id = @id`
	args := pgx.NamedArgs{
		"id":    image.Id,
		"title": image.Title,
		"image": image.Data,
	}

	tag, err := r.db.Exec(ctx, sqlStatement, args)
	if tag.RowsAffected() == 0 {
		r.logger.Debug("image not found", "op", op)
		return fmt.Errorf("%s: %w", op, crud_errors.ErrNotFound)
	}

	if err != nil {
		r.logger.Error("failed update image by id", logger.Err(err), "op", op)
		return fmt.Errorf("%s: failed exec query: %v", op, err)
	}

	return nil
}

func (r *ImageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	op := "repository.postgres.productRepository.Delete"
	sqlStatement := "DELETE FROM image WHERE id=@id"
	arg := pgx.NamedArgs{
		"id": id,
	}

	_, err := r.db.Exec(ctx, sqlStatement, arg)
	if err != nil {
		r.logger.Error("failed delete image by id", logger.Err(err), "op", op)
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
