package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/dto"
	"CRUD-HOME-APPLIANCE-STORE/internal/mapper"
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService interface {
	AddProduct(ctx context.Context, dto *dto.ProductDTO) (*model.Product, error)
	GetAllProducts(ctx context.Context, limit, offset int) ([]dto.ProductDTO, error)
	GetProductById(ctx context.Context, id string) (*dto.ProductDTO, error)
	DecreaseStock(ctx context.Context, id string, decrease int) error
	DeleteProductById(ctx context.Context, id string) error
}

type ProductServiceImpl struct {
	Repo repositories.ProductRepository
}

func NewProductService(rep repositories.ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{
		Repo: rep,
	}
}

func (ps *ProductServiceImpl) AddProduct(ctx context.Context, dto *dto.ProductDTO) (*model.Product, error) {
	product := mapper.ToProductModel(dto)
	result, err := ps.Repo.AddProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		product.Id = objectID.Hex()

	} else {
		return nil, fmt.Errorf("failed to parse inserted ID")
	}

	return product, nil
}

func (ps *ProductServiceImpl) GetAllProducts(ctx context.Context, limit, offset int) ([]dto.ProductDTO, error) {
	if limit < 0 || offset < 0 {
		return nil, errors.New("limit and offset cannot be less of 0")
	}

	products, err := ps.Repo.GetAllProducts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	dto := mapper.ToProductDTOs(products)

	return dto, nil
}

func (ps *ProductServiceImpl) GetProductById(ctx context.Context, id string) (*dto.ProductDTO, error) {
	product, err := ps.Repo.GetProductById(ctx, id)
	if product == nil && err == nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	dto := mapper.ToProductDTO(product)

	return dto, nil
}

func (ps *ProductServiceImpl) DecreaseStock(ctx context.Context, id string, decrease int) error {
	if decrease <= 0 {
		return errors.New("decrease value must be greater that 0")
	}

	return ps.Repo.DecreaseParametr(ctx, id, decrease)
}

func (ps *ProductServiceImpl) DeleteProductById(ctx context.Context, id string) error {
	return ps.Repo.DeleteProductById(ctx, id)
}
