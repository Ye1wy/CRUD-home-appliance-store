package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"context"
	"errors"
)

type ImagesService interface {
	AddImage(ctx context.Context, image *model.Image) error
	GetProductImages(ctx context.Context, id int) (*model.ProductImage, error)
	GetImageById(ctx context.Context, id string) (*model.Image, error)
	ChangeImage(ctx context.Context, id string, newImage []byte) error
	DeleteImageById(ctx context.Context, id int) error
}

type ImagesServiceImpl struct {
	Repo *repositories.ImagesRepository
}

func NewImagesServiceImpl(repo *repositories.ImagesRepository) *ImagesServiceImpl {
	return &ImagesServiceImpl{
		Repo: repo,
	}
}

func (s *ImagesServiceImpl) AddImage(ctx context.Context, image *model.Image) error {
	if len(image.Image) == 0 {
		return errors.New("image cannot be empty")
	}

	return s.Repo.AddImage(ctx, image)
}

func (s *ImagesServiceImpl) GetProductImages(ctx context.Context, id int) (*model.ProductImage, error) {
	image, err := s.Repo.GetProductImages(ctx, id)
	if err != nil {
		return nil, err
	}

	if image == nil {
		return nil, errors.New("product or image not found")
	}

	return image, nil
}

func (s *ImagesServiceImpl) GetImageById(ctx context.Context, id string) (*model.Image, error) {
	image, err := s.Repo.GetImageById(ctx, id)
	if err != nil {
		return nil, err
	}

	if image == nil {
		return nil, errors.New("image not found")
	}

	return image, nil
}

func (s *ImagesServiceImpl) ChangeImage(ctx context.Context, id string, newImage []byte) error {
	err := s.Repo.ChangeImage(ctx, id, newImage)
	if err != nil {
		return err
	}

	return nil
}

func (s *ImagesServiceImpl) DeleteImageById(ctx context.Context, id string) error {
	return s.Repo.DeleteImageById(ctx, id)
}
