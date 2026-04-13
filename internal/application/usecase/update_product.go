package usecase

import (
	"errors"
	"html"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrUnauthorizedUpdateProduct = errors.New("Unauthorized")
)

type UpdateProductRequest struct {
	ID          string `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Stock       int    `json:"stock"`
	SellerID    string `json:"-"`
}

type UpdateProductResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Price       int64  `json:"price"`
	Stock       int    `json:"stock"`
	Description string `json:"description"`
	SellerID    string `json:"seller_id"`
}

type UpdateProductUseCase struct {
	ProductRepository domain.ProductRepository
}

func NewUpdateProductUseCase(repo domain.ProductRepository) *UpdateProductUseCase {
	return &UpdateProductUseCase{
		ProductRepository: repo,
	}
}

func (uc *UpdateProductUseCase) Execute(req UpdateProductRequest) (*UpdateProductResponse, error) {
	product, err := uc.ProductRepository.GetProductByID(req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	if product.SellerID != req.SellerID {
		return nil, ErrUnauthorizedUpdateProduct
	}

	product.Name = html.EscapeString(req.Name)
	product.Description = html.EscapeString(req.Description)
	product.Price = req.Price
	product.Stock = req.Stock

	if err := product.Validate(); err != nil {
		return nil, err
	}

	err = uc.ProductRepository.UpdateProduct(product)
	if err != nil {
		return nil, err
	}

	return &UpdateProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Stock:       product.Stock,
		Description: product.Description,
		SellerID:    product.SellerID,
	}, nil
}
