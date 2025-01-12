package services

import (
	"CRUD-HOME-APPLIANCE-STORE/internal/model"
	"CRUD-HOME-APPLIANCE-STORE/internal/repositories"
	"context"
	"errors"
)

type ProductsService interface {
	AddProduct(ctx context.Context, product *model.Product) error
	GetAllProducts(ctx context.Context, limit, offset int) ([]model.Product, error)
	GetProductById(ctx context.Context, id int) (*model.Product, error)
	DecreaseStock(ctx context.Context, id int, decrease int) error
	DeleteProductById(ctx context.Context, id int) error
}

type ProductsServiceImpl struct {
	Rep *repositories.ProductsRepository
}

func NewProductService(rep *repositories.ProductsRepository) *ProductsServiceImpl {
	return &ProductsServiceImpl{
		Rep: rep,
	}
}

func (ps *ProductsServiceImpl) AddProduct(ctx context.Context, product model.Product) error {
	if product.Name == "" || product.Category == "" {
		return errors.New("product name and category cannot be emply")
	}

	if product.Price < 0 {
		return errors.New("product price cannot be negative")
	}

	return ps.Rep.AddProduct(ctx, &product)
}

func (ps *ProductsServiceImpl) GetAllProducts(ctx context.Context, limit, offset int) ([]model.Product, error) {
	if limit < 0 || offset < 0 {
		return nil, errors.New("limit and offset cannot be less of 0")
	}

	return ps.Rep.GetAllProducts(ctx, limit, offset)
}

func (ps *ProductsServiceImpl) GetProductById(ctx context.Context, id int) (*model.Product, error) {
	product, err := ps.Rep.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New("product not found")
	}

	return product, nil
}

func (ps *ProductsServiceImpl) DecreaseStock(ctx context.Context, id int, decrease int) error {
	if decrease <= 0 {
		return errors.New("decrease value must be greater that 0")
	}

	return ps.Rep.DecreaseParametr(ctx, id, decrease)
}

func (ps *ProductsServiceImpl) DeleteProductById(ctx context.Context, id int) error {
	return ps.Rep.DeleteProductById(ctx, id)
}
